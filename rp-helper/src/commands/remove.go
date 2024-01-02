package commands

import (
	"github.com/anhgelus/discord-bots/rp-helper/src/redis"
	"github.com/anhgelus/discord-bots/rp-helper/src/utils"
	"github.com/bwmarrin/discordgo"
)

func Remove(client *discordgo.Session, i *discordgo.InteractionCreate) {
	resp := responseBuilder{I: i, C: client}
	resp.IsEphemeral()
	optMap := generateOptionMap(i)
	opt, ok := optMap["user"]
	var err error
	if !ok {
		err = resp.Message("Remove not given").Send()
		if err != nil {
			utils.SendAlert("remove.go - No member reply", err.Error())
		}
		return
	}
	u := opt.UserValue(client)
	if u.Bot {
		err = resp.Message("Impossible to remove a bot!").Send()
		if err != nil {
			utils.SendAlert("remove.go - Bot given reply", err.Error())
		}
		return
	}
	p := redis.Player{
		DiscordID: u.ID,
		GuildID:   i.GuildID,
	}
	err = p.Load()
	if err != nil {
		err = resp.Message("Error while loading the player").Send()
		if err != nil {
			utils.SendAlert("remove.go - Loading player", err.Error())
		}
		return
	}
	err = p.Remove()
	if err != nil {
		utils.SendAlert("remove.go - Saving player", err.Error())
		err = resp.Message("Error while removing player").Send()
		if err != nil {
			utils.SendAlert("remove.go - Player not saved reply", err.Error())
		}
		return
	}
	err = resp.Message("Player removed!").Send()
	if err != nil {
		utils.SendAlert("add.go - Player saved reply", err.Error())
	}
}
