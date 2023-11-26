package commands

import (
	"fmt"
	"github.com/anhgelus/discord-bots/rp-helper/src/config"
	"github.com/anhgelus/discord-bots/rp-helper/src/redis"
	"github.com/anhgelus/discord-bots/rp-helper/src/utils"
	"github.com/bwmarrin/discordgo"
)

func Profile(client *discordgo.Session, i *discordgo.InteractionCreate) {
	p := redis.Player{
		DiscordID: i.Member.User.ID,
		GuildID:   i.GuildID,
	}
	err := p.Load()
	resp := responseBuilder{
		C: client,
		I: i,
	}
	if err != nil {
		utils.SendAlert("profile.go - Loading player", err.Error())
		err = resp.IsEphemeral().Message("Error, please report this bug").Send()
		if err != nil {
			utils.SendAlert("profile.go - Sending error interaction 1", err.Error())
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
				Name:   fmt.Sprintf("Secondary %d", i+1),
				Value:  s,
				Inline: true,
			}
			if s == "" {
				f.Value = "Error"
			}
			fields = append(fields, f)
		}
	}
	message := ""
	if len(fields) == 0 {
		message = "You have no goals :("
	}
	err = resp.IsEphemeral().Embeds([]*discordgo.MessageEmbed{
		{
			Title:       "Your profile",
			Description: message,
			Color:       utils.Success,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: i.Member.User.AvatarURL("512"),
			},
			Fields: fields,
		},
	}).Send()
	if err != nil {
		utils.SendAlert("profile.go - Sending information", err.Error())
		err = resp.IsEphemeral().Message("Error, please report this bug").Send()
		if err != nil {
			utils.SendAlert("profile.go - Sending error interaction 2", err.Error())
		}
	}
}
