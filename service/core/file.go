package core_service

import (
	"mime/multipart"

	"github.com/shordem/api.thryvo/dto"
	"github.com/shordem/api.thryvo/lib/config"
)

type FileServiceInterface interface {
	UploadFile(file *multipart.FileHeader) (string, error)
	GetFile(fileName string) (dto.GetFileDTO, error)
}

type fileService struct {
	mediaConfig config.FileConfigInterface
}

func NewFileService(mediaConfig config.FileConfigInterface) FileServiceInterface {
	return &fileService{mediaConfig: mediaConfig}
}

func (s *fileService) UploadFile(file *multipart.FileHeader) (string, error) {
	return s.mediaConfig.UploadFile(file)
}

func (s *fileService) GetFile(fileName string) (dto.GetFileDTO, error) {
	return s.mediaConfig.GetObject(fileName)
}
