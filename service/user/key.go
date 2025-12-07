package user_service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/shordem/api.thryvo/lib/helper"
	"github.com/shordem/api.thryvo/model"
	user_repository "github.com/shordem/api.thryvo/repository/user"
)

type SubscriptionChecker interface {
	GetSubscriptionStatus(ctx context.Context, userID uuid.UUID) (*SubscriptionStatusResponse, error)
}

type SubscriptionStatusResponse struct {
	IsActive bool
}

type keyService struct {
	keyRepository       user_repository.KeyRepositoryInterface
	subscriptionChecker SubscriptionChecker
}

type KeyServiceInterface interface {
	GetKey(userId uuid.UUID) (string, error)
	CreateKey(userId uuid.UUID) error
}

func NewKeyService(keyRepository user_repository.KeyRepositoryInterface, subscriptionChecker SubscriptionChecker) KeyServiceInterface {
	return &keyService{
		keyRepository:       keyRepository,
		subscriptionChecker: subscriptionChecker,
	}
}

func (k *keyService) GetKey(userId uuid.UUID) (string, error) {
	// Check if user has active subscription
	if k.subscriptionChecker != nil {
		status, err := k.subscriptionChecker.GetSubscriptionStatus(context.Background(), userId)
		if err != nil {
			return "", errors.New("failed to check subscription status")
		}

		if !status.IsActive {
			return "", errors.New("active subscription required to access API key")
		}
	}

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
