// Package redis Redis 客户端封装。可复用，不依赖具体业务模块。
package redis

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"
)

// Client Redis 客户端封装。
type Client struct {
	client *goredis.Client
}

// NewClient 创建并探测 Redis 连接。
func NewClient(host string, port int, password string, db int) (*Client, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("连接 Redis 失败: %w", err)
	}

	return &Client{client: client}, nil
}

// Raw 返回原始 redis 客户端，供仓储层读写缓存。
func (c *Client) Raw() *goredis.Client {
	return c.client
}

// Close 关闭连接。
func (c *Client) Close() error {
	if c == nil || c.client == nil {
		return nil
	}
	return c.client.Close()
}
