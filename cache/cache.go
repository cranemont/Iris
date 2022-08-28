package cache

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v9"
)

type Cache interface {
	Get(key string) []byte
	Set(key string, value interface{})
	IsExist(key string) bool
}

type cache struct {
	ctx    context.Context
	client redis.Client
}

func NewCache(ctx context.Context) *cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // TODO: export to ENV
		Password: "",
		DB:       0,
	})
	return &cache{ctx, *rdb}
}

func (c *cache) Get(key string) []byte {
	val, err := c.client.Get(c.ctx, key).Bytes()
	if err != nil {
		log.Println(err)
	} else if err == redis.Nil {
		fmt.Println("k2 does not exist")
		return nil
	}
	return val
}

func (c *cache) Set(key string, value interface{}) {
	fmt.Println("set cache: ", key)
	err := c.client.Set(c.ctx, key, value, 0).Err()
	if err != nil {
		log.Println(err)
	}
}

func (c *cache) IsExist(key string) bool {
	val, err := c.client.Exists(c.ctx, key).Result()
	if val > 0 {
		return true
	}
	if err != nil {
		log.Println(err)
	}
	return false
}
