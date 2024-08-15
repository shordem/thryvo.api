package core_repository

import (
	"errors"

	"github.com/google/uuid"

	"github.com/shordem/api.thryvo/lib/database"
	"github.com/shordem/api.thryvo/model"
)

type FolderRepositoryInterface interface {
	CreateFolder(folder model.Folder) (model.Folder, error)
	FindFoldersByUserId(userId uuid.UUID) ([]model.Folder, error)
	FindFoldersByParentId(userId uuid.UUID, parentId uuid.UUID) ([]model.Folder, error)
	FindFolderById(id uuid.UUID) (model.Folder, error)
	UpdateFolder(folder model.Folder) (model.Folder, error)
	DeleteFolder(id uuid.UUID, userId uuid.UUID) error
}

type folderRepository struct {
	database database.DatabaseInterface
}

func NewFolderRepository(database database.DatabaseInterface) FolderRepositoryInterface {
	return &folderRepository{database: database}
}

// CreateFolder implements FolderRepositoryInterface.
func (f *folderRepository) CreateFolder(folder model.Folder) (model.Folder, error) {
	folder.Prepare()

	if folder.ParentID != nil {
		if err := f.database.Connection().First(&model.Folder{}, "id = ?", folder.ParentID).Error; err != nil {
			return model.Folder{}, errors.New("parent folder not found")
		}
	}

	if err := f.database.Connection().Create(&folder).Error; err != nil {
		return model.Folder{}, err
	}

	return folder, nil
}

// FindFoldersByUserId implements FolderRepositoryInterface.
func (f *folderRepository) FindFoldersByUserId(userId uuid.UUID) ([]model.Folder, error) {
	var folders []model.Folder

	if err := f.database.Connection().
		Where("user_id = ?", userId).
		Order("id DESC").
		Find(&folders).
		Error; err != nil {
		return nil, err
	}

	return folders, nil
}

// FindFoldersByParentId implements FolderRepositoryInterface.
func (f *folderRepository) FindFoldersByParentId(userId uuid.UUID, parentId uuid.UUID) ([]model.Folder, error) {
	var folders []model.Folder

	if err := f.database.Connection().
		Where("user_id = ? AND parent_id = ?", userId, parentId).
		Order("id DESC").
		Find(&folders).
		Error; err != nil {
		return nil, err
	}

	return folders, nil
}

// FindFolderById implements FolderRepositoryInterface.
func (f *folderRepository) FindFolderById(id uuid.UUID) (model.Folder, error) {
	var folder model.Folder

	if err := f.database.Connection().First(&folder, "id = ?", id).Error; err != nil {
		return model.Folder{}, err
	}

	return folder, nil
}

// UpdateFolder implements FolderRepositoryInterface.
func (f *folderRepository) UpdateFolder(folder model.Folder) (model.Folder, error) {

	if folder.ParentID != nil {
		if err := f.database.Connection().First(&model.Folder{}, "id = ?", folder.ParentID).Error; err != nil {
			return model.Folder{}, errors.New("parent folder not found")
		}
	}

	if err := f.database.Connection().Updates(&folder).Error; err != nil {
		return model.Folder{}, err
	}

	return folder, nil
}

// FIXME: delete folder should set all files and folders parent_id to null
// DeleteFolder implements FolderRepositoryInterface.
func (f *folderRepository) DeleteFolder(id uuid.UUID, userId uuid.UUID) error {

	if err := f.database.Connection().Model(&model.Folder{}).Delete("id = ? AND user_id = ?", id, userId).Error; err != nil {
		return err
	}

	return nil
}
