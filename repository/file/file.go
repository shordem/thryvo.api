package core_repository

import (
	"strings"

	"github.com/google/uuid"

	"github.com/shordem/api.thryvo/lib/database"
	"github.com/shordem/api.thryvo/model"
	"github.com/shordem/api.thryvo/repository"
)

type FileRepositoryInterface interface {
	Create(file model.File) (model.File, error)
	FindAllFiles(pageable repository.Pageable) ([]model.File, repository.Pagination, error)
	FindFileById(uuid uuid.UUID) (model.File, error)
	UpdateFile(file model.File) (model.File, error)
	DeleteFile(uuid uuid.UUID) error
}

type fileRepository struct {
	database database.DatabaseInterface
}

func NewFileRepository(database database.DatabaseInterface) FileRepositoryInterface {
	return &fileRepository{database: database}
}

// Create implements FileRepositoryInterface.
func (f *fileRepository) Create(file model.File) (model.File, error) {
	file.Prepare()

	err := f.database.Connection().Create(&file).Error

	if err != nil {
		return model.File{}, err
	}

	return file, err
}

// FindFileById implements FileRepositoryInterface.
func (f *fileRepository) FindFileById(uuid uuid.UUID) (model.File, error) {
	var file model.File

	err := f.database.Connection().Where("id = ?", uuid).First(&file).Error

	return file, err
}

// DeleteFile implements FileRepositoryInterface.
func (f *fileRepository) DeleteFile(uuid uuid.UUID) error {

	file, err := f.FindFileById(uuid)

	if err != nil {
		return err
	}

	err = f.database.Connection().Delete(&file).Error

	if err != nil {
		return err
	}

	return nil
}

// FindAllFiles implements FileRepositoryInterface.
func (f *fileRepository) FindAllFiles(pageable repository.Pageable) (files []model.File, pagination repository.Pagination, err error) {
	var file model.File

	pagination.CurrentPage = int64(pageable.Page)
	pagination.TotalItems = 0
	pagination.TotalPages = 1

	search := strings.TrimSpace(pageable.Search)
	offset := (pageable.Page - 1) * pageable.Size
	model := f.database.Connection().Model(&file)

	if len(search) > 0 {
		model = model.Where("original_name LIKE ?", "%"+search+"%")
	}

	if err = model.Count(&pagination.TotalItems).Error; err != nil {
		return nil, pagination, err
	}

	// apply pagination
	paginatedQuery := model.
		Offset(offset).
		Limit(pageable.Size).
		Order(pageable.SortBy + " " + pageable.SortDirection)

	if err = paginatedQuery.Find(&files).Error; err != nil {
		return nil, pagination, err
	}

	if pagination.TotalItems > 0 {
		pagination.TotalPages = (pagination.TotalItems + int64(pageable.Size) - 1) / int64(pageable.Size)
	} else {
		pagination.TotalPages = 1
	}

	return files, pagination, err
}

// UpdateFile implements FileRepositoryInterface.
func (f *fileRepository) UpdateFile(file model.File) (model.File, error) {
	file.Prepare()

	err := f.database.Connection().Save(&file).Error

	if err != nil {
		return model.File{}, err
	}

	return file, err
}
