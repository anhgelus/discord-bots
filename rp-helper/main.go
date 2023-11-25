package main

import (
	_ "embed"
	"github.com/anhgelus/discord-bots/rp-helper/src/config"
	start "github.com/anhgelus/discord-bots/rp-helper/src/init"
	"github.com/anhgelus/discord-bots/rp-helper/src/redis"
	"github.com/anhgelus/discord-bots/rp-helper/src/utils"
	"os"
)

//go:embed resources/config.toml
var defaultConfig string

//go:embed resources/objectives.toml
var defaultObjectives string

func main() {
	var cfg config.Config
	err := config.Get(&cfg, defaultConfig)
	if err != nil {
		utils.SendDebug("error")
		utils.SendError(err)
		return
	}

	utils.SendDebug(cfg.Redis.Address, cfg.Main.Debug)

	start.Debug = cfg.Main.Debug

	client, err := cfg.Redis.Get()
	if err != nil {
		utils.SendError(err)
		return
	}
	client.Close()
	redis.Credentials = cfg.Redis

	err = config.GetObjectives(&config.Objs, defaultObjectives)
	if err != nil {
		utils.SendError(err)
		return
	}

	start.Bot(os.Args[1])
}
