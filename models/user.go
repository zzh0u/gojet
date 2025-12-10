package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`                           // 用户ID
	Username  string    `json:"username" binding:"required"`  // 用户登录名称
	NickName  string    `json:"nick_name" binding:"required"` // 用户全名
	Password  string    `json:"password" binding:"required"`  // 用户登录密码
	Email     string    `json:"email" binding:"required"`     // 用户电子邮箱
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
}

func (*User) TableName() string {
	return "user"
}

// CompareSimple 使用 bcrypt 验证密码
func (u *User) CompareSimple(password string) bool {
	// 使用 bcrypt 比较密码
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// HashPassword 使用 bcrypt 生成密码哈希
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
