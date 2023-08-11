package commands

import (
	"fmt"
	"github.com/anhgelus/discord-bots/les-copaings/src/db/sql"
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/anhgelus/discord-bots/les-copaings/src/xp"
	"github.com/bwmarrin/discordgo"
)

func ResetUserXP(client *discordgo.Session, i *discordgo.InteractionCreate) {
	optionMap := generateOptionMap(i)

	var member *discordgo.Member
	copaing := sql.Copaing{GuildID: i.GuildID}
	if opt, ok := optionMap["membre"]; ok {
		user := opt.UserValue(client)
		var err error
		member, err = client.GuildMember(i.GuildID, user.ID)
		if err != nil {
			utils.SendAlert("reset_user_xp.go - Fetching the member", err.Error())
			err = respondInteraction(client, i, "Impossible de récupérer le membre!")
			if err != nil {
				utils.SendAlert("reset_user_xp.go - Respond interaction bot xp", err.Error())
			}
			return
		}
		if user.Bot {
			err = respondInteraction(client, i, "Les bots n'ont pas de niveau !")
			if err != nil {
				utils.SendAlert("reset_user_xp.go - Respond interaction bot xp", err.Error())
			}
			return
		}
		copaing.UserID = user.ID
	} else {
		err := respondInteraction(client, i, "Le membre n'a pas été renseigné !")
		if err != nil {
			utils.SendAlert("reset_user_xp.go - Respond interaction no member", err.Error())
		}
		return
	}
	result := sql.DB.FirstOrCreate(&copaing, copaing)
	if result.Error != nil {
		utils.SendAlert("reset_user_xp.go - Querying or creating copaing", result.Error.Error())
		return
	}

	copaing.XP = 0
	sql.Save(&copaing)
	xp.UpdateRolesNoMessage(&copaing, client)

	err := respondInteraction(client, i, fmt.Sprintf("%s a bien été reset", member.User.Username))
	if err != nil {
		utils.SendAlert("reset_user_xp.go - Respond interaction", err.Error())
	}
}
