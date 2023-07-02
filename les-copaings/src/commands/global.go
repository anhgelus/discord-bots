package commands

import "github.com/bwmarrin/discordgo"

func respondInteraction(client *discordgo.Session, i *discordgo.InteractionCreate, msg string) error {
	return client.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
}
