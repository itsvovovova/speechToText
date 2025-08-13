package cache

import (
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisSession struct {
	SessionId string
	Client    *redis.Client
}

type RedisSessionProvider struct {
	Client *redis.Client
}

type RedisSessionManager struct {
	Provider    *RedisSessionProvider
	Cookie      string
	MaxLifetime time.Duration
}
