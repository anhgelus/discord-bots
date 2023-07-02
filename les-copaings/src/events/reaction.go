package events

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func ReactionAdd(client *discordgo.Session, event *discordgo.MessageReactionAdd){
	err := client.GuildMemberRoleAdd(event.GuildID, event.UserID, "1085295165054402690")
	if err != nil { fmt.Println(err) }
}