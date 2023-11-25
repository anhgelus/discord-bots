package commands

import (
	"fmt"
	"github.com/anhgelus/discord-bots/rp-helper/src/config"
	"github.com/anhgelus/discord-bots/rp-helper/src/redis"
	"github.com/anhgelus/discord-bots/rp-helper/src/utils"
	"github.com/bwmarrin/discordgo"
)

func Me(client *discordgo.Session, i *discordgo.InteractionCreate) {
	p := redis.Player{
		DiscordID: i.Member.User.ID,
		GuildID:   i.GuildID,
	}
	err := p.Load()
	resp := responseBuilder{}
	if err != nil {
		utils.SendAlert("me.go - Loading player", err.Error())
		err = resp.IsEphemeral().Message("Error, please report this bug").Send(client, i)
		if err != nil {
			utils.SendAlert("me.go - Sending error interaction 1", err.Error())
		}
		return
	}
	var fields []*discordgo.MessageEmbedField
	if p.Goals.Main != config.UnsetGoal {
		f := &discordgo.MessageEmbedField{
			Name:   "Main",
			Value:  p.Goals.Main,
			Inline: false,
		}
		if p.Goals.Main == "" {
			f.Value = "Error"
		}
		fields = append(fields, f)
	}
	for i, s := range p.Goals.Secondaries {
		if s != config.UnsetGoal {
			f := &discordgo.MessageEmbedField{
				Name:   fmt.Sprintf("Secondary %d", i),
				Value:  s,
				Inline: false,
			}
			if s == "" {
				f.Value = "Error"
			}
			fields = append(fields, f)
		}
	}
	err = resp.IsEphemeral().Embeds([]*discordgo.MessageEmbed{
		{
			Title:  i.Member.User.Username + " profile",
			Fields: fields,
		},
	}).Send(client, i)
	if err != nil {
		utils.SendAlert("me.go - Sending information", err.Error())
		err = resp.IsEphemeral().Message("Error, please report this bug").Send(client, i)
		if err != nil {
			utils.SendAlert("me.go - Sending error interaction 2", err.Error())
		}
	}
}
