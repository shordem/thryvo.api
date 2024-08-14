package core_service

import (
	"github.com/google/uuid"

	"github.com/shordem/api.thryvo/dto"
	"github.com/shordem/api.thryvo/model"
	core_repository "github.com/shordem/api.thryvo/repository/core"
	user_repository "github.com/shordem/api.thryvo/repository/user"
)

type FolderServiceInterface interface {
	CreateFolder(folderDto dto.FolderDTO) (dto.FolderDTO, error)
	FindFoldersByUserId(userId uuid.UUID) ([]dto.FolderDTO, error)
	FindFoldersByParentId(userId uuid.UUID, parentId uuid.UUID) ([]dto.FolderDTO, error)
	UpdateFolder(folderDto dto.FolderDTO) (dto.FolderDTO, error)
	DeleteFolder(id uuid.UUID, userId uuid.UUID) error
}

type folderService struct {
	folderRepository core_repository.FolderRepositoryInterface
	userRepository   user_repository.UserRepositoryInterface
}

func NewFolderService(
	folderRepository core_repository.FolderRepositoryInterface,
	userRepository user_repository.UserRepositoryInterface,
) FolderServiceInterface {
	return &folderService{
		folderRepository: folderRepository,
		userRepository:   userRepository,
	}
}

func (f *folderService) ConvertToDTO(folder model.Folder) dto.FolderDTO {
	var folderDto dto.FolderDTO

	folderDto.ID = folder.ID
	folderDto.UserID = folder.UserID
	folderDto.ParentID = folder.ParentID
	folderDto.Name = folder.Name
	folderDto.CreatedAt = folder.CreatedAt
	folderDto.UpdatedAt = folder.UpdatedAt
	folderDto.DeletedAt = folder.DeletedAt.Time

	return folderDto
}

func (f *folderService) ConvertToModel(folderDto dto.FolderDTO) model.Folder {
	var folder model.Folder

	folder.ID = folderDto.ID
	folder.UserID = folderDto.UserID
	folder.ParentID = folderDto.ParentID
	folder.Name = folderDto.Name
	folder.CreatedAt = folderDto.CreatedAt
	folder.UpdatedAt = folderDto.UpdatedAt
	folder.DeletedAt.Time = folderDto.DeletedAt

	return folder
}

func (f *folderService) CreateFolder(folderDto dto.FolderDTO) (dto.FolderDTO, error) {
	folder := f.ConvertToModel(folderDto)

	folder, err := f.folderRepository.CreateFolder(folder)
	if err != nil {
		return dto.FolderDTO{}, err
	}

	return f.ConvertToDTO(folder), nil
}

func (f *folderService) FindFoldersByUserId(userId uuid.UUID) ([]dto.FolderDTO, error) {
	folders, err := f.folderRepository.FindFoldersByUserId(userId)
	if err != nil {
		return nil, err
	}

	folderDtos := []dto.FolderDTO{}
	for _, folder := range folders {
		folderDtos = append(folderDtos, f.ConvertToDTO(folder))
	}

	return folderDtos, nil
}

func (f *folderService) FindFoldersByParentId(userId uuid.UUID, parentId uuid.UUID) ([]dto.FolderDTO, error) {
	folders, err := f.folderRepository.FindFoldersByParentId(userId, parentId)
	if err != nil {
		return nil, err
	}

	folderDtos := []dto.FolderDTO{}
	for _, folder := range folders {
		folderDtos = append(folderDtos, f.ConvertToDTO(folder))
	}

	return folderDtos, nil
}

func (f *folderService) UpdateFolder(folderDto dto.FolderDTO) (dto.FolderDTO, error) {
	folder := f.ConvertToModel(folderDto)

	folder, err := f.folderRepository.UpdateFolder(folder)
	if err != nil {
		return dto.FolderDTO{}, err
	}

	return f.ConvertToDTO(folder), nil
}

func (f *folderService) DeleteFolder(id uuid.UUID, userId uuid.UUID) error {
	return f.folderRepository.DeleteFolder(id, userId)
}
