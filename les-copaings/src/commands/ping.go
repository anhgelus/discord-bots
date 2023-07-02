package commands

import (
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/bwmarrin/discordgo"
)

func Ping(client *discordgo.Session, i *discordgo.InteractionCreate) {
	err := respondInteraction(client, i, "Pong !")
	if err != nil {
		utils.SendAlert(err.Error())
	}
}
