package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/anhgelus/discord-bots/rp-helper/src/config"
	"github.com/anhgelus/discord-bots/rp-helper/src/utils"
	"github.com/redis/go-redis/v9"
	"strings"
)

var Credentials config.RedisCredentials

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

var (
	ErrGuildIDDiscordIDNotPresent = errors.New("guild_id or discord_id not informed")
	ErrWrongSecondaries           = errors.New("wrong number of secondaries")
)

func (p *Player) GenKey() string {
	return fmt.Sprintf("%s:%s", p.GuildID, p.DiscordID)
}

func (p *Player) Save() error {
	if len(p.Goals.Secondaries) != config.Objs.Settings.NumberOfSecondaries {
		return ErrWrongSecondaries
	}
	c, _ := Credentials.Get()
	defer c.Close()
	key := p.GenKey()
	err := c.Set(Ctx, key+":main", p.Goals.Main, 0).Err()
	if err != nil {
		utils.SendAlert("redis.go - Saving player main", err.Error())
	}
	for i := 1; i <= config.Objs.Settings.NumberOfSecondaries; i++ {
		err = c.Set(Ctx, fmt.Sprintf("%s:sec%d", key, i), p.Goals.Secondaries[i-1], 0).Err()
		if err != nil {
			utils.SendAlert(fmt.Sprintf("redis.go - Saving player sec%d", i), err.Error())
		}
	}
	ps, err := c.Get(Ctx, p.GuildID).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}
	sp := strings.Split(ps, ",")
	if !utils.AStringContains(sp, p.DiscordID) {
		sp = append(sp, p.DiscordID)
		ps = strings.Join(sp, ",")
		err = c.Set(Ctx, p.GuildID, ps, 0).Err()
		if err != nil {
			return err
		}
	}
	return nil
}
func (p *Player) Remove() error {
	if len(p.Goals.Secondaries) != config.Objs.Settings.NumberOfSecondaries {
		return ErrWrongSecondaries
	}
	c, _ := Credentials.Get()
	defer c.Close()
	key := p.GenKey()
	err := c.Del(Ctx, key+":main").Err()
	if err != nil {
		utils.SendAlert("redis.go - Removing player main", err.Error())
	}
	for i := 1; i <= config.Objs.Settings.NumberOfSecondaries; i++ {
		err = c.Del(Ctx, fmt.Sprintf("%s:sec%d", key, i)).Err()
		if err != nil {
			utils.SendAlert(fmt.Sprintf("redis.go - Removing player sec%d", i), err.Error())
		}
	}
	ps, err := c.Get(Ctx, p.GuildID).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}
	// Hello, World, Today
	nps := strings.ReplaceAll(ps, ", ", ",") + ","
	// Hello,World,Today,
	nps = strings.Replace(nps, p.DiscordID+",", "", 1)
	// Hello,Today,
	nps, _ = strings.CutSuffix(nps, ",")
	// Hello,Today
	//nps = strings.ReplaceAll(k, ",", ", ")
	// Hello, Today
	err = c.Set(Ctx, p.GuildID, nps, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (p *Player) Load() error {
	if p.GuildID == "" || p.DiscordID == "" {
		return ErrGuildIDDiscordIDNotPresent
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
	for i := 1; i <= config.Objs.Settings.NumberOfSecondaries; i++ {
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

func GetPlayers(guildID string) []Player {
	c, _ := Credentials.Get()
	defer c.Close()
	ps, err := c.Get(Ctx, guildID).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		utils.SendAlert("redis.go - Getting all players", err.Error())
		return nil
	}
	sp := strings.Split(ps, ",")
	var players []Player
	for _, s := range sp {
		if s == "" {
			continue
		}
		p := Player{DiscordID: s, GuildID: guildID}
		err = p.Load()
		if err != nil {
			utils.SendAlert("redis.go - Loading player", err.Error())
			continue
		}
		players = append(players, p)
	}
	return players
}
