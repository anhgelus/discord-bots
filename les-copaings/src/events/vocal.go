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
	user := redis.GenerateConnectedUser(event.Member)
	user.Disconnect()
	exp := xp.CalcExperienceFromVocal(user.TimeConnected)

	copaing := sql.Copaing{UserID: event.UserID, GuildID: event.GuildID}
	result := sql.DB.FirstOrCreate(&copaing, copaing)
	if result.Error != nil {
		utils.SendAlert(result.Error.Error())
		return
	}
	oldLvl := xp.CalcLevel(copaing.XP)
	copaing.XP += exp
	if oldLvl != xp.CalcLevel(copaing.XP) {
		//TODO: handle level up on vocal
	}
	result = sql.DB.Save(&copaing)
	if result.Error != nil {
		utils.SendAlert(result.Error.Error())
		return
	}
}
