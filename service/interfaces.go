package service

import (
	"gojet/models"
)

// User 用户接口
type User interface {
	Create(user *models.User) error
	CreateBatch(users []*models.User) error
	GetAll() ([]*models.User, error)
	GetByID(id uint) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
}
