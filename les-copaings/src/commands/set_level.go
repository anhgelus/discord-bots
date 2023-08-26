package commands

import (
	"fmt"
	"github.com/anhgelus/discord-bots/les-copaings/src/db/sql"
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/anhgelus/discord-bots/les-copaings/src/xp"
	"github.com/bwmarrin/discordgo"
)

func SetLevel(client *discordgo.Session, i *discordgo.InteractionCreate) {
	optionMap := generateOptionMap(i)

	var member *discordgo.Member
	copaing := sql.Copaing{GuildID: i.GuildID}
	if opt, ok := optionMap["membre"]; ok {
		user := opt.UserValue(client)
		var err error
		member, err = client.GuildMember(i.GuildID, user.ID)
		if err != nil {
			utils.SendAlert("set_level.go - Fetching the member", err.Error())
			err = respondInteraction(client, i, "Impossible de récupérer le membre!")
			if err != nil {
				utils.SendAlert("set_level.go - Respond interaction bot xp", err.Error())
			}
			return
		}
		if user.Bot {
			err = respondInteraction(client, i, "Les bots n'ont pas de niveau !")
			if err != nil {
				utils.SendAlert("set_level.go - Respond interaction bot xp", err.Error())
			}
			return
		}
		copaing.UserID = user.ID
	} else {
		err := respondInteraction(client, i, "Le membre n'a pas été renseigné !")
		if err != nil {
			utils.SendAlert("set_level.go - Respond interaction no member", err.Error())
		}
		return
	}

	var lvl uint
	if opt, ok := optionMap["level"]; ok {
		lvl = uint(opt.UintValue())
	} else {
		err := respondInteraction(client, i, "Le niveau n'a pas été renseigné !")
		if err != nil {
			utils.SendAlert("set_level.go - Respond interaction no level", err.Error())
		}
		return
	}

	result := sql.DB.FirstOrCreate(&copaing, copaing)
	if result.Error != nil {
		utils.SendAlert("set_level.go - Querying or creating copaing", result.Error.Error())
		return
	}

	copaing.XP = xp.CalcXpForLevel(lvl)
	sql.Save(&copaing)
	xp.UpdateRolesNoMessage(&copaing, client)

	err := respondInteraction(client, i, fmt.Sprintf("Le niveau de %s a bien été mis-à-jour", member.User.Username))
	if err != nil {
		utils.SendAlert("set_level.go - Respond interaction", err.Error())
	}
}
