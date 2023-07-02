package start

import (
	cmds "github.com/anhgelus/discord-bots/les-copaings/src/commands"
	"github.com/bwmarrin/discordgo"
)

var (
	commandsHandler = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			cmds.Ping(s, i)
		},
		"avatar": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			cmds.Avatar(s, i)
		},
		"rank": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			cmds.Rank(s, i)
		},
		"top": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			cmds.Top(s, i)
		},
	}
)

func CommandHandlers(s *discordgo.Session) {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandsHandler[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}
