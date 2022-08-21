package cache

import "fmt"

type Cache interface {
	Get() string
	Set(data string) error
}

type cache struct {
}

func (c *cache) Get() string {
	return "get from cache"
}

func (c *cache) Set(data string) error {
	fmt.Println("set cache: ", data)
	return nil
}

func NewCache() *cache {
	return &cache{}
}
