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
	cfg := sql.Config{GuildID: event.GuildID}
	sql.LoadConfig(&cfg)
	if utils.AStringContains(cfg.DisabledXpChannelsSlice(), event.ChannelID) {
		return
	}
	content := utils.TrimMessage(event.Message.Content)
	event.Member.User = event.Author
	exp := xp.CalcExperience(calcPower(content))

	copaing := sql.GetCopaing(event.Author.ID, event.GuildID)
	data := xp.NewXp(event.Member, &copaing, exp, true)
	if data.IsNewLevel {
		if data.LevelUp {
			err := client.MessageReactionAdd(event.ChannelID, event.Message.ID, "⬆")
			if err != nil {
				utils.SendAlert("message.go - Reaction add", err.Error())
			}
		}
		xp.UpdateRoles(&copaing, client, event)
	}
	sql.Save(&copaing)
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
