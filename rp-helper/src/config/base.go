package config

import (
	"fmt"
	"github.com/anhgelus/discord-bots/rp-helper/src/redis"
	"github.com/anhgelus/discord-bots/rp-helper/src/utils"
	"github.com/pelletier/go-toml/v2"
	"os"
)

type Config struct {
	Redis redis.RedisCredentials
}

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
