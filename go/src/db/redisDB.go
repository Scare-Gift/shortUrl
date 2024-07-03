package db

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var (
	Rdb  *redis.Client
	Rctx context.Context
)

func init() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "192.168.104:6379",
		Password: "",
		DB:       0,
	})
	Rctx = context.Background()
}
