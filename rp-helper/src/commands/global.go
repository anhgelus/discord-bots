package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

type responseBuilder struct {
	content       string
	ephemeral     bool
	deferred      bool
	messageEmbeds []*discordgo.MessageEmbed
	I             *discordgo.InteractionCreate
	C             *discordgo.Session
}

func (res *responseBuilder) Send() error {
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
	return res.C.InteractionRespond(res.I.Interaction, r)
}

func (res *responseBuilder) Interaction(i *discordgo.InteractionCreate) *responseBuilder {
	res.I = i
	return res
}

func (res *responseBuilder) Client(c *discordgo.Session) *responseBuilder {
	res.C = c
	return res
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
	t := time.Now()
	txt := fmt.Sprintf("%d/%d/%d", t.Day(), t.Month(), t.Year())
	av := res.I.Member.User.AvatarURL("")
	for _, em := range e {
		em.Footer = &discordgo.MessageEmbedFooter{
			Text:    txt,
			IconURL: av,
		}
	}
	res.messageEmbeds = e
	return res
}
