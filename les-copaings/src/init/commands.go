package start

import (
	"fmt"
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/bwmarrin/discordgo"
)

type Cmd struct {
	discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var cmds []Cmd

func Command(client *discordgo.Session) {
	registeredCommands := make([]*discordgo.ApplicationCommand, len(cmds))
	o := 0
	for i, v := range cmds {
		cmd, err := client.ApplicationCommandCreate(client.State.User.ID, "", &v.ApplicationCommand)
		if err != nil {
			utils.SendAlert(err.Error())
			return
		}
		registeredCommands[i] = cmd
		utils.SendSuccess(fmt.Sprintf("[COMMAND] : %s initialized", v.Name))
		o += 1
	}
	utils.SendSuccess(fmt.Sprintf("[Recaps] %d/%d commands has been loaded", o, len(cmds)))
}
