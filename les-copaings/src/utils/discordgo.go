package utils

import (
	"github.com/bwmarrin/discordgo"
	"strconv"
	"time"
)

func FetchGuildUser(s *discordgo.Session, guildID string) []*discordgo.Member {
	member, err := s.GuildMembers(guildID, "", 1000)
	if err != nil {
		SendAlert("discordgo.go - Failed to fetch guild users", err.Error())
	}
	return member
}

func GetTimestampFromId(id string) (time.Time, error) {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return time.UnixMilli(0), err
	}

	// https://discord.com/developers/docs/reference#snowflakes-snowflake-id-format-structure-left-to-right
	timestamp := (idInt >> 22) + 1420070400000

	return time.UnixMilli(timestamp), nil
}
