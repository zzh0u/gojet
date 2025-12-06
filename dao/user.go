package dao

import (
	"errors"

	"gojet/models"
	"gojet/service"

	"gorm.io/gorm"
)

// userRepository implements service.UserRepository using GORM
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) service.UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(user *models.User) error {
	result := r.db.Create(user)
	return result.Error
}

// CreateBatch creates multiple users
func (r *userRepository) CreateBatch(users []*models.User) error {
	result := r.db.CreateInBatches(users, len(users))
	return result.Error
}

// GetAll retrieves all users (excluding soft-deleted)
func (r *userRepository) GetAll() ([]*models.User, error) {
	var users []*models.User
	result := r.db.Find(&users)
	return users, result.Error
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, result.Error
}

// Update updates a user
func (r *userRepository) Update(user *models.User) error {
	result := r.db.Save(user)
	return result.Error
}

// Delete soft-deletes a user
func (r *userRepository) Delete(id uint) error {
	result := r.db.Delete(&models.User{}, id)
	return result.Error
}
