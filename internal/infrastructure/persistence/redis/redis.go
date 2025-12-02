package redis

import (
	"context"
	"fmt"

	"github.com/HiroLiang/goat-server/internal/config"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() error {
	if &config.App().Redis.Addr == nil || config.App().Redis.Addr == "" {
		fmt.Println("without using redis")
		return nil
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.App().Redis.Addr,
		Password: config.App().Redis.Password,
		DB:       config.App().Redis.DB,
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("failed to connect Redis: %v", err)
	}

	fmt.Println("Connected to Redis")
	return nil
}
