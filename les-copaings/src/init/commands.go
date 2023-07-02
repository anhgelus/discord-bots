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
		Description: "Basic command to check if the bot respond",
	})
	commands = append(commands, &discordgo.ApplicationCommand{
		Name:        "avatar",
		Description: "Display self avatar",
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
