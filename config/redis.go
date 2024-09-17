package config

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"strconv"
	"time"
)

var RedisClient *redis.Client
var ctx = context.Background()

func InitRedis() {
	db, er := strconv.Atoi(os.Getenv("REDIS_DB"))
	if er != nil {
		log.Fatalf("Could not get cache db: %v", er)
	}
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       db,
	})

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	fmt.Println("Connected to Redis")
}

func SetKey(ctx *gin.Context, key string, value string, time time.Duration) {
	err := RedisClient.Set(ctx, key, value, time).Err()
	if err != nil {
		log.Fatalf("Could not set cache: %v", err)
	}
}

func GetValue(ctx *gin.Context, key string) string {
	val, err := RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		fmt.Println("Key does not exist")
	} else if err != nil {
		log.Fatalf("Could not get cache: %v", err)
	}
	return val
}
