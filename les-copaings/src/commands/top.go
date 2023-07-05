package commands

import (
	"fmt"
	"github.com/anhgelus/discord-bots/les-copaings/src/db/sql"
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/anhgelus/discord-bots/les-copaings/src/xp"
	"github.com/bwmarrin/discordgo"
	"time"
)

func Top(client *discordgo.Session, i *discordgo.InteractionCreate) {
	var tops []sql.Copaing
	sql.DB.Order("xp desc").Limit(10).Where("guild_id = ?", i.GuildID).Find(&tops)
	var msg string
	for i, top := range tops {
		user, err := client.User(top.UserID)
		if err != nil {
			utils.SendAlert("top.go - Failed to get user", err.Error())
			return
		}
		msg += fmt.Sprintf("%d. **%s** - niveau : %d\n", i+1, user.Username, xp.CalcLevel(top.XP))
	}
	err := client.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Top",
					Description: "Les membres les plus actifs du serveur !\n" + msg,
					Author:      &discordgo.MessageEmbedAuthor{Name: i.Member.User.Username},
					Footer: &discordgo.MessageEmbedFooter{
						Text: "© 2023 - Les Copaings",
					},
					Color:     utils.Success,
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		},
	})
	if err != nil {
		utils.SendAlert("top.go - Respond interaction", err.Error())
	}
}
