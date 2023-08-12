package commands

import (
	"fmt"
	"github.com/anhgelus/discord-bots/les-copaings/src/db/sql"
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/anhgelus/discord-bots/les-copaings/src/xp"
	"github.com/bwmarrin/discordgo"
)

func Rank(client *discordgo.Session, i *discordgo.InteractionCreate) {
	optionMap := generateOptionMap(i)

	var hasOption bool
	var member *discordgo.Member
	copaing := sql.Copaing{GuildID: i.GuildID}
	if opt, ok := optionMap["membre"]; ok {
		hasOption = ok
		user := opt.UserValue(client)
		var err error
		member, err = client.GuildMember(i.GuildID, user.ID)
		if err != nil {
			utils.SendAlert("rank.go - Fetching the member", err.Error())
			err = respondInteraction(client, i, "Impossible de récupérer le membre!")
			if err != nil {
				utils.SendAlert("rank.go - Respond interaction bot xp", err.Error())
			}
			return
		}
		if user.Bot {
			err = respondInteraction(client, i, "Les bots n'ont pas de niveau !")
			if err != nil {
				utils.SendAlert("rank.go - Respond interaction bot xp", err.Error())
			}
			return
		}
		copaing.UserID = user.ID
	} else {
		hasOption = ok
		member = i.Member
		copaing = sql.Copaing{UserID: i.Member.User.ID, GuildID: i.GuildID}
	}
	result := sql.DB.FirstOrCreate(&copaing, copaing)
	if result.Error != nil {
		utils.SendAlert("rank.go - Querying or creating copaing", result.Error.Error())
		return
	}

	data := xp.NewXp(member, &copaing, 0, hasOption)
	if data.IsNewLevel {
		xp.UpdateRolesNoMessage(&copaing, client)
	}

	var msg string
	level := xp.CalcLevel(copaing.XP)
	nextLvlXp := xp.CalcXpForLevel(level + 1)
	if hasOption {
		msg = fmt.Sprintf("Le niveau de %s", member.User.Username)
	} else {
		msg = "Votre niveau"
	}
	msg = fmt.Sprintf("%s : **%d**\n> XP : %d\n> Prochain niveau dans %d XP",
		msg, level, copaing.XP, nextLvlXp-copaing.XP)
	err := respondInteraction(client, i, msg)
	if err != nil {
		utils.SendAlert("rank.go - Respond interaction", err.Error())
	}
}
