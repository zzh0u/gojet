//revive:disable var-naming
package api

import (
	"gojet/models"
)

// User 用户服务接口 - 定义 API 层需要的用户操作
// 这个接口让 API 层不直接依赖具体的服务实现，方便测试和解耦
type User interface {
	CreateUser(name string) (*models.User, error)
	CreateInitialData() error
	GetAllUsers() ([]*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	UpdateUser(id uint, name string) (*models.User, error)
	DeleteUser(id uint) error
}
