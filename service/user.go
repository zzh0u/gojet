package service

import (
	"gojet/models"
	"gojet/util/apperror"
	"gojet/util/response"
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
		slog.Error("创建用户失败", "用户", user, "error", err)
		return nil, apperror.Wrap(err, 500, response.MsgUserCreateFailed)
	}
	return user, nil
}

// CreateInitialData 创建初始学生数据
func (s *UserService) CreateInitialData() error {
	existingUsers, err := s.userRepo.GetAll()
	if err != nil {
		// 重要：遇到错误应该返回，而不是继续执行
		return apperror.Wrap(err, 500, "检查现有数据失败")
	}
	if len(existingUsers) > 0 {
		slog.Info("初始数据已存在，跳过插入")
		return nil // 数据已存在，跳过
	}

	users := []*models.User{
		{Name: "包子"},
		{Name: "玉米"},
		{Name: "花卷"},
		{Name: "吐司"},
	}

	if err := s.userRepo.CreateBatch(users); err != nil {
		slog.Error("创建初始数据失败", "error", err)
		return apperror.Wrap(err, 500, response.MsgDBInsertError)
	}

	slog.Info("初始数据创建成功", "count", len(users))
	return nil
}

// GetAllUsers 获取所有用户
func (s *UserService) GetAllUsers() ([]*models.User, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, apperror.Wrap(err, 500, "获取用户列表失败")
	}
	return users, nil
}

// GetUserByID 根据 ID 获取用户
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		// DAO 层已经包装了错误，直接返回
		return nil, err
	}
	return user, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(id uint, name string) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	user.Name = name

	if err := s.userRepo.Update(user); err != nil {
		slog.Error("更新用户失败", "id", id, "error", err)
		return nil, apperror.Wrap(err, 500, response.MsgUserUpdateFailed)
	}

	slog.Info("更新用户成功", "id", id, "name", name)
	return user, nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id uint) error {
	if err := s.userRepo.Delete(id); err != nil {
		slog.Error("删除用户失败", "id", id, "error", err)
		return apperror.Wrap(err, 500, response.MsgUserDeleteFailed)
	}
	slog.Info("删除用户成功", "id", id)
	return nil
}
