package xp

import (
	"fmt"
	"github.com/anhgelus/discord-bots/les-copaings/src/db/sql"
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/bwmarrin/discordgo"
	"sync"
)

// UpdateRoles will update the role of the copaing and send a message
func UpdateRoles(copaing *sql.Copaing, client *discordgo.Session, event *discordgo.MessageCreate) {
	added, lost := updateRoles(copaing, client)
	if added > 1 {
		_, err := client.ChannelMessageSend(event.ChannelID,
			fmt.Sprintf("%s a gagné %d rôles !", event.Author.Mention(), added),
		)
		if err != nil {
			utils.SendAlert("roles.go - Cannot send message for roles added", err.Error())
			return
		}
	} else if added == 1 {
		_, err := client.ChannelMessageSend(event.ChannelID,
			fmt.Sprintf("%s a gagné 1 rôle !", event.Author.Mention()),
		)
		if err != nil {
			utils.SendAlert("roles.go - Cannot send message for role added", err.Error())
			return
		}
	}

	if lost > 1 {
		_, err := client.ChannelMessageSend(event.ChannelID,
			fmt.Sprintf("%s a perdu %d rôles !", event.Author.Mention(), lost),
		)
		if err != nil {
			utils.SendAlert("roles.go - Cannot send message for roles lost", err.Error())
			return
		}
	} else if lost == 1 {
		_, err := client.ChannelMessageSend(event.ChannelID,
			fmt.Sprintf("%s a perdu 1 rôle !", event.Author.Mention()),
		)
		if err != nil {
			utils.SendAlert("roles.go - Cannot send message for role lost", err.Error())
			return
		}
	}
}

// UpdateRolesNoMessage will update the role of the copaing
func UpdateRolesNoMessage(copaing *sql.Copaing, client *discordgo.Session) {
	_, _ = updateRoles(copaing, client)
}

// updateRoles will update roles and return the total of modified role
func updateRoles(copaing *sql.Copaing, client *discordgo.Session) (uint, uint) {
	cfg := sql.Config{GuildID: copaing.GuildID}
	sql.DB.Model(&sql.Config{}).Where("guild_id = ?", cfg.GuildID).Preload("XpRoles").FirstOrCreate(&cfg)

	roles := make(chan string)
	notRoles := make(chan string)

	member, err := client.GuildMember(copaing.GuildID, copaing.UserID)
	if err != nil {
		utils.SendAlert("roles.go - Cannot fetch the member", err.Error())
	}

	go sql.GetNewRoles(copaing, &cfg, member.Roles, roles, notRoles)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	added := 0
	lost := 0

	go func() {
		for role := range roles {
			err = client.GuildMemberRoleAdd(copaing.GuildID, copaing.UserID, role)
			if err != nil {
				utils.SendAlert("roles.go - Role add no msg", err.Error())
				continue
			}
			added++
		}
		wg.Done()
	}()

	go func() {
		for role := range notRoles {
			err = client.GuildMemberRoleRemove(copaing.GuildID, copaing.UserID, role)
			if err != nil {
				utils.SendAlert("message.go - Role remove no msg", err.Error())
				continue
			}
			lost++
		}
		wg.Done()
	}()

	wg.Wait()

	return uint(added), uint(lost)
}
