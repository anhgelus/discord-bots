package handlers

import (
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/bwmarrin/discordgo"
)

func Ping(client *discordgo.Session, i *discordgo.InteractionCreate) {
	err := client.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Response gived everythings seems fine here",
		},
	})
	if err != nil {
		utils.SendAlert(err.Error())
	}
}
