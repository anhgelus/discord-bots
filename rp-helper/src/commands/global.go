package commands

import "github.com/bwmarrin/discordgo"

type responseBuilder struct {
	content       string
	ephemeral     bool
	deferred      bool
	messageEmbeds []*discordgo.MessageEmbed
}

func (res *responseBuilder) Send(client *discordgo.Session, i *discordgo.InteractionCreate) error {
	r := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: res.content,
			Embeds:  res.messageEmbeds,
		},
	}
	if res.deferred {
		r.Type = discordgo.InteractionResponseDeferredChannelMessageWithSource
	}
	if res.ephemeral {
		r.Data.Flags = discordgo.MessageFlagsEphemeral
	}
	return client.InteractionRespond(i.Interaction, r)
}

func (res *responseBuilder) IsEphemeral() *responseBuilder {
	res.ephemeral = true
	return res
}

func (res *responseBuilder) NotEphemeral() *responseBuilder {
	res.ephemeral = false
	return res
}

func (res *responseBuilder) IsDeferred() *responseBuilder {
	res.deferred = true
	return res
}

func (res *responseBuilder) NotDeferred() *responseBuilder {
	res.deferred = false
	return res
}

func (res *responseBuilder) Message(s string) *responseBuilder {
	res.content = s
	return res
}

func (res *responseBuilder) Embeds(e []*discordgo.MessageEmbed) *responseBuilder {
	res.messageEmbeds = e
	return res
}
