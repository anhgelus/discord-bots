package commands

import (
	"github.com/anhgelus/discord-bots/les-copaings/src/db/sql"
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/bwmarrin/discordgo"
)

func ResetXP(client *discordgo.Session, i *discordgo.InteractionCreate) {
	sql.DB.Where("guild_id = ?", i.GuildID).Delete(sql.Copaing{})
	err := respondEphemeralInteraction(client, i, "XP Reset")
	if err != nil {
		utils.SendAlert("reset_xp.go - Interaction respond", err.Error())
	}
}
