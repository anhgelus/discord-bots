package commands

import (
	"fmt"
	"github.com/anhgelus/discord-bots/les-copaings/src/db/sql"
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/anhgelus/discord-bots/les-copaings/src/xp"
	"github.com/bwmarrin/discordgo"
)

func Rank(client *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var hasOption bool
	var user *discordgo.User
	copaing := sql.Copaing{GuildID: i.GuildID}
	if opt, ok := optionMap["membre"]; ok {
		hasOption = ok
		user = opt.UserValue(client)
		if user.Bot {
			err := respondInteraction(client, i, "Les bots n'ont pas de niveau !")
			if err != nil {
				utils.SendAlert(err.Error())
			}
			return
		}
		copaing.UserID = user.ID
	} else {
		hasOption = ok
		copaing = sql.Copaing{UserID: i.Member.User.ID, GuildID: i.GuildID}
	}
	result := sql.DB.FirstOrCreate(&copaing, copaing)
	if result.Error != nil {
		utils.SendAlert(result.Error.Error())
		return
	}

	var msg string
	level := xp.CalcLevel(copaing.XP)
	nextLvlXp := xp.CalcXpForLevel(level + 1)
	if hasOption {
		msg = fmt.Sprintf("Le niveau de %s", user.Username)
	} else {
		msg = "Votre niveau"
	}
	msg = fmt.Sprintf("%s : **%d**\n> XP : %d\n> Prochain niveau dans %d XP",
		msg, level, copaing.XP, nextLvlXp-copaing.XP)
	err := respondInteraction(client, i, msg)
	if err != nil {
		utils.SendAlert(err.Error())
	}
}
