package events

import (
	"fmt"
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

	copaing := sql.Copaing{UserID: event.Author.ID, GuildID: event.GuildID}
	result := sql.DB.FirstOrCreate(&copaing, copaing)
	if result.Error != nil {
		utils.SendAlert("message.go - Querying/Creating copaing", result.Error.Error())
		return
	}
	oldLvl := xp.CalcLevel(copaing.XP)
	copaing.XP += exp
	if oldLvl != xp.CalcLevel(copaing.XP) {
		err := client.MessageReactionAdd(event.ChannelID, event.Message.ID, "⬆")
		if err != nil {
			utils.SendAlert("message.go - Reaction add", err.Error())
		}
		updateRoles(&copaing, client, event)
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

func updateRoles(copaing *sql.Copaing, client *discordgo.Session, event *discordgo.MessageCreate) {
	cfg := sql.Config{GuildID: copaing.GuildID}
	sql.DB.Model(&sql.Config{}).Preload("XpRoles").FirstOrCreate(&cfg)

	roles := make(chan string)

	go sql.GetNewRoles(copaing, &cfg, event.Member.Roles, roles)

	for role := range roles {
		err := client.GuildMemberRoleAdd(copaing.GuildID, copaing.UserID, role)
		if err != nil {
			utils.SendAlert("message.go - Role add", err.Error())
			_, err = client.ChannelMessageSend(event.ChannelID, "Impossible de vous ajouter le rôle "+role)
			if err != nil {
				utils.SendAlert("message.go - Message send role failed", err.Error())
			}
			continue
		}
		_, err = client.ChannelMessageSend(event.ChannelID,
			fmt.Sprintf("<@%s>, vous venez d'obtenir un nouveau rôle !", copaing.UserID),
		)
		if err != nil {
			utils.SendAlert("message.go - New role message", err.Error())
		}
	}
}
