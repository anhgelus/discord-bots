package sql

import (
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
)

func GetNewRoles(copaing *Copaing, cfg *Config, roles []string, c chan<- string, n chan<- string) {
	for _, xpr := range cfg.XpRoles {
		if copaing.XP < xpr.XP {
			if !utils.AStringContains(roles, xpr.Role) {
				continue
			}
			n <- xpr.Role
			continue
		}
		if utils.AStringContains(roles, xpr.Role) {
			continue
		}
		c <- xpr.Role
	}
	close(c)
	close(n)
}

func GetCopaing(userID string, guildID string) Copaing {
	copaing := Copaing{UserID: userID, GuildID: guildID}
	result := DB.Where("user_id = ? AND guild_id = ?", copaing.UserID, copaing.GuildID).FirstOrCreate(&copaing, copaing)
	if result.Error != nil {
		utils.SendAlert("message.go - Querying/Creating copaing", result.Error.Error())
		return Copaing{}
	}
	return copaing
}
