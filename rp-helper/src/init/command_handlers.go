package start

import (
	"github.com/anhgelus/discord-bots/rp-helper/src/utils"
	"github.com/bwmarrin/discordgo"
)

var (
	commandsHandler = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
)

func CommandHandlers(s *discordgo.Session) {
	if len(commandsHandler) == 0 {
		genCmds()
	}
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandsHandler[i.ApplicationCommandData().Name]; ok {
			v, err := utils.ComesFromDM(s, i.ChannelID)
			if err != nil {
				utils.SendAlert("command_handlers.go - Checking if is coming from dm", err.Error())
			}
			if v {
				return
			}
			h(s, i)
		}
	})
}

func genCmds() {
	for _, c := range cmds {
		commandsHandler[c.Name] = c.Handler
	}
}
