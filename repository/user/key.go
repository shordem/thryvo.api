package user_repository

import (
	"github.com/google/uuid"

	"github.com/shordem/api.thryvo/lib/database"
	"github.com/shordem/api.thryvo/model"
)

type KeyRepositoryInterface interface {
	Create(key model.Key) (model.Key, error)
	FindKeyById(uuid uuid.UUID) (model.Key, error)
	FindKeyByUserID(uuid uuid.UUID) (model.Key, error)
	FindUserIDByKey(key string) (model.Key, error)
	UpdateKey(key model.Key) (model.Key, error)
	DeleteKey(uuid uuid.UUID) error
}

type keyRepository struct {
	database database.DatabaseInterface
}

func NewKeyRepository(database database.DatabaseInterface) KeyRepositoryInterface {
	return &keyRepository{database: database}
}

// Create implements KeyRepositoryInterface.
func (k *keyRepository) Create(key model.Key) (model.Key, error) {
	key.Prepare()

	err := k.database.Connection().Create(&key).Error

	if err != nil {
		return model.Key{}, err
	}

	return key, err
}

// FindKeyById implements KeyRepositoryInterface.
func (k *keyRepository) FindKeyById(uuid uuid.UUID) (model.Key, error) {
	var key model.Key

	err := k.database.Connection().Where("id = ?", uuid).First(&key).Error

	return key, err
}

// FindKeyByUserID implements KeyRepositoryInterface.
func (k *keyRepository) FindKeyByUserID(uuid uuid.UUID) (model.Key, error) {
	var key model.Key

	err := k.database.Connection().Where("user_id = ?", uuid).First(&key).Error

	return key, err
}

// FindUserIDByKey implements KeyRepositoryInterface.
func (k *keyRepository) FindUserIDByKey(key string) (model.Key, error) {
	var keyModel model.Key

	err := k.database.Connection().Where("key = ?", key).First(&k).Error

	return keyModel, err
}

// DeleteKey implements KeyRepositoryInterface.
func (k *keyRepository) DeleteKey(uuid uuid.UUID) error {

	key, err := k.FindKeyById(uuid)

	if err != nil {
		return err
	}

	err = k.database.Connection().Delete(&key).Error

	if err != nil {
		return err
	}

	return nil
}

// UpdateKey implements KeyRepositoryInterface.
func (k *keyRepository) UpdateKey(key model.Key) (model.Key, error) {

	err := k.database.Connection().Save(&key).Error

	if err != nil {
		return model.Key{}, err
	}

	return key, err
}
