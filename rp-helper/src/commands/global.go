package commands

import "github.com/bwmarrin/discordgo"

type responseBuilder struct {
	Content       string
	Ephemeral     bool
	Deferred      bool
	MessageEmbeds []*discordgo.MessageEmbed
}

func (res *responseBuilder) Send(client *discordgo.Session, i *discordgo.InteractionCreate) error {
	r := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: res.Content,
			Embeds:  res.MessageEmbeds,
		},
	}
	if res.Deferred {
		r.Type = discordgo.InteractionResponseDeferredChannelMessageWithSource
	}
	if res.Ephemeral {
		r.Data.Flags = discordgo.MessageFlagsEphemeral
	}
	return client.InteractionRespond(i.Interaction, r)
}

func (res *responseBuilder) IsEphemeral() *responseBuilder {
	res.Ephemeral = true
	return res
}

func (res *responseBuilder) IsDeferred() *responseBuilder {
	res.Deferred = true
	return res
}

func (res *responseBuilder) Message(s string) *responseBuilder {
	res.Content = s
	return res
}

func (res *responseBuilder) Embeds(e []*discordgo.MessageEmbed) *responseBuilder {
	res.MessageEmbeds = e
	return res
}
