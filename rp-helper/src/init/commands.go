package start

import (
	"fmt"
	"github.com/anhgelus/discord-bots/rp-helper/src/utils"
	"github.com/bwmarrin/discordgo"
)

type Cmd struct {
	discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var cmds []Cmd

func RegisterCommands(client *discordgo.Session) {
	registeredCommands := make([]*discordgo.ApplicationCommand, len(cmds))
	o := 0
	guildID := ""
	if Debug {
		gs, err := client.UserGuilds(1, "", "")
		if err != nil {
			utils.SendAlert("commands.go - Fecthing guilds for debug", err.Error())
		} else {
			guildID = gs[0].ID
		}
	}
	for i, v := range cmds {
		cmd, err := client.ApplicationCommandCreate(client.State.User.ID, guildID, &v.ApplicationCommand)
		if err != nil {
			utils.SendAlert("commands.go - Create application command", err.Error())
			continue
		}
		registeredCommands[i] = cmd
		utils.SendSuccess(fmt.Sprintf("[COMMAND]: %s initialized", v.Name))
		o += 1
	}
	utils.SendSuccess(fmt.Sprintf("[Recaps] %d/%d commands has been loaded", o, len(cmds)))
}
