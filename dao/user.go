package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gojet/models"
	"gojet/util/apperror"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// UserCacheConfig 用户缓存配置
var (
	// UserCacheDuration 缓存过期时间
	UserCacheDuration = 10 * time.Minute

	// UserCachePrefix 缓存键前缀
	UserCachePrefix = "user:"
)

// UserRepository 用户仓库
type UserRepository struct {
	db    *gorm.DB
	cache *redis.Client // Redis 客户端，用于缓存
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db, cache: nil}
}

// SetCache 设置缓存客户端（用于运行时启用缓存）
func (r *UserRepository) SetCache(cache *redis.Client) {
	r.cache = cache
}

// userCacheKey 生成用户缓存键
func userCacheKey(id int) string {
	return fmt.Sprintf("%s%d", UserCachePrefix, id)
}

// Create 创建用户
func (r *UserRepository) Create(user *models.User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		return apperror.Wrap(result.Error, 500, apperror.DBInsertError)
	}
	return nil
}

// CreateBatch 批量创建用户
func (r *UserRepository) CreateBatch(users []*models.User) error {
	result := r.db.CreateInBatches(users, len(users))
	if result.Error != nil {
		return apperror.Wrap(result.Error, 500, apperror.DBInsertError)
	}
	return nil
}

// GetAll 获取所有用户
func (r *UserRepository) GetAll() ([]*models.User, error) {
	var users []*models.User
	// GORM 默认不会查询软删除的记录
	result := r.db.Find(&users)
	if result.Error != nil {
		return nil, apperror.Wrap(result.Error, 500, apperror.DBQueryError)
	}
	return users, nil
}

// GetByID 根据 ID 获取用户 - 支持缓存
func (r *UserRepository) GetByID(id int) (*models.User, error) {
	ctx := context.Background()

	// 如果有缓存，先从缓存获取
	if r.cache != nil {
		key := userCacheKey(id)
		cached, err := r.cache.Get(ctx, key).Result()
		if err == nil {
			// 缓存命中
			var user models.User
			if err := json.Unmarshal([]byte(cached), &user); err == nil {
				return &user, nil
			}
		}
	}

	// 缓存未命中，从数据库获取
	user, err := r.getByIDFromDB(ctx, id)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if r.cache != nil && user != nil {
		data, _ := json.Marshal(user)
		r.cache.Set(ctx, userCacheKey(id), data, UserCacheDuration)
	}

	return user, nil
}

// getByIDFromDB 从数据库获取用户
func (r *UserRepository) getByIDFromDB(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, apperror.New(404, apperror.RecordNotFound)
	}
	if result.Error != nil {
		return nil, apperror.Wrap(result.Error, 500, apperror.DBQueryError)
	}
	return &user, nil
}

// GetUserByUserName 根据用户名获取用户 - 支持缓存
func (r *UserRepository) GetUserByUserName(username string) (*models.User, error) {
	ctx := context.Background()

	// 如果有缓存，通过用户名缓存键查询
	if r.cache != nil {
		key := fmt.Sprintf("%s%s", UserCachePrefix, "username:"+username)
		cached, err := r.cache.Get(ctx, key).Result()
		if err == nil {
			var user models.User
			if err := json.Unmarshal([]byte(cached), &user); err == nil {
				return &user, nil
			}
		}
	}

	// 缓存未命中，从数据库获取
	var user models.User
	result := r.db.Where("username = ?", username).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, apperror.New(404, apperror.RecordNotFound)
	}
	if result.Error != nil {
		return nil, apperror.Wrap(result.Error, 500, apperror.DBQueryError)
	}

	// 写入缓存
	if r.cache != nil {
		// 按 ID 缓存
		data, _ := json.Marshal(&user)
		r.cache.Set(ctx, userCacheKey(user.ID), data, UserCacheDuration)
		// 按用户名缓存
		key := fmt.Sprintf("%s%s", UserCachePrefix, "username:"+username)
		r.cache.Set(ctx, key, data, UserCacheDuration)
	}

	return &user, nil
}

// Update 更新用户 - 保存用户信息到数据库
func (r *UserRepository) Update(user *models.User) error {
	result := r.db.Save(user)
	if result.Error != nil {
		return apperror.Wrap(result.Error, 500, apperror.DBUpdateError)
	}

	// 清除用户缓存
	if r.cache != nil {
		r.invalidateCache(user.ID, user.Username)
	}

	return nil
}

// Delete 删除用户 - 软删除指定 ID 的用户
func (r *UserRepository) Delete(id int) error {
	// 先获取用户信息，用于清除用户名缓存
	var user models.User
	if err := r.db.First(&user, id).Error; err == nil {
		result := r.db.Delete(&models.User{}, id)
		if result.Error != nil {
			return apperror.Wrap(result.Error, 500, apperror.DBDeleteError)
		}
		// 清除缓存
		if r.cache != nil {
			r.invalidateCache(user.ID, user.Username)
		}
		return nil
	}

	result := r.db.Delete(&models.User{}, id)
	if result.Error != nil {
		return apperror.Wrap(result.Error, 500, apperror.DBDeleteError)
	}

	// 清除缓存
	if r.cache != nil {
		r.invalidateCache(id, "")
	}

	return nil
}

// GetUserByEmail 根据邮箱获取用户
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	ctx := context.Background()

	// 如果有缓存，通过邮箱缓存键查询
	if r.cache != nil {
		key := fmt.Sprintf("%s%s", UserCachePrefix, "email:"+email)
		cached, err := r.cache.Get(ctx, key).Result()
		if err == nil {
			var user models.User
			if err = json.Unmarshal([]byte(cached), &user); err == nil {
				return &user, nil
			}
		}
	}

	// 缓存未命中，从数据库获取
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, apperror.New(404, apperror.RecordNotFound)
	}
	if result.Error != nil {
		return nil, apperror.Wrap(result.Error, 500, apperror.DBQueryError)
	}

	// 写入缓存
	if r.cache != nil {
		data, _ := json.Marshal(&user)
		r.cache.Set(ctx, userCacheKey(user.ID), data, UserCacheDuration)
		key := fmt.Sprintf("%s%s", UserCachePrefix, "email:"+email)
		r.cache.Set(ctx, key, data, UserCacheDuration)
	}

	return &user, nil
}

// invalidateCache 清除用户缓存
func (r *UserRepository) invalidateCache(id int, username string) {
	ctx := context.Background()
	keys := []string{userCacheKey(id)}
	if username != "" {
		keys = append(keys, fmt.Sprintf("%s%s", UserCachePrefix, "username:"+username))
	}
	r.cache.Del(ctx, keys...)
}
