package events

import (
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
	exp := xp.CalcExperience(calcPower(content))

	author := event.Author

	copaing := sql.Copaing{UserID: author.ID, GuildID: event.GuildID}
	result := sql.DB.FirstOrCreate(&copaing, copaing)
	if result.Error != nil {
		utils.SendAlert(result.Error.Error())
		return
	}
	oldLvl := xp.CalcLevel(copaing.XP)
	copaing.XP += exp
	if oldLvl != xp.CalcLevel(copaing.XP) {
		err := client.MessageReactionAdd(event.ChannelID, event.Message.ID, "⬆")
		if err != nil {
			utils.SendAlert(err.Error())
		}
	}
	result = sql.DB.Save(&copaing)
	if result.Error != nil {
		utils.SendAlert(result.Error.Error())
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
