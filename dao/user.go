package dao

import (
	"errors"

	"gojet/models"
	"gojet/service"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB // GORM 数据库连接实例
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *gorm.DB) service.UserRepository {
	return &userRepository{db: db}
}

// Create 创建用户
func (r *userRepository) Create(user *models.User) error {
	result := r.db.Create(user)
	return result.Error
}

// CreateBatch 批量创建用户
func (r *userRepository) CreateBatch(users []*models.User) error {
	result := r.db.CreateInBatches(users, len(users))
	return result.Error
}

// GetAll 获取所有用户
func (r *userRepository) GetAll() ([]*models.User, error) {
	var users []*models.User
	// GORM 默认不会查询软删除的记录
	result := r.db.Find(&users)
	return users, result.Error
}

// GetByID 根据 ID 获取用户 - 查询指定 ID 的用户信息
func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, id)
	// 如果记录不存在，返回 nil 而不是错误
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, result.Error
}

// Update 更新用户 - 保存用户信息到数据库
func (r *userRepository) Update(user *models.User) error {
	result := r.db.Save(user)
	return result.Error
}

// Delete 删除用户 - 软删除指定 ID 的用户
func (r *userRepository) Delete(id uint) error {
	result := r.db.Delete(&models.User{}, id)
	return result.Error
}
