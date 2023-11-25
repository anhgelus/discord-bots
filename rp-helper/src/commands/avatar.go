package commands

import (
	"github.com/anhgelus/discord-bots/rp-helper/src/utils"
	"github.com/bwmarrin/discordgo"
)

func Avatar(client *discordgo.Session, i *discordgo.InteractionCreate) {
	resp := responseBuilder{}
	err := resp.Embeds([]*discordgo.MessageEmbed{
		{
			Title:  i.Member.User.Username + "'s  avatar",
			Color:  utils.Success,
			Image:  &discordgo.MessageEmbedImage{URL: i.Member.User.AvatarURL("")},
			Footer: &discordgo.MessageEmbedFooter{Text: i.Member.User.Username, IconURL: i.Member.User.AvatarURL("")},
		},
	}).Send(client, i)
	if err != nil {
		utils.SendAlert("avatar.go - Interaction respond", err.Error())
	}
}
