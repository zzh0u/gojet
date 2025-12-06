package service

import (
	"gojet/models"
)

// UserRepository defines the interface for user data operations that service layer needs
type UserRepository interface {
	// Create creates a new user
	Create(user *models.User) error

	// CreateBatch creates multiple users
	CreateBatch(users []*models.User) error

	// GetAll retrieves all users
	GetAll() ([]*models.User, error)

	// GetByID retrieves a user by ID
	GetByID(id uint) (*models.User, error)

	// Update updates a user
	Update(user *models.User) error

	// Delete soft-deletes a user
	Delete(id uint) error
}
