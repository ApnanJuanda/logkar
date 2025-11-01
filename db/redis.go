package db

import (
	"bsnack/lib/env"
	"github.com/redis/go-redis/v9"
	"log"
)

func RedisNewClient() *redis.Client {
	options, err := redis.ParseURL(env.String("REDIS_URL", "redis://127.0.0.1:6379/0"))
	if err != nil {
		log.Fatal(err)
	}
	return redis.NewClient(options)
}
