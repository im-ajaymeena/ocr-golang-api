package common

import (
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_SERVER_URL"),
	Password: "", // no password set
	DB:       0,  // use default DB
})
