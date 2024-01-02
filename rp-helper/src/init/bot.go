package start

import (
	"fmt"
	cmd "github.com/anhgelus/discord-bots/rp-helper/src/commands"
	"github.com/anhgelus/discord-bots/rp-helper/src/utils"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var Debug bool

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
				err = dg.UpdateWatchStatus(0, "RP Helper")
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
			Description: "Get the ping of the bot",
		},
		Handler: cmd.Ping,
	}, Cmd{
		ApplicationCommand: discordgo.ApplicationCommand{
			Name:        "profile",
			Description: "Get information about yourself",
		},
		Handler: cmd.Profile,
	}, Cmd{
		ApplicationCommand: discordgo.ApplicationCommand{
			Name:        "list",
			Description: "List players",
		},
		Handler: cmd.List,
	}, Cmd{
		ApplicationCommand: discordgo.ApplicationCommand{
			Name:                     "generate_goals",
			Description:              "Generate goals for every players",
			DefaultMemberPermissions: &adminPerm,
		},
		Handler: cmd.GenerateGoals,
	}, Cmd{
		ApplicationCommand: discordgo.ApplicationCommand{
			Name:        "add",
			Description: "Add a player",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "User to add",
					Required:    true,
				},
			},
			DefaultMemberPermissions: &adminPerm,
		},
		Handler: cmd.Add,
	}, Cmd{
		ApplicationCommand: discordgo.ApplicationCommand{
			Name:        "remove",
			Description: "Remove a player",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "User to remove",
					Required:    true,
				},
			},
			DefaultMemberPermissions: &adminPerm,
		},
		Handler: cmd.Remove,
	})
}
