package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User 用户数据模型。
type User struct {
	ID        int       `json:"id"`                                             // 用户 ID
	Username  string    `json:"username" binding:"required" gorm:"uniqueIndex"` // 用户登录名称
	NickName  string    `json:"nick_name" binding:"required"`                   // 用户全名
	Password  string    `json:"password" binding:"required"`                    // 用户登录密码
	Email     string    `json:"email" binding:"required" gorm:"uniqueIndex"`    // 用户电子邮箱
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
}

func (*User) TableName() string {
	return "user"
}

// CompareSimple 使用 bcrypt 验证密码。
func (u *User) CompareSimple(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// HashPassword 使用 bcrypt 生成密码哈希。
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CreateUserRequest 创建用户请求结构体。
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	NickName string `json:"nick_name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

// UpdateUserRequest 更新用户请求结构体。
type UpdateUserRequest struct {
	Username string `json:"username" binding:"required"`
	NickName string `json:"nick_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// UserResponse 用户响应结构体。
type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	NickName string `json:"nick_name"`
	Email    string `json:"email"`
}
