package service

import (
	"gojet/models"
	"log/slog"
)

// UserService 用户服务
type UserService struct {
	userRepo UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// CreateUser 创建新用户
func (s *UserService) CreateUser(name string) (*models.User, error) {
	user := &models.User{
		Name: name,
	}

	if err := s.userRepo.Create(user); err != nil {
		slog.Error("创建用户失败", "用户", user)
		return nil, err
	}
	return user, nil
}

// CreateInitialData 创建初始学生数据
func (s *UserService) CreateInitialData() error {
	existingUsers, err := s.userRepo.GetAll()
	if err == nil && len(existingUsers) > 0 {
		return nil
	}

	users := []*models.User{
		{Name: "包子"},
		{Name: "玉米"},
		{Name: "花卷"},
		{Name: "吐司"},
	}

	if err := s.userRepo.CreateBatch(users); err != nil {
		return err
	}
	return nil
}

// GetAllUsers 获取所有用户
func (s *UserService) GetAllUsers() ([]*models.User, error) {
	return s.userRepo.GetAll()
}

// GetUserByID 根据 ID 获取用户
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

// UpdateUser 更新用户信息
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

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}
