package events

import (
	"github.com/anhgelus/discord-bots/les-copaings/src/db/redis"
	"github.com/anhgelus/discord-bots/les-copaings/src/db/sql"
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/anhgelus/discord-bots/les-copaings/src/xp"
	"github.com/bwmarrin/discordgo"
)

func ConnectionVocal(client *discordgo.Session, event *discordgo.VoiceStateUpdate) {
	if event.BeforeUpdate != nil {
		return
	}
	if event.Member.User.Bot {
		return
	}
	cfg := sql.Config{GuildID: event.GuildID}
	sql.LoadConfig(&cfg)
	if utils.AStringContains(cfg.DisabledXpChannelsSlice(), event.ChannelID) {
		return
	}
	user := redis.GenerateConnectedUser(event.Member)
	user.Connect()
}

func DisconnectionVocal(client *discordgo.Session, event *discordgo.VoiceStateUpdate) {
	if event.ChannelID != "" {
		return
	}
	if event.Member.User.Bot {
		return
	}
	cfg := sql.Config{GuildID: event.GuildID}
	sql.LoadConfig(&cfg)
	if utils.AStringContains(cfg.DisabledXpChannelsSlice(), event.ChannelID) {
		return
	}
	user := redis.GenerateConnectedUser(event.Member)
	user.Disconnect()
	exp := xp.CalcExperienceFromVocal(user.TimeConnected)

	copaing := sql.GetCopaing(event.UserID, event.GuildID)
	data := xp.NewXp(event.Member, &copaing, exp, true)
	if data.IsNewLevel {
		xp.UpdateRolesNoMessage(&copaing, client)
	}
	sql.Save(&copaing)
}
