package utils

import (
	"github.com/bwmarrin/discordgo"
)

func FetchGuildUser(s *discordgo.Session, guildID string) []*discordgo.Member {
	member, err := s.GuildMembers(guildID, "", 1000)
	if err != nil {
		SendAlert("discordgo.go - Failed to fetch guild users", err.Error())
	}
	return member
}
