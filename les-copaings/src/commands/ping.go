package commands

import (
	"fmt"
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/bwmarrin/discordgo"
	"time"
)

func Ping(client *discordgo.Session, i *discordgo.InteractionCreate) {
	err := client.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		utils.SendAlert("ping.go - Respond interaction", err.Error())
	}

	response, err := client.InteractionResponse(i.Interaction)
	if err != nil {
		utils.SendAlert("ping.go - Interaction response", err.Error())
	}

	var msg string

	interactionTimestamp, err := utils.GetTimestampFromId(i.ID)
	if err != nil {
		utils.SendAlert("ping.go - Get timestamp from ID", err.Error())
		msg = "� Pong !"
	} else {
		utils.SendDebug(interactionTimestamp.Format(time.UnixDate))
		msg = fmt.Sprintf(
			":ping_pong: Pong !\nLatence du bot : `%d ms`\nLatence de l'API discord : `%d ms`",
			response.Timestamp.Sub(interactionTimestamp).Milliseconds(),
			client.HeartbeatLatency().Milliseconds(),
		)
	}

	_, err = client.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})

	if err != nil {
		utils.SendAlert("ping.go - Interaction response edit", err.Error())
	}
}
