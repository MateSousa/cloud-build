package initializers

import (
	"fmt"
	"log"
	"github.com/go-redis/redis"
)

var Redis *redis.Client

func ConnectRedis(config *Config) {
	Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword,
		DB:       0,
	})

	_, err := Redis.Ping().Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis")
	}
	fmt.Println("? Connected Successfully to Redis")
}