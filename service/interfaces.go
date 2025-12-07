package service

import (
	"gojet/models"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	Create(user *models.User) error
	CreateBatch(users []*models.User) error
	GetAll() ([]*models.User, error)
	GetByID(id uint) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
}
