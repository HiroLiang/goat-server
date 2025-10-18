package database

import (
	"context"
	"fmt"

	"github.com/HiroLiang/goat-server/internal/config"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() error {
	if &config.Cfg.Redis.Addr == nil || config.Cfg.Redis.Addr == "" {
		fmt.Println("without using redis")
		return nil
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.Cfg.Redis.Addr,
		Password: config.Cfg.Redis.Password,
		DB:       config.Cfg.Redis.DB,
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("failed to connect Redis: %v", err)
	}

	fmt.Println("Connected to Redis")
	return nil
}
