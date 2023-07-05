package start

import (
	"fmt"
	cmd "github.com/anhgelus/discord-bots/les-copaings/src/commands"
	event "github.com/anhgelus/discord-bots/les-copaings/src/events"
	"github.com/anhgelus/discord-bots/les-copaings/src/timers"
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

func Bot(token string, resetEach uint) {
	dg, err := discordgo.New("Bot " + token) // Define connection to discord API with bot token
	if err != nil {
		utils.SendAlert("bot.go - Token", err.Error())
	}

	err = dg.Open() // Bot start
	if err != nil {
		utils.SendAlert("bot.go - Start", err.Error())
	}
	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		utils.SendSuccess(fmt.Sprintf("Bot started as %s", s.State.User.Username))
	})

	initCommands()
	utils.SendSuccess("Commands generated")
	go func() {
		RegisterCommands(dg)
		utils.SendSuccess("Commands registered")
	}()
	CommandHandlers(dg)
	dg.AddHandler(event.ReactionAdd)
	dg.AddHandler(event.MessageSent)
	dg.AddHandler(event.ConnectionVocal)
	dg.AddHandler(event.DisconnectionVocal)
	dg.AddHandler(event.GuildCreate)
	dg.AddHandler(event.GuildDelete)

	timers.SetupTimers(resetEach, dg)

	dg.Identify.Intents = discordgo.IntentsAll

	dg.StateEnabled = true

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	err = dg.Close() // Bot Shutown
	if err != nil {
		utils.SendAlert("bot.go - Shutdown", err.Error())
	}
}

func initCommands() {
	var adminPerm int64 = discordgo.PermissionManageServer
	cmds = append(cmds, Cmd{
		ApplicationCommand: discordgo.ApplicationCommand{
			Name:        "ping",
			Description: "Juste un ping",
		},
		Handler: cmd.Ping,
	}, Cmd{
		ApplicationCommand: discordgo.ApplicationCommand{
			Name:        "avatar",
			Description: "Obtenez votre avatar",
		},
		Handler: cmd.Avatar,
	}, Cmd{
		ApplicationCommand: discordgo.ApplicationCommand{
			Name:        "top",
			Description: "Obtenez le top du serveur",
		},
		Handler: cmd.Top,
	}, Cmd{
		ApplicationCommand: discordgo.ApplicationCommand{
			Name:        "rank",
			Description: "Obtenez votre rang",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "membre",
					Description: "Le rang du membre que vous souhaitez obtenir",
					Required:    false,
				},
			},
		},
		Handler: cmd.Rank,
	}, Cmd{
		ApplicationCommand: discordgo.ApplicationCommand{
			Name:        "purge",
			Description: "Purgez les membres",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "whitelist",
					Description: "Les rôles que vous souhaitez garder (forme : `ID,ID,ID`)",
					Required:    true,
				},
			},
		},
		Handler: cmd.Purge,
	}, Cmd{
		ApplicationCommand: discordgo.ApplicationCommand{
			Name:        "config",
			Description: "Mise à jour de la configuration",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "id",
					Description: "ID à mettre à jour",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Rôles liés à l'XP",
							Value: "xp-roles",
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "valeur",
					Description: "Nouvelle valeur",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "arg1",
					Description: "Argument 1",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "arg2",
					Description: "Argument 2",
					Required:    false,
				},
			},
			DefaultMemberPermissions: &adminPerm,
		},
		Handler: cmd.Config,
	})
}
