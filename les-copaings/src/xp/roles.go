package xp

import (
	"fmt"
	"github.com/anhgelus/discord-bots/les-copaings/src/db/sql"
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/bwmarrin/discordgo"
	"sync"
)

func UpdateRoles(copaing *sql.Copaing, client *discordgo.Session, event *discordgo.MessageCreate) {
	cfg := sql.Config{GuildID: copaing.GuildID}
	sql.DB.Model(&sql.Config{}).Where("guild_id = ?", cfg.GuildID).Preload("XpRoles").FirstOrCreate(&cfg)

	roles := make(chan string)
	notRoles := make(chan string)

	go sql.GetNewRoles(copaing, &cfg, event.Member.Roles, roles, notRoles)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for role := range roles {
			err := client.GuildMemberRoleAdd(copaing.GuildID, copaing.UserID, role)
			if err != nil {
				utils.SendAlert("roles.go - Role add", err.Error())
				_, err = client.ChannelMessageSend(event.ChannelID, "Impossible de vous ajouter le rôle "+role)
				if err != nil {
					utils.SendAlert("roles.go - Message send role failed", err.Error())
				}
				continue
			}
			_, err = client.ChannelMessageSend(event.ChannelID,
				fmt.Sprintf("<@%s>, vous venez d'obtenir un nouveau rôle !", copaing.UserID),
			)
			if err != nil {
				utils.SendAlert("roles.go - New role message", err.Error())
			}
		}
		wg.Done()
	}()

	go func() {
		for role := range notRoles {
			err := client.GuildMemberRoleRemove(copaing.GuildID, copaing.UserID, role)
			if err != nil {
				utils.SendAlert("roles.go - Role remove", err.Error())
				_, err = client.ChannelMessageSend(event.ChannelID, "Impossible de vous supprimer le rôle "+role)
				if err != nil {
					utils.SendAlert("roles.go - Message send role failed", err.Error())
				}
				continue
			}
			_, err = client.ChannelMessageSend(event.ChannelID,
				fmt.Sprintf("<@%s>, vous avez perdu un rôle !", copaing.UserID),
			)
			if err != nil {
				utils.SendAlert("roles.go - Role lost message", err.Error())
			}
		}
		wg.Done()
	}()

	wg.Wait()
}

func UpdateRolesNoMessage(copaing *sql.Copaing, client *discordgo.Session) {
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

	go func() {
		for role := range roles {
			err = client.GuildMemberRoleAdd(copaing.GuildID, copaing.UserID, role)
			if err != nil {
				utils.SendAlert("roles.go - Role add no msg", err.Error())
				continue
			}
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
		}
		wg.Done()
	}()

	wg.Wait()
}
