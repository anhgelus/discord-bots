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

func GuildMemberLeave(s *discordgo.Session, event *discordgo.GuildMemberRemove) {
	copaing := sql.Copaing{GuildID: event.GuildID, UserID: event.User.ID}
	sql.DB.FirstOrCreate(&copaing)
	copaing.OldXP = copaing.XP
	copaing.XP = 0
	sql.Save(&copaing)
}

func GuildMemberJoin(s *discordgo.Session, event *discordgo.GuildMemberAdd) {
	copaing := sql.Copaing{GuildID: event.GuildID, UserID: event.User.ID}
	sql.DB.FirstOrCreate(&copaing)
	copaing.XP = copaing.OldXP
	sql.Save(&copaing)
}
