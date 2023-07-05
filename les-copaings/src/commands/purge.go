package commands

import (
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/bwmarrin/discordgo"
	"strings"
)

var removed = 0

func Purge(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := generateOptionMap(i)

	var exclude []string
	if opt, ok := options["whitelist"]; ok {
		exclude = strings.Split(opt.StringValue(), ",")
		if len(exclude) == 0 {
			exclude = []string{strings.Trim(opt.StringValue(), ",")}
		}
		for i, e := range exclude {
			exclude[i] = strings.Trim(e, " ")
		}
	} else {
		err := respondEphemeralInteraction(s, i, "L'argument whitelist n'a pas été renseigné")
		if err != nil {
			utils.SendAlert("purge.go - Respond interaction whitelist", err.Error())
		}
		return
	}

	guild, err := s.State.Guild(i.GuildID)
	if err != nil {
		utils.SendAlert("purge.go - Failed to the guild", err.Error())
		return
	}
	members := utils.FetchGuildUser(s, i.GuildID)
	if len(members) == 0 {
		utils.SendAlert("purge.go - Fetch members", "they are no members")
		return
	}
	err = respondLoadingInteraction(s, i, "Je m'en occupe !")
	if err != nil {
		utils.SendAlert("purge.go - Failed to send loading", err.Error())
		return
	}

	c := make(chan *discordgo.Member)

	go sortMembers(guild, members, exclude, c)

	msg := ""
	for rm := range c {
		if msg == "" {
			msg = "Membres purgés : "
		}
		err = s.GuildMemberDeleteWithReason(i.GuildID, rm.User.ID, "Purge")
		if err != nil {
			msg += rm.User.Username + "(not purged see the console), "
			utils.SendAlert("purge.go - Purging users", err.Error())
			continue
		}
		msg += rm.User.Username + ", "
	}
	if msg == "" {
		msg = "Aucun membre purgé  "
	}
	content := msg[:len(msg)-2]
	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
	if err != nil {
		utils.SendAlert("purge.go - Respond interaction", err.Error())
	}
}

func sortMembers(guild *discordgo.Guild, members []*discordgo.Member, whitelist []string, c chan *discordgo.Member) {
	ownID := guild.OwnerID
	for _, member := range members {
		if member.User.Bot || member.User.ID == ownID {
			continue
		}
		did := false
		for _, r := range member.Roles {
			if did {
				continue
			}
			if utils.AStringContains(whitelist, r) {
				did = true
			}
		}
		if !did {
			c <- member
		}
	}
	close(c)
}
