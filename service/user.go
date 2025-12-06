package service

import (
	"log/slog"

	"gojet/models"
)

// UserService handles business logic for users
type UserService struct {
	userRepo UserRepository
	logger   *slog.Logger
}

// NewUserService creates a new UserService
func NewUserService(userRepo UserRepository, logger *slog.Logger) *UserService {
	return &UserService{userRepo: userRepo, logger: logger}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(name string) (*models.User, error) {
	s.logger.Info("Creating new user", "name", name)

	user := &models.User{
		Name: name,
	}

	if err := s.userRepo.Create(user); err != nil {
		s.logger.Error("Failed to create user", "error", err, "name", name)
		return nil, err
	}

	s.logger.Info("User created successfully", "id", user.ID, "name", name)
	return user, nil
}

// CreateInitialData creates initial student data
func (s *UserService) CreateInitialData() error {
	s.logger.Info("Initializing sample student data")

	// Check if data already exists
	existingUsers, err := s.userRepo.GetAll()
	if err == nil && len(existingUsers) > 0 {
		s.logger.Info("Database already contains data, skipping initialization")
		return nil
	}

	users := []*models.User{
		{Name: "包子"},
		{Name: "玉米"},
		{Name: "花卷"},
		{Name: "吐司"},
	}

	if err := s.userRepo.CreateBatch(users); err != nil {
		s.logger.Error("Failed to create initial data", "error", err)
		return err
	}

	s.logger.Info("Initial data created successfully", "count", len(users))
	return nil
}

// GetAllUsers returns all users
func (s *UserService) GetAllUsers() ([]*models.User, error) {
	return s.userRepo.GetAll()
}

// GetUserByID returns a user by ID
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

// UpdateUser updates a user's information
func (s *UserService) UpdateUser(id uint, name string) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	user.Name = name

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}
