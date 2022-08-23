package cache

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
)

type Cache interface {
	Get() string
	Set(data string) error
}

type cache struct {
	ctx    context.Context
	client redis.Client
}

func NewCache() *cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // export to ENV
		Password: "",
		DB:       0,
	})
	return &cache{context.Background(), *rdb}
}

func (c *cache) Get(key string, dest interface{}) error {
	// p, err := c.client.HSet()
}

func (c *cache) Set(data string) error {
	fmt.Println("set cache: ", data)
	return nil
}
