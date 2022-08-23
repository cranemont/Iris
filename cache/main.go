package cache

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
)

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	err := rdb.Set(ctx, "k", "v", 0).Err()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	val, err := rdb.Get(ctx, "k").Result()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("k", val)

	val2, err := rdb.Get(ctx, "k2").Result()
	if err == redis.Nil {
		fmt.Println("k2 does not exist")
	} else if err != nil {
		fmt.Println(err)
		panic(err)
	} else {
		fmt.Println("k2", val2)
	}
}
