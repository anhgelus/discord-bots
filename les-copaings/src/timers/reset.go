package timers

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

func SetupTimers(s *discordgo.Session, interval uint) {
	utils.NewTimers(1*time.Minute, func(_ chan struct{}) {
		checkReset(s, interval)
	})
}

func intervalToUnix(interval uint) int {
	return int(interval * 30 * 24 * 60 * 60)
}

func checkReset(s *discordgo.Session, interval uint) {
	checkGuilds(s, s.State.Guilds, interval)
}

func checkGuilds(s *discordgo.Session, guilds []*discordgo.Guild, interval uint) {
	client, _ := redis.Credentials.GetClient()
	for _, guild := range guilds {
		val := client.Get(redis.Ctx, GenLastResetKey(guild.ID))
		if val.Err() == rdb.Nil {
			InitializeGuild(guild.ID, client)
			continue
		} else if val.Err() != nil {
			utils.SendAlert("reset.go - Get last key", val.Err().Error())
			continue
		}
		last, err := strconv.Atoi(val.Val())
		if err != nil {
			utils.SendAlert("reset.go - Str to Int conversion", err.Error())
			continue
		}
		if time.Now().Unix() >= int64(last+intervalToUnix(interval)) {
			ResetGuild(s, guild)
		}
	}
}

func GenLastResetKey(guildID string) string {
	return strings.Replace(LastResetKey, "guild_id", guildID, -1)
}

func InitializeGuild(guildID string, client *rdb.Client) {
	client.Set(redis.Ctx, GenLastResetKey(guildID), time.Now().Unix(), 0)
}

func ResetGuild(s *discordgo.Session, guild *discordgo.Guild) {
	// reset the xp of all members
	sql.DB.Model(sql.Copaing{}).Where("guild_id = ?", guild.ID).Updates(sql.Copaing{XP: 0})

	// reset roles
	cfg := sql.Config{GuildID: guild.ID}
	sql.DB.FirstOrCreate(&cfg)
	members := utils.FetchGuildUser(s, guild.ID)
	for _, member := range members {
		for _, role := range member.Roles {
			for _, xpr := range cfg.XpRoles {
				if xpr.Role == role {
					err := s.GuildMemberRoleRemove(guild.ID, member.User.ID, role)
					if err != nil {
						utils.SendAlert("reset.go - Reset roles", err.Error())
						continue
					}
				}
			}
		}
	}

	//TODO: broadcast this news
}
