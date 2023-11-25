package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/anhgelus/discord-bots/rp-helper/src/config"
	"github.com/anhgelus/discord-bots/rp-helper/src/utils"
	"github.com/redis/go-redis/v9"
)

type RedisCredentials struct {
	Address  string
	Password string
	DB       int
}

var Credentials RedisCredentials

var Ctx = context.Background()

type Player struct {
	DiscordID string
	GuildID   string
	Goals     PlayerGoal
}

type PlayerGoal struct {
	Main        string
	Secondaries []string
}

func (rc *RedisCredentials) Get() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     rc.Address,
		Password: rc.Password,
		DB:       rc.DB,
	})
	err := client.Ping(Ctx).Err()
	return client, err
}

func (p *Player) GenKey() string {
	return fmt.Sprintf("%s:%s", p.GuildID, p.DiscordID)
}

func (p *Player) Save() {
	if len(p.Goals.Secondaries) != config.Objs.NumberOfSecondaries {
		utils.SendAlert("redis.go - Saving player", "too much secondaries")
		return
	}
	c, _ := Credentials.Get()
	defer c.Close()
	key := p.GenKey()
	err := c.Set(Ctx, key+":main", p.Goals.Main, 0).Err()
	if err != nil {
		utils.SendAlert("redis.go - Saving player main", err.Error())
	}
	for i := 1; i <= config.Objs.NumberOfSecondaries; i++ {
		err = c.Set(Ctx, fmt.Sprintf("%s:sec%d", key, i), p.Goals.Secondaries[i-1], 0).Err()
		if err != nil {
			utils.SendAlert(fmt.Sprintf("redis.go - Saving player sec%d", i), err.Error())
		}
	}
}

func (p *Player) Load() error {
	if p.GuildID == "" || p.DiscordID == "" {
		return errors.New("guild_id and discord_id not informed")
	}
	c, _ := Credentials.Get()
	defer c.Close()
	key := p.GenKey()
	var err error
	p.Goals.Main, err = c.Get(Ctx, key+":main").Result()
	if errors.Is(err, redis.Nil) {
		p.Goals.Main = config.UnsetGoal
	} else if err != nil {
		utils.SendAlert("redis.go - Loading player main", err.Error())
	}
	var secondaries []string
	for i := 1; i <= config.Objs.NumberOfSecondaries; i++ {
		v, err := c.Get(Ctx, fmt.Sprintf("%s:sec%d", key, i)).Result()
		if errors.Is(err, redis.Nil) {
			v = config.UnsetGoal
		} else if err != nil {
			utils.SendAlert(fmt.Sprintf("redis.go - Loading player sec%d", i), err.Error())
		}
		secondaries = append(secondaries, v)
	}
	p.Goals.Secondaries = secondaries
	return nil
}
