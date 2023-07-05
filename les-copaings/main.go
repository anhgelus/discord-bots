package main

import (
	_ "embed"
	"fmt"
	"github.com/anhgelus/discord-bots/les-copaings/src/db/redis"
	"github.com/anhgelus/discord-bots/les-copaings/src/db/sql"
	start "github.com/anhgelus/discord-bots/les-copaings/src/init"
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/pelletier/go-toml/v2"
	"os"
)

//go:embed resources/config.toml
var defaultConfig string

type Config struct {
	SQL       sql.DBCredentials
	Redis     redis.RedisCredentials
	ResetEach uint
}

func main() {
	c, err := os.ReadFile("/config/config.toml")
	if err != nil {
		utils.SendAlert("main.go - Create file", "Error during reading the file, creating a new one.")
		err = os.WriteFile("/config/config.toml", []byte(defaultConfig), 0666)
		if err != nil {
			utils.SendError(err)
			return
		}
		return
	}
	var cfg Config
	err = toml.Unmarshal(c, &cfg)
	if err != nil {
		utils.SendError(err)
		return
	}
	sql.DB = cfg.SQL.Connect()
	if sql.DB == nil {
		utils.SendError(fmt.Errorf("the database is nil"))
		return
	}
	err = sql.DB.AutoMigrate(&sql.Copaing{}, &sql.Config{}, &sql.XpRole{})
	if err != nil {
		utils.SendError(err)
		return
	}
	client, err := cfg.Redis.GetClient()
	if err != nil {
		utils.SendError(err)
		return
	}
	client.Close()
	redis.Credentials = cfg.Redis
	start.Bot(os.Args[1], cfg.ResetEach)
}
