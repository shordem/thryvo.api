package user_repository

import (
	"strings"

	"github.com/google/uuid"

	"github.com/shordem/api.thryvo/lib/database"
	models "github.com/shordem/api.thryvo/model"
	"github.com/shordem/api.thryvo/repository"
)

type UserRepositoryInterface interface {
	Create(user models.User) (models.User, error)
	FindAllUsers(pageable repository.Pageable) ([]models.User, repository.Pagination, error)
	FindUserById(uuid uuid.UUID) (models.User, error)
	FindUserByEmail(email string) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
	DeleteUser(uuid uuid.UUID) error
}

type userRepository struct {
	database database.DatabaseInterface
}

func NewUserRepository(database database.DatabaseInterface) UserRepositoryInterface {
	return &userRepository{database: database}
}

// Create implements UserRepositoryInterface.
func (u *userRepository) Create(user models.User) (models.User, error) {
	user.Prepare()

	err := u.database.Connection().Create(&user).Error

	if err != nil {

		return models.User{}, err
	}

	return user, err
}

// DeleteUser implements UserRepositoryInterface.
func (u *userRepository) DeleteUser(uuid uuid.UUID) error {

	user, err := u.FindUserById(uuid)

	if err != nil {
		return err
	}

	err = u.database.Connection().Delete(&user).Error

	if err != nil {

		return err
	}

	return nil
}

// FindAllUsers implements UserRepositoryInterface.
func (u *userRepository) FindAllUsers(pageable repository.Pageable) (users []models.User, pagination repository.Pagination, err error) {
	var user models.User

	pagination.CurrentPage = int64(pageable.Page)
	pagination.TotalItems = 0
	pagination.TotalPages = 1

	search := strings.TrimSpace(pageable.Search)
	offset := (pageable.Page - 1) * pageable.Size
	model := u.database.Connection().Model(&user)

	// Apply search filters
	if len(search) > 0 {
		model = model.Where("first_name LIKE ? OR last_name LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Get total items
	if err = model.Count(&pagination.TotalItems).Error; err != nil {
		return nil, pagination, err
	}

	// apply pagination
	paginatedQuery := model.
		Select("id", "first_name", "last_name", "referral_code", "email", "is_email_verified", "created_at").
		Offset(int(offset)).
		Limit(int(pageable.Size)).
		Order(pageable.SortBy + " " + pageable.SortDirection)

	// execute query
	if err = paginatedQuery.Find(&users).Error; err != nil {
		return nil, pagination, err
	}

	// calculate total pages
	if pagination.TotalItems > 0 {
		pagination.TotalPages = (pagination.TotalItems + int64(pageable.Size) - 1) / int64(pageable.Size)
	} else {
		pagination.TotalPages = 1
	}

	return users, pagination, nil
}

// FindUserByEmail implements UserRepositoryInterface.
func (u *userRepository) FindUserByEmail(email string) (user models.User, err error) {
	err = u.database.Connection().Model(&models.User{}).Where("email = ?", email).First(&user).Error

	return user, err
}

// FindUserById implements UserRepositoryInterface.
func (u *userRepository) FindUserById(uuid uuid.UUID) (user models.User, err error) {
	err = u.database.Connection().Model(&models.User{}).Where("id = ?", uuid).First(&user).Error

	return user, err
}

// UpdateUser implements UserRepositoryInterface.
func (u *userRepository) UpdateUser(user models.User) (models.User, error) {

	checkRow, err := u.FindUserById(user.ID)

	if err != nil {
		return checkRow, err
	}

	err = u.database.Connection().Model(&checkRow).Updates(user).Error

	if err != nil {

		return models.User{}, err
	}

	return checkRow, err

}
