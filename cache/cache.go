package cache

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v9"
)

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, value interface{}) error
	IsExist(key string) (bool, error)
}

type cache struct {
	ctx    context.Context
	client redis.Client
}

func NewCache(ctx context.Context) *cache {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: "",
		DB:       0,
	})
	return &cache{ctx, *rdb}
}

func (c *cache) Get(key string) ([]byte, error) {
	val, err := c.client.Get(c.ctx, key).Bytes()
	if err != nil {
		return nil, fmt.Errorf("failed to get key: %w", err)
	} else if err == redis.Nil {
		fmt.Println("key does not exist")
		return nil, fmt.Errorf("key does not exist")
	}
	return val, nil
}

func (c *cache) Set(key string, value interface{}) error {
	fmt.Println("set cache: ", key)
	err := c.client.Set(c.ctx, key, value, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to set key: %w", err)
	}
	return nil
}

func (c *cache) IsExist(key string) (bool, error) {
	val, err := c.client.Exists(c.ctx, key).Result()
	if val > 0 {
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check existance: %w", err)
	}
	return false, nil
}
