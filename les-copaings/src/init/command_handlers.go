package start

import (
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
			//user := redis.GenerateConnectedUser(i.Member)
			//time := user.TimeSinceLastEvent()
			//reduce := xp.CalcXpLose(utils.HoursOfUnix(time))
			//user.UpdateLastEvent()

			h(s, i)
		}
	})
}

func genCmds() {
	for _, c := range cmds {
		commandsHandler[c.Name] = c.Handler
	}
}
