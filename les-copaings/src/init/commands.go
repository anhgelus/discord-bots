package start

import (
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/bwmarrin/discordgo"
	"strconv"
)

var commands []*discordgo.ApplicationCommand

func init() {
	commands = append(commands, &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Juste un ping",
	})
	commands = append(commands, &discordgo.ApplicationCommand{
		Name:        "avatar",
		Description: "Obtenez votre avatar",
	})
	commands = append(commands, &discordgo.ApplicationCommand{
		Name:        "top",
		Description: "Obtenez le top du serveur",
	})
	commands = append(commands, &discordgo.ApplicationCommand{
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
	})
}

func Command(client *discordgo.Session) {
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	o := 0
	for i, v := range commands {
		cmd, err := client.ApplicationCommandCreate(client.State.User.ID, "", v)
		if err != nil {
			utils.SendAlert(err.Error())
			return
		}
		registeredCommands[i] = cmd
		utils.SendSuccess("[COMMAND] : " + v.Name + " initialized")
		o += 1
	}
	utils.SendSuccess("[Recaps] " + strconv.Itoa(o) + "/" + strconv.Itoa(len(commands)) + " commands has been loaded")
}
