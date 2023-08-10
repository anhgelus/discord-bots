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
	time := user.TimeSinceLastEvent()
	reduce := xp.CalcXpLose(utils.HoursOfUnix(time))
	user.UpdateLastEvent()
	user.Disconnect()
	exp := xp.CalcExperienceFromVocal(user.TimeConnected)

	copaing := sql.Copaing{UserID: event.UserID, GuildID: event.GuildID}
	result := sql.DB.Where("user_id = ? AND guild_id = ?", copaing.UserID, copaing.GuildID).FirstOrCreate(&copaing, copaing)
	if result.Error != nil {
		utils.SendAlert("vocal.go - Querying/Creating copaing", result.Error.Error())
		return
	}
	oldLvl := xp.CalcLevel(copaing.XP)
	if int(copaing.XP)-int(reduce) < 0 {
		copaing.XP = 0
	} else {
		copaing.XP -= reduce
	}
	copaing.XP += exp
	if oldLvl != xp.CalcLevel(copaing.XP) {
		xp.UpdateRolesNoMessage(&copaing, client)
	}
	result = sql.DB.Save(&copaing)
	if result.Error != nil {
		utils.SendAlert("vocal.go - Save copaing", result.Error.Error())
		return
	}
}
