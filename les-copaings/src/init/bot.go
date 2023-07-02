package start

import (
	event "github.com/anhgelus/discord-bots/les-copaings/src/events"
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

func Bot(token string) {
	dg, err := discordgo.New("Bot " + token) // Define connection to discord API with bot token
	if err != nil {
		utils.SendAlert(err.Error())
	}

	err = dg.Open() // Bot start
	if err != nil {
		utils.SendAlert(err.Error())
	}

	utils.SendSucces("Bot started")
	Command(dg)
	CommandHandlers(dg)
	dg.AddHandler(event.ReactionAdd)

	dg.Identify.Intents = discordgo.IntentMessageContent | discordgo.IntentsMessageContent | discordgo.IntentGuildMembers | discordgo.IntentGuildMessages

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	err = dg.Close() // Bot Shutown
	if err != nil {
		utils.SendAlert(err.Error())
	}
}
