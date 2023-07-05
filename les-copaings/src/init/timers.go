package start

import (
	"github.com/anhgelus/discord-bots/les-copaings/src/db/redis"
	"github.com/anhgelus/discord-bots/les-copaings/src/db/sql"
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/bwmarrin/discordgo"
	rdb "github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

const LastResetKey = "last_reset"

func SetupTimers(interval uint, s *discordgo.Session) {
	ticker := time.NewTicker(1 * time.Hour)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				checkReset(interval, s)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func intervalToUnix(interval uint) int {
	return int(interval * 30 * 24 * 60 * 60)
}

func checkReset(interval uint, s *discordgo.Session) {
	client, _ := redis.Credentials.GetClient()
	val := client.Get(redis.Ctx, LastResetKey)
	if val.Err() == rdb.Nil {
		initialize(s, client)
		return
	}
	last, err := strconv.Atoi(val.Val())
	if err != nil {
		utils.SendAlert("timers.go - Str to Int conversion", err.Error())
		return
	}
	if time.Now().Unix() >= int64(last+intervalToUnix(interval)) {
		reset(s)
	}
}

func initialize(s *discordgo.Session, client *rdb.Client) {
	client.Set(redis.Ctx, LastResetKey, time.Now().Unix(), 0)
}

func reset(s *discordgo.Session) {
	guilds, err := s.UserGuilds(0, "", "")
	if err != nil {
		utils.SendAlert("timers.go - Guilds", err.Error())
		return
	}
	for _, guild := range guilds {
		resetGuild(guild)
	}
}

func resetGuild(guild *discordgo.UserGuild) {
	//TODO: send message to broadcast the reset when config will be implemented
	//TODO: reset roles when roles will be implemented

	// reset the xp of all members
	sql.DB.Model(sql.Copaing{}).Where("guild_id = ?", guild.ID).Updates(sql.Copaing{XP: 0})
}
