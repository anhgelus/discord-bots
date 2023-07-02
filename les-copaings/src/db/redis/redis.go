package redis

import "github.com/redis/go-redis/v9"

type RedisCredentials struct {
	Address  string
	Password string
	DB       int
}

var Credentials RedisCredentials

func (rc *RedisCredentials) GetClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     rc.Address,
		Password: rc.Password,
		DB:       rc.DB,
	})
}
