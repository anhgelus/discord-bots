package commands

import (
	"github.com/anhgelus/discord-bots/rp-helper/src/redis"
	"github.com/anhgelus/discord-bots/rp-helper/src/utils"
	"github.com/bwmarrin/discordgo"
)

func Add(client *discordgo.Session, i *discordgo.InteractionCreate) {
	resp := responseBuilder{I: i, C: client}
	optMap := generateOptionMap(i)
	opt, ok := optMap["user"]
	var err error
	if !ok {
		err = resp.Message("Member not given").IsEphemeral().Send()
		if err != nil {
			utils.SendAlert("add.go - No member reply", err.Error())
		}
		return
	}
	u := opt.UserValue(client)
	if u.Bot {
		err = resp.Message("Impossible to add a bot!").IsEphemeral().Send()
		if err != nil {
			utils.SendAlert("add.go - Bot given reply", err.Error())
		}
		return
	}
	p := redis.Player{
		DiscordID: u.ID,
		GuildID:   i.GuildID,
	}
	err = p.Load()
	if err != nil {
		err = resp.Message("Error while loading the player").IsEphemeral().Send()
		if err != nil {
			utils.SendAlert("add.go - Loading player", err.Error())
		}
		return
	}
	p.Save()
	err = resp.Message("Player saved!").IsEphemeral().Send()
	if err != nil {
		utils.SendAlert("add.go - Player saved reply", err.Error())
	}
}
