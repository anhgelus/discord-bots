package config

import (
	"github.com/anhgelus/discord-bots/rp-helper/src/redis"
	"github.com/anhgelus/discord-bots/rp-helper/src/utils"
	"github.com/pelletier/go-toml/v2"
	"os"
)

type Config struct {
	Redis redis.RedisCredentials
}

func Get(cfg *Config, defaultConfig string) error {
	c, err := os.ReadFile("/config/config.toml")
	if err != nil {
		utils.SendAlert("main.go - Create file", "Error during reading the file, creating a new one.")
		err = os.WriteFile("/config/config.toml", []byte(defaultConfig), 0666)
		if err != nil {
			return err
		}
		return nil
	}
	return toml.Unmarshal(c, &cfg)
}
