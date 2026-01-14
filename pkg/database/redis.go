package database

import (
	"context"
	"fmt"
	"log"
	"shifty-backend/configs"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(cfg *configs.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0,
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Can not connect to Redis: ", err)
	}
	log.Println("Connect Redis successful!")
	return rdb
}
