package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisCredentials struct {
	Address  string
	Password string
	DB       int
}

var Credentials RedisCredentials

var Ctx = context.Background()

func (rc *RedisCredentials) GetClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     rc.Address,
		Password: rc.Password,
		DB:       rc.DB,
	})
	err := client.Ping(Ctx).Err()
	return client, err
}
