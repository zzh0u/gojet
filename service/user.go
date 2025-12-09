package service

import (
	"gojet/dao"
	"gojet/models"
	"gojet/util/apperror"
	"gojet/util/response"
	"log/slog"
)

// userRepo 包级变量，存储用户仓库实例
var userRepo *dao.UserRepository

// InitService 初始化服务层，设置依赖的数据仓库
func InitService(repo *dao.UserRepository) {
	userRepo = repo
}

// CreateUser 创建新用户
func CreateUser(name string) (*models.User, error) {
	user := &models.User{
		Name: name,
	}

	if err := userRepo.Create(user); err != nil {
		slog.Error("创建用户失败", "用户", user, "error", err)
		return nil, apperror.Wrap(err, 500, response.MsgUserCreateFailed)
	}
	return user, nil
}

// CreateInitialData 创建初始学生数据
func CreateInitialData() error {
	existingUsers, err := userRepo.GetAll()
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

	if err := userRepo.CreateBatch(users); err != nil {
		slog.Error("创建初始数据失败", "error", err)
		return apperror.Wrap(err, 500, response.MsgDBInsertError)
	}

	slog.Info("初始数据创建成功", "count", len(users))
	return nil
}

// GetAllUsers 获取所有用户
func GetAllUsers() ([]*models.User, error) {
	users, err := userRepo.GetAll()
	if err != nil {
		return nil, apperror.Wrap(err, 500, "获取用户列表失败")
	}
	return users, nil
}

// GetUserByID 根据 ID 获取用户
func GetUserByID(id uint) (*models.User, error) {
	user, err := userRepo.GetByID(id)
	if err != nil {
		// DAO 层已经包装了错误，直接返回
		return nil, err
	}
	return user, nil
}

// UpdateUser 更新用户信息
func UpdateUser(id uint, name string) (*models.User, error) {
	user, err := userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	user.Name = name

	if err := userRepo.Update(user); err != nil {
		slog.Error("更新用户失败", "id", id, "error", err)
		return nil, apperror.Wrap(err, 500, response.MsgUserUpdateFailed)
	}

	slog.Info("更新用户成功", "id", id, "name", name)
	return user, nil
}

// DeleteUser 删除用户
func DeleteUser(id uint) error {
	if err := userRepo.Delete(id); err != nil {
		slog.Error("删除用户失败", "id", id, "error", err)
		return apperror.Wrap(err, 500, response.MsgUserDeleteFailed)
	}
	slog.Info("删除用户成功", "id", id)
	return nil
}
