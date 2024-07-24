package core_service

import (
	"mime/multipart"

	"github.com/shordem/api.thryvo/dto"
	"github.com/shordem/api.thryvo/lib/config"
)

type MediaServiceInterface interface {
	UploadMedia(file *multipart.FileHeader) (string, error)
	GetMedia(fileName string) (dto.GetFileDTO, error)
}

type mediaService struct {
	mediaConfig config.FileConfigInterface
}

func NewMediaService(mediaConfig config.FileConfigInterface) MediaServiceInterface {
	return &mediaService{mediaConfig: mediaConfig}
}

func (s *mediaService) UploadMedia(file *multipart.FileHeader) (string, error) {
	return s.mediaConfig.UploadFile(file)
}

func (s *mediaService) GetMedia(fileName string) (dto.GetFileDTO, error) {
	return s.mediaConfig.GetObject(fileName)
}
