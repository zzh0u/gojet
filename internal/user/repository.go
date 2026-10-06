package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gojet/internal/infra/apperror"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// UserCacheConfig 用户缓存配置。
var (
	// UserCacheDuration 缓存过期时间。
	UserCacheDuration = 10 * time.Minute

	// UserCachePrefix 缓存键前缀。
	UserCachePrefix = "user:"
)

// UserRepository 用户仓库。
type UserRepository struct {
	db    *gorm.DB
	cache *redis.Client
}

// NewUserRepository 创建用户仓库实例。cache 可为 nil，表示不启用缓存。
func NewUserRepository(db *gorm.DB, cache *redis.Client) *UserRepository {
	return &UserRepository{db: db, cache: cache}
}

func userCacheKey(id int) string {
	return fmt.Sprintf("%s%d", UserCachePrefix, id)
}

// Create 创建用户。
func (r *UserRepository) Create(user *User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		return apperror.Wrap(result.Error, 500, apperror.DBInsertError)
	}
	return nil
}

// CreateBatch 批量创建用户。
func (r *UserRepository) CreateBatch(users []*User) error {
	result := r.db.CreateInBatches(users, len(users))
	if result.Error != nil {
		return apperror.Wrap(result.Error, 500, apperror.DBInsertError)
	}
	return nil
}

// GetAll 获取所有用户。
func (r *UserRepository) GetAll() ([]*User, error) {
	var users []*User
	result := r.db.Find(&users)
	if result.Error != nil {
		return nil, apperror.Wrap(result.Error, 500, apperror.DBQueryError)
	}
	return users, nil
}

// GetByID 根据 ID 获取用户，优先读取缓存。
func (r *UserRepository) GetByID(id int) (*User, error) {
	ctx := context.Background()

	if r.cache != nil {
		key := userCacheKey(id)
		cached, err := r.cache.Get(ctx, key).Result()
		if err == nil {
			var user User
			if err := json.Unmarshal([]byte(cached), &user); err == nil {
				return &user, nil
			}
		}
	}

	user, err := r.getByIDFromDB(id)
	if err != nil {
		return nil, err
	}

	if r.cache != nil && user != nil {
		data, _ := json.Marshal(user)
		r.cache.Set(ctx, userCacheKey(id), data, UserCacheDuration)
	}

	return user, nil
}

func (r *UserRepository) getByIDFromDB(id int) (*User, error) {
	var user User
	result := r.db.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, apperror.New(404, apperror.RecordNotFound)
	}
	if result.Error != nil {
		return nil, apperror.Wrap(result.Error, 500, apperror.DBQueryError)
	}
	return &user, nil
}

// GetUserByUserName 根据用户名获取用户，优先读取缓存。
func (r *UserRepository) GetUserByUserName(username string) (*User, error) {
	ctx := context.Background()

	if r.cache != nil {
		key := fmt.Sprintf("%s%s", UserCachePrefix, "username:"+username)
		cached, err := r.cache.Get(ctx, key).Result()
		if err == nil {
			var user User
			if err := json.Unmarshal([]byte(cached), &user); err == nil {
				return &user, nil
			}
		}
	}

	var user User
	result := r.db.Where("username = ?", username).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, apperror.New(404, apperror.RecordNotFound)
	}
	if result.Error != nil {
		return nil, apperror.Wrap(result.Error, 500, apperror.DBQueryError)
	}

	if r.cache != nil {
		data, _ := json.Marshal(&user)
		r.cache.Set(ctx, userCacheKey(user.ID), data, UserCacheDuration)
		key := fmt.Sprintf("%s%s", UserCachePrefix, "username:"+username)
		r.cache.Set(ctx, key, data, UserCacheDuration)
	}

	return &user, nil
}

// Update 更新用户，并清除对应缓存。
func (r *UserRepository) Update(user *User) error {
	result := r.db.Save(user)
	if result.Error != nil {
		return apperror.Wrap(result.Error, 500, apperror.DBUpdateError)
	}

	if r.cache != nil {
		r.invalidateCache(user.ID, user.Username)
	}

	return nil
}

// Delete 软删除指定 ID 的用户，并清除对应缓存。
func (r *UserRepository) Delete(id int) error {
	var user User
	if err := r.db.First(&user, id).Error; err == nil {
		result := r.db.Delete(&User{}, id)
		if result.Error != nil {
			return apperror.Wrap(result.Error, 500, apperror.DBDeleteError)
		}
		if r.cache != nil {
			r.invalidateCache(user.ID, user.Username)
		}
		return nil
	}

	result := r.db.Delete(&User{}, id)
	if result.Error != nil {
		return apperror.Wrap(result.Error, 500, apperror.DBDeleteError)
	}

	if r.cache != nil {
		r.invalidateCache(id, "")
	}

	return nil
}

// GetUserByEmail 根据邮箱获取用户，优先读取缓存。
func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	ctx := context.Background()

	if r.cache != nil {
		key := fmt.Sprintf("%s%s", UserCachePrefix, "email:"+email)
		cached, err := r.cache.Get(ctx, key).Result()
		if err == nil {
			var user User
			if err = json.Unmarshal([]byte(cached), &user); err == nil {
				return &user, nil
			}
		}
	}

	var user User
	result := r.db.Where("email = ?", email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, apperror.New(404, apperror.RecordNotFound)
	}
	if result.Error != nil {
		return nil, apperror.Wrap(result.Error, 500, apperror.DBQueryError)
	}

	if r.cache != nil {
		data, _ := json.Marshal(&user)
		r.cache.Set(ctx, userCacheKey(user.ID), data, UserCacheDuration)
		key := fmt.Sprintf("%s%s", UserCachePrefix, "email:"+email)
		r.cache.Set(ctx, key, data, UserCacheDuration)
	}

	return &user, nil
}

func (r *UserRepository) invalidateCache(id int, username string) {
	ctx := context.Background()
	keys := []string{userCacheKey(id)}
	if username != "" {
		keys = append(keys, fmt.Sprintf("%s%s", UserCachePrefix, "username:"+username))
	}
	r.cache.Del(ctx, keys...)
}
