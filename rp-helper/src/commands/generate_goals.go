package commands

import (
	"fmt"
	"github.com/anhgelus/discord-bots/rp-helper/src/config"
	"github.com/anhgelus/discord-bots/rp-helper/src/redis"
	"github.com/anhgelus/discord-bots/rp-helper/src/utils"
	"github.com/bwmarrin/discordgo"
	"math"
	"math/rand"
	"slices"
)

func GenerateGoals(client *discordgo.Session, i *discordgo.InteractionCreate) {
	resp := responseBuilder{C: client, I: i}
	ps := redis.GetPlayers(i.GuildID)
	var pss []string
	for _, p := range ps {
		pss = append(pss, fmt.Sprintf("<@%s>", p.DiscordID))
	}

	n := len(ps)
	if n/2+1 > len(config.Objs.Mains) {
		err := resp.Message("Too much players for main objectives").IsEphemeral().Send()
		if err != nil {
			utils.SendAlert("generate_goals.go - Too much players 1", err.Error())
		}
		return
	}
	second := int(math.Ceil(float64(len(config.Objs.Secondaries)) * 1.5))
	if n > second {
		err := resp.Message("Too much players for secondary objectives").IsEphemeral().Send()
		if err != nil {
			utils.SendAlert("generate_goals.go - Too much players 2", err.Error())
		}
		return
	}
	err := resp.IsDeferred().Send()
	if err != nil {
		utils.SendAlert("generate_goals.go - Sending defer", err.Error())
	}
	objMainsBrut := config.Objs.Mains
	l := len(config.Objs.Mains)
	var objMains []string
	if n == len(config.Objs.Mains)/2 || n == len(config.Objs.Mains)/2+1 {
		objMains = config.GenerateMainGoals(config.Objs.Mains, pss)
	} else if n%2 == 0 {
		for i := 0; i < l-n/2; i++ {
			r := rand.Intn(len(objMainsBrut))
			slices.Delete(objMains, r, r+1)
		}
	} else {
		for i := 0; i < l-n/2+1; i++ {
			r := rand.Intn(len(objMainsBrut))
			slices.Delete(objMains, r, r+1)
		}
	}
	objMains = config.GenerateMainGoals(objMainsBrut, pss)
	objSecs := config.Objs.Secondaries
	lS := len(config.Objs.Secondaries)
	for i := 0; i < lS-second; i++ {
		r := rand.Intn(len(objSecs))
		slices.Delete(objSecs, r, r+1)
	}

	hasError := false
	for _, p := range ps {
		r := rand.Intn(len(objMains))
		p.Goals.Main = objMains[r]
		slices.Delete(objMains, r, r+1)
		for i := 0; i < config.Objs.Settings.NumberOfSecondaries; i++ {
			r = rand.Intn(len(objSecs))
			p.Goals.Secondaries[i] = objSecs[r].GenerateGoal(pss)
		}
		err := p.Save()
		if err != nil {
			hasError = true
			utils.SendAlert("generate_goals.go - Saving player", err.Error())
		}
	}
	if hasError {
		resp.Message("Goals generated! (But there is an error)")
	} else {
		resp.Message("Goals generated!")
	}
	err = resp.IsEdit().Send()
	if err != nil {
		utils.SendAlert("generate_goals.go - Sending reply", err.Error())
	}
}
