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
	// 检查用户名是否已存在
	existingUser, err := userRepo.GetUserByUserName(user.Username)
	if err != nil && err.Error() != apperror.RecordNotFound {
		return nil, apperror.Wrap(err, 500, apperror.DBQueryError)
	}
	if existingUser != nil {
		return nil, apperror.New(409, apperror.UserNameExists)
	}

	existingUserByEmail, err := userRepo.GetUserByEmail(user.Email)
	if err != nil && err.Error() != apperror.RecordNotFound {
		return nil, apperror.Wrap(err, 500, apperror.DBQueryError)
	}
	if existingUserByEmail != nil {
		return nil, apperror.New(409, apperror.EmailExists)
	}

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
func GetUserByID(id int) (*models.User, error) {
	user, err := userRepo.GetByID(id)
	if err != nil {
		// DAO 层已经包装了错误，直接返回
		return nil, err
	}
	return user, nil
}

// UpdateUserRequest 更新用户请求结构体
type UpdateUserRequest struct {
	Username string `json:"username" binding:"required"`
	NickName string `json:"nick_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// UpdateUser 更新用户信息
func UpdateUser(id int, req UpdateUserRequest) (*models.User, error) {
	user, err := userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 如果用户名有变化，检查新用户名是否已被其他用户使用
	if user.Username != req.Username {
		existingUser, err := userRepo.GetUserByUserName(req.Username)
		if err != nil && err.Error() != apperror.RecordNotFound {
			// 查询出错（不是记录不存在的错误）
			return nil, apperror.Wrap(err, 500, apperror.DBQueryError)
		}
		// 如果找到的用户不是当前用户，说明用户名已被占用
		if existingUser != nil && existingUser.ID != id {
			return nil, apperror.New(409, apperror.UserNameExists)
		}
	}

	// 如果邮箱有变化，检查新邮箱是否已被其他用户使用
	if user.Email != req.Email {
		existingUserByEmail, err := userRepo.GetUserByEmail(req.Email)
		if err != nil && err.Error() != apperror.RecordNotFound {
			return nil, apperror.Wrap(err, 500, apperror.DBQueryError)
		}
		if existingUserByEmail != nil && existingUserByEmail.ID != id {
			return nil, apperror.New(409, apperror.EmailExists)
		}
	}

	user.Username = req.Username
	user.NickName = req.NickName
	user.Email = req.Email

	if err := userRepo.Update(user); err != nil {
		slog.Error("更新用户失败", "id", id, "error", err)
		return nil, apperror.Wrap(err, 500, apperror.UserUpdateFailed)
	}

	slog.Info("更新用户成功", "id", id, "name", req.Username)
	return user, nil
}

// DeleteUser 删除用户
func DeleteUser(id int) error {
	if err := userRepo.Delete(id); err != nil {
		slog.Error("删除用户失败", "id", id, "error", err)
		return apperror.Wrap(err, 500, apperror.UserDeleteFailed)
	}
	slog.Info("删除用户成功", "id", id)
	return nil
}
