package commands

import (
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func Purge(client *discordgo.Session, i *discordgo.InteractionCreate) {
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
		err := respondEphemeralInteraction(client, i, "L'argument whitelist n'a pas été renseigné")
		if err != nil {
			utils.SendAlert("purge.go - Respond interaction whitelist", err.Error())
		}
		return
	}

	err := client.RequestGuildMembers(i.GuildID, "", 0, "", false)
	if err != nil {
		utils.SendAlert("purge.go - Failed to request guild members", err.Error())
		return
	}
	guild, err := client.State.Guild(i.GuildID)
	if err != nil {
		utils.SendAlert("purge.go - Failed to the guild", err.Error())
		return
	}
	members, err := client.GuildMembers(i.GuildID, "", 0)
	if err != nil {
		utils.SendAlert("purge.go - Failed to get guild members", err.Error())
		return
	}
	toRemove := members
	ownID := guild.OwnerID
	for id, member := range members {
		utils.SendDebug("for each", member.User.Username, id)
		if member.User.Bot || member.User.ID == ownID {
			toRemove = removeMember(toRemove, member, id)
			continue
		}
		did := false
		for _, r := range member.Roles {
			if did {
				continue
			}
			if utils.AStringContains(exclude, r) {
				utils.SendDebug("remove", member.User.Username)
				toRemove = removeMember(toRemove, member, id)
				did = true
			}
		}
	}
	msg := ""
	for _, rm := range toRemove {
		//err = client.GuildMemberDeleteWithReason(i.GuildID, rm.User.ID, "Purge")
		//if err != nil {
		//	msg += rm.User.Username + "(not purged see the console), "
		//	utils.SendAlert(err.Error())
		//	continue
		//}
		msg += rm.User.Username + ", "
	}
	if msg == "" {
		msg = "Aucun membre purgé  "
	}
	err = respondEphemeralInteraction(client, i, msg[:len(msg)-2])
	if err != nil {
		utils.SendAlert("purge.go - Respond interaction", err.Error())
	}
}

func removeMember(arr []*discordgo.Member, member *discordgo.Member, id int) []*discordgo.Member {
	if id == len(arr)-1 {
		arr = arr[:id]
		return arr
	}
	arr = append(arr[:id], arr[id+1:]...)
	return arr
}
