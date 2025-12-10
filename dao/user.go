package dao

import (
	"errors"

	"gojet/models"
	"gojet/util/apperror"
	"gojet/util/response"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB // GORM 数据库连接实例
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建用户
func (r *UserRepository) Create(user *models.User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		return apperror.Wrap(result.Error, 500, response.MsgDBInsertError)
	}
	return nil
}

// CreateBatch 批量创建用户
func (r *UserRepository) CreateBatch(users []*models.User) error {
	result := r.db.CreateInBatches(users, len(users))
	if result.Error != nil {
		return apperror.Wrap(result.Error, 500, response.MsgDBInsertError)
	}
	return nil
}

// GetAll 获取所有用户
func (r *UserRepository) GetAll() ([]*models.User, error) {
	var users []*models.User
	// GORM 默认不会查询软删除的记录
	result := r.db.Find(&users)
	if result.Error != nil {
		return nil, apperror.Wrap(result.Error, 500, response.MsgDBQueryError)
	}
	return users, nil
}

// GetByID 根据 ID 获取用户
func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, apperror.New(404, response.MsgRecordNotFound)
	}
	if result.Error != nil {
		return nil, apperror.Wrap(result.Error, 500, response.MsgDBQueryError)
	}
	return &user, nil
}

// GetUserByUserName 根据用户名获取用户
func (r *UserRepository) GetUserByUserName(username string) (*models.User, error) {
	var user models.User
	result := r.db.Where("username = ?", username).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, apperror.New(404, response.MsgRecordNotFound)
	}
	if result.Error != nil {
		return nil, apperror.Wrap(result.Error, 500, response.MsgDBQueryError)
	}
	return &user, nil
}

// Update 更新用户 - 保存用户信息到数据库
func (r *UserRepository) Update(user *models.User) error {
	result := r.db.Save(user)
	if result.Error != nil {
		return apperror.Wrap(result.Error, 500, response.MsgDBUpdateError)
	}
	return nil
}

// Delete 删除用户 - 软删除指定 ID 的用户
func (r *UserRepository) Delete(id uint) error {
	result := r.db.Delete(&models.User{}, id)
	if result.Error != nil {
		return apperror.Wrap(result.Error, 500, response.MsgDBDeleteError)
	}
	return nil
}
