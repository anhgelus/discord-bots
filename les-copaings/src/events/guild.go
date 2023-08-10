package events

import (
	"fmt"
	"github.com/anhgelus/discord-bots/les-copaings/src/db/redis"
	"github.com/anhgelus/discord-bots/les-copaings/src/db/sql"
	"github.com/bwmarrin/discordgo"
)

func GuildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {

}

func GuildDelete(s *discordgo.Session, event *discordgo.GuildDelete) {
	guildID := event.ID
	client, _ := redis.Credentials.GetClient()
	client.Del(redis.Ctx, fmt.Sprintf("%s:*", guildID))
	sql.DB.Model(sql.Copaing{}).Where("guild_id = ?", guildID).Delete(sql.Copaing{})
	sql.DB.Model(sql.Config{}).Where("guild_id = ?", guildID).Delete(sql.Config{})
	client.Close()
}
