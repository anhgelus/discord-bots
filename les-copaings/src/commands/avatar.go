package handlers

import (
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/bwmarrin/discordgo"
)

func Avatar(client *discordgo.Session, i *discordgo.InteractionCreate) {
	err := client.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:  i.Member.User.Username + "'s  avatar",
					Color:  15418782,
					Image:  &discordgo.MessageEmbedImage{URL: i.Member.User.AvatarURL("")},
					Footer: &discordgo.MessageEmbedFooter{Text: i.Member.User.Username, IconURL: i.Member.User.AvatarURL("")},
				},
			},
		},
	})
	if err != nil {
		utils.SendAlert(err.Error())
	}
}
