package start

import (
	"github.com/anhgelus/discord-bots/les-copaings/src/db/redis"
	"github.com/anhgelus/discord-bots/les-copaings/src/db/sql"
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/bwmarrin/discordgo"
	rdb "github.com/redis/go-redis/v9"
	"strconv"
	"strings"
	"time"
)

const LastResetKey = "guild_id:last_reset"

func SetupTimers(interval uint, s *discordgo.Session) {
	utils.NewTimers(1*time.Hour, func(_ chan struct{}) {
		checkReset(interval, s)
	})
}

func intervalToUnix(interval uint) int {
	return int(interval * 30 * 24 * 60 * 60)
}

func checkReset(interval uint, s *discordgo.Session) {
	checkGuilds(s.State.Guilds, interval, s)
}

func checkGuilds(guilds []*discordgo.Guild, interval uint, s *discordgo.Session) {
	client, _ := redis.Credentials.GetClient()
	for _, guild := range guilds {
		val := client.Get(redis.Ctx, genLastResetKey(guild.ID))
		if val.Err() == rdb.Nil {
			initialize(guild.ID, client)
			continue
		} else if val.Err() != nil {
			utils.SendAlert("timers.go - Get last key", val.Err().Error())
			continue
		}
		last, err := strconv.Atoi(val.Val())
		if err != nil {
			utils.SendAlert("timers.go - Str to Int conversion", err.Error())
			continue
		}
		if time.Now().Unix() >= int64(last+intervalToUnix(interval)) {
			reset(s)
		}
	}
}

func genLastResetKey(guildID string) string {
	return strings.Replace(LastResetKey, "guild_id", guildID, -1)
}

func initialize(guildID string, client *rdb.Client) {
	client.Set(redis.Ctx, genLastResetKey(guildID), time.Now().Unix(), 0)
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
