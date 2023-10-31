package commands

import (
	"fmt"
	"github.com/anhgelus/discord-bots/les-copaings/src/db/sql"
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/bwmarrin/discordgo"
)

func RemoveUserId(client *discordgo.Session, i *discordgo.InteractionCreate) {
	optionMap := generateOptionMap(i)

	if opt, ok := optionMap["id"]; ok {
		err := sql.DB.Model(&sql.Copaing{}).Delete(&sql.Copaing{}).Where("user_id = ?", opt.Value).Error
		if err != nil {
			utils.SendAlert("remove_user_id.go", err.Error())
		}
		err = respondInteraction(client, i, fmt.Sprintf("<@%s> a bien été reset", opt.Value))
		if err != nil {
			utils.SendAlert("remove_user_id.go - Respond interaction", err.Error())
		}
	} else {
		err := respondInteraction(client, i, "L'id n'a pas été renseigné !")
		if err != nil {
			utils.SendAlert("remove_user_id.go - Respond interaction no member", err.Error())
		}
	}
}
