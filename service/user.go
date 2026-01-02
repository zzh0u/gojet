package service

import (
	"gojet/dao"
	"gojet/models"
	"gojet/util/apperror"
	"log/slog"
)

// userRepo 包级变量，存储用户仓库实例
var userRepo *dao.UserRepository

// InitService 初始化服务层，设置依赖的数据仓库
func InitService(repo *dao.UserRepository) {
	userRepo = repo
}

// CreateUser 使用完整的用户信息创建用户
func CreateUser(user *models.User) (*models.User, error) {
	if err := userRepo.Create(user); err != nil {
		slog.Error("创建用户失败", "用户", user.Username, "error", err)
		return nil, apperror.Wrap(err, 500, apperror.UserCreateFailed)
	}

	slog.Info("创建用户成功", "id", user.ID, "username", user.Username)
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
		{Username: "包子", NickName: "包子", Password: "123456", Email: "baozi@example.com"},
		{Username: "玉米", NickName: "玉米", Password: "123456", Email: "corn@example.com"},
		{Username: "花卷", NickName: "花卷", Password: "123456", Email: "flower@example.com"},
		{Username: "吐司", NickName: "吐司", Password: "123456", Email: "toast@example.com"},
	}

	// 对每个用户的密码进行哈希处理
	for _, user := range users {
		hashedPassword, err := models.HashPassword(user.Password)
		if err != nil {
			slog.Error("密码哈希失败", "username", user.Username, "error", err)
			return apperror.Wrap(err, 500, "密码哈希失败")
		}
		user.Password = hashedPassword
	}

	if err := userRepo.CreateBatch(users); err != nil {
		slog.Error("创建初始数据失败", "error", err)
		return apperror.Wrap(err, 500, apperror.DBInsertError)
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

	user.Username = name

	if err := userRepo.Update(user); err != nil {
		slog.Error("更新用户失败", "id", id, "error", err)
		return nil, apperror.Wrap(err, 500, apperror.UserUpdateFailed)
	}

	slog.Info("更新用户成功", "id", id, "name", name)
	return user, nil
}

// DeleteUser 删除用户
func DeleteUser(id uint) error {
	if err := userRepo.Delete(id); err != nil {
		slog.Error("删除用户失败", "id", id, "error", err)
		return apperror.Wrap(err, 500, apperror.UserDeleteFailed)
	}
	slog.Info("删除用户成功", "id", id)
	return nil
}
