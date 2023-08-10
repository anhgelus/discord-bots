package events

import (
	"github.com/anhgelus/discord-bots/les-copaings/src/db/redis"
	"github.com/anhgelus/discord-bots/les-copaings/src/db/sql"
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/anhgelus/discord-bots/les-copaings/src/xp"
	"github.com/bwmarrin/discordgo"
)

func MessageSent(client *discordgo.Session, event *discordgo.MessageCreate) {
	if event.Author.Bot {
		return
	}
	content := event.Message.Content
	event.Member.User = event.Author
	user := redis.GenerateConnectedUser(event.Member)
	time := user.TimeSinceLastEvent()
	reduce := xp.CalcXpLose(utils.HoursOfUnix(time))
	user.UpdateLastEvent()
	exp := xp.CalcExperience(calcPower(content))

	copaing := sql.Copaing{UserID: event.Author.ID, GuildID: event.GuildID}
	result := sql.DB.Where("user_id = ? AND guild_id = ?", copaing.UserID, copaing.GuildID).FirstOrCreate(&copaing, copaing)
	if result.Error != nil {
		utils.SendAlert("message.go - Querying/Creating copaing", result.Error.Error())
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
		err := client.MessageReactionAdd(event.ChannelID, event.Message.ID, "⬆")
		if err != nil {
			utils.SendAlert("message.go - Reaction add", err.Error())
		}
		xp.UpdateRoles(&copaing, client, event)
	}
	result = sql.DB.Save(&copaing)
	if result.Error != nil {
		utils.SendAlert("message.go - Save copaing", result.Error.Error())
		return
	}
}

func calcPower(message string) (uint, uint) {
	var chars []rune
	for _, c := range []rune(message) {
		toAdd := true
		for _, ch := range chars {
			if ch == c {
				toAdd = false
			}
		}
		if toAdd {
			chars = append(chars, c)
		}
	}
	return uint(len(message)), uint(len(chars))
}
