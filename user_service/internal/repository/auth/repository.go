package auth

import "github.com/go-redis/redis"

type redisStorage struct {
	client *redis.Client
}