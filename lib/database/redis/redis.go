package redis

import (
	log "github.com/mafazans/kumparan/lib/logger"

	"github.com/go-redis/redis"
)

func InitRedis(logger log.Logger) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, _ := client.Ping().Result()

	logger.Info("Redis Initialized: ", pong)

	return client
}
