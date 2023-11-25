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
	err := respondLoadingInteraction(client, i)
	if err != nil {
		utils.SendAlert("top.go - Failed to make response defer", err.Error())
		return
	}
	go func() {
		m := getTops(client, i)
		_, err = client.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Top",
					Description: "Les membres les plus actifs du serveur !\n" + m,
					Author:      &discordgo.MessageEmbedAuthor{Name: i.Member.User.Username},
					Footer: &discordgo.MessageEmbedFooter{
						Text: "© 2023 - Les Copaings",
					},
					Color:     utils.Success,
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		})
		if err != nil {
			utils.SendAlert("top.go - Respond interaction", err.Error())
		}
	}()
}

func getTops(client *discordgo.Session, i *discordgo.InteractionCreate) string {
	var tops []sql.Copaing
	sql.DB.Order("xp desc").Limit(10).Where("guild_id = ?", i.GuildID).Find(&tops)
	var msg string
	reload := false
	for i, top := range tops {
		member, err := client.GuildMember(top.GuildID, top.UserID)
		if err != nil {
			utils.SendAlert("top.go - Failed to get member", err.Error())
			msg += fmt.Sprintf("%d. **<@%s>** - niveau : %d\n", i+1, top.UserID, xp.CalcLevel(top.XP))
			continue
		}
		//data := xp.NewXpNoUpdate(member, &top, 0)
		//if data.IsNewLevel {
		//	reload = true
		//}
		msg += fmt.Sprintf("%d. **%s** - niveau : %d\n", i+1, member.User.Username, xp.CalcLevel(top.XP))
	}
	if reload {
		return getTops(client, i)
	}
	return msg
}
