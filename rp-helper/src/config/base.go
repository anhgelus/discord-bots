package config

import (
	"context"
	"errors"
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
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		utils.SendAlert("base.go - Create file", "Error during reading the file, creating a new one.")
		c = []byte(defaultConfig)
		err = os.WriteFile(path, c, 0666)
		if err != nil {
			return err
		}
	}
	return toml.Unmarshal(c, &cfg)
}

type Config struct {
	Main  MainSettings
	Redis RedisCredentials
}

type MainSettings struct {
	Debug bool
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
