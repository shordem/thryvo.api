package user_service

import (
	"github.com/google/uuid"
	"github.com/shordem/api.thryvo/lib/helper"
	"github.com/shordem/api.thryvo/model"
	user_repository "github.com/shordem/api.thryvo/repository/user"
)

type keyService struct {
	keyRepository user_repository.KeyRepositoryInterface
}

type KeyServiceInterface interface {
	GetKey(userId uuid.UUID) (string, error)
	CreateKey(userId uuid.UUID) error
}

func NewKeyService(keyRepository user_repository.KeyRepositoryInterface) KeyServiceInterface {
	return &keyService{keyRepository: keyRepository}
}

func (k *keyService) GetKey(userId uuid.UUID) (string, error) {

	key, err := k.keyRepository.FindKeyByUserID(userId)

	return key.Key, err
}

func (k *keyService) CreateKey(userId uuid.UUID) error {

	key := model.Key{
		UserID: userId,
		Key:    helper.GenerateRandomHexStr(32),
	}

	_, err := k.keyRepository.Create(key)

	return err
}
