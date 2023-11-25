package commands

import (
	"fmt"
	"github.com/anhgelus/discord-bots/rp-helper/src/redis"
	"github.com/anhgelus/discord-bots/rp-helper/src/utils"
	"github.com/bwmarrin/discordgo"
)

func List(client *discordgo.Session, i *discordgo.InteractionCreate) {
	resp := responseBuilder{I: i, C: client}
	ps := redis.GetPlayers(i.GuildID)
	var msg string
	for i, p := range ps {
		msg += fmt.Sprintf("- <@%s>", p.DiscordID)
		if i != len(ps)-1 {
			msg += "\n"
		}
	}
	err := resp.Embeds([]*discordgo.MessageEmbed{
		{
			Title:       "List of players",
			Description: msg,
			Color:       utils.Success,
		},
	}).Send()
	if err != nil {
		utils.SendAlert("list.go - Responding to interaction", err.Error())
	}
}
