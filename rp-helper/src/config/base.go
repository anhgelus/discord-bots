package config

import (
	"context"
	"fmt"
	"github.com/anhgelus/discord-bots/rp-helper/src/utils"
	"github.com/pelletier/go-toml/v2"
	"github.com/redis/go-redis/v9"
	"os"
)

func Get(cfg *Config, defaultConfig string) error {
	return get(cfg, defaultConfig, "config")
}

func GetObjectives(cfg *Objectives, defaultConfig string) error {
	return get(cfg, defaultConfig, "objectives")
}

func get(cfg any, defaultConfig string, name string) error {
	path := fmt.Sprintf("/config/%s.toml", name)
	c, err := os.ReadFile(path)
	if err != nil {
		utils.SendAlert("base.go - Create file", "Error during reading the file, creating a new one.")
		err = os.WriteFile(path, []byte(defaultConfig), 0666)
		if err != nil {
			return err
		}
		return nil
	}
	return toml.Unmarshal(c, &cfg)
}

type Config struct {
	Redis RedisCredentials
}

type RedisCredentials struct {
	Address  string
	Password string
	DB       int
}

func (rc *RedisCredentials) Get() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     rc.Address,
		Password: rc.Password,
		DB:       rc.DB,
	})
	err := client.Ping(context.Background()).Err()
	return client, err
}
