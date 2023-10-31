package start

import (
	"fmt"
	cmd "github.com/anhgelus/discord-bots/les-copaings/src/commands"
	event "github.com/anhgelus/discord-bots/les-copaings/src/events"
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Bot(token string) {
	dg, err := discordgo.New("Bot " + token) // Define connection to discord API with bot token
	if err != nil {
		utils.SendAlert("bot.go - Token", err.Error())
	}

	err = dg.Open() // Bot start
	if err != nil {
		utils.SendAlert("bot.go - Start", err.Error())
	}
	go func() {
		time.Sleep(30 * time.Second)
		utils.SendSuccess(fmt.Sprintf("Bot started as %s", dg.State.User.Username))
		utils.NewTimers(30*time.Second, func(_ chan struct{}) {
			rand.NewSource(time.Now().Unix())
			r := rand.Intn(3)
			switch r {
			case 0:
			case 1:
				err = dg.UpdateWatchStatus(0, "Les Copaings")
				if err != nil {
					utils.SendAlert("bot.go - Update status", err.Error())
				}
			case 2:
				err = dg.UpdateGameStatus(0, "Dev by @anhgelus")
				if err != nil {
					utils.SendAlert("bot.go - Update status", err.Error())
				}
			}
		})
	}()

	initCommands()
	utils.SendSuccess("Commands generated")
	go func() {
		RegisterCommands(dg)
		utils.SendSuccess("Commands registered")
	}()
	CommandHandlers(dg)
	//dg.AddHandler(event.ReactionAdd)
	dg.AddHandler(event.MessageSent)
	dg.AddHandler(event.ConnectionVocal)
	dg.AddHandler(event.DisconnectionVocal)
	dg.AddHandler(event.GuildCreate)
	dg.AddHandler(event.GuildDelete)
	dg.AddHandler(event.GuildMemberJoin)
	dg.AddHandler(event.GuildMemberLeave) // event for leave, ban and kick

	dg.Identify.Intents = discordgo.IntentsAll

	dg.StateEnabled = true

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	err = dg.Close() // Bot Shutdown
	if err != nil {
		utils.SendAlert("bot.go - Shutdown", err.Error())
	}
}

func initCommands() {
	var adminPerm int64 = discordgo.PermissionManageServer
	cmds = append(cmds, Cmd{
		ApplicationCommand: discordgo.ApplicationCommand{
			Name:        "ping",
			Description: "Obtenez le ping du bot",
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
			Name:                     "reset_xp",
			Description:              "Reset l'XP",
			DefaultMemberPermissions: &adminPerm,
		},
		Handler: cmd.ResetXP,
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
			Name:        "reset_user_xp",
			Description: "Reset l'XP d'un membre",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "membre",
					Description: "Le membre a reset",
					Required:    true,
				},
			},
			DefaultMemberPermissions: &adminPerm,
		},
		Handler: cmd.ResetUserXP,
	}, Cmd{
		ApplicationCommand: discordgo.ApplicationCommand{
			Name:        "set_level",
			Description: "Modifie le niveau d'un membre",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "membre",
					Description: "Le membre a modifié",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "level",
					Description: "Le niveau",
					Required:    true,
				},
			},
			DefaultMemberPermissions: &adminPerm,
		},
		Handler: cmd.SetLevel,
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
			DefaultMemberPermissions: &adminPerm,
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
						{
							Name:  "Voir les paramètres",
							Value: "show",
						},
						{
							Name:  "Change le salon des news",
							Value: "set-broadcast",
						},
						{
							Name:  "Salons où l'XP est désactivée",
							Value: "disabled-xp-channels",
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "value",
					Description: "Nouvelle valeur",
					Required:    false,
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
