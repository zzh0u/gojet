package api

import (
	"gojet/models"
)

// UserService defines the interface for user operations that API layer needs
type UserService interface {
	// CreateUser creates a new user
	CreateUser(name string) (*models.User, error)

	// CreateInitialData creates initial student data
	CreateInitialData() error

	// GetAllUsers returns all users
	GetAllUsers() ([]*models.User, error)

	// GetUserByID returns a user by ID
	GetUserByID(id uint) (*models.User, error)

	// UpdateUser updates a user's information
	UpdateUser(id uint, name string) (*models.User, error)

	// DeleteUser deletes a user by ID
	DeleteUser(id uint) error
}
