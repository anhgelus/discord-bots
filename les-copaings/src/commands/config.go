package commands

import (
	"github.com/anhgelus/discord-bots/les-copaings/src/db/sql"
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/anhgelus/discord-bots/les-copaings/src/xp"
	"github.com/bwmarrin/discordgo"
	"strconv"
)

var (
	configs    = []string{"xp-roles"}
	subXpRoles = []string{"add", "edit", "remove"}
)

type configData struct {
	id      string
	value   string
	arg1    string
	arg2    string
	options map[string]*discordgo.ApplicationCommandInteractionDataOption
}

func Config(client *discordgo.Session, i *discordgo.InteractionCreate) {
	options := generateOptionMap(i)
	data := configData{}
	data.options = options
	if opt, ok := options["id"]; ok {
		data.id = opt.StringValue()
	} else {
		err := respondEphemeralInteraction(client, i, "L'argument id n'a pas été renseigné")
		if err != nil {
			utils.SendAlert("config.go - Respond interaction id", err.Error())
		}
		return
	}

	valid := false
	for _, o := range configs {
		if o == data.id {
			valid = true
		}
	}
	if !valid {
		err := respondEphemeralInteraction(client, i, "L'argument id est invalide")
		if err != nil {
			utils.SendAlert("config.go - Respond interaction invalid id", err.Error())
		}
		return
	}

	if opt, ok := options["valeur"]; ok {
		data.value = opt.StringValue()
	} else {
		err := respondEphemeralInteraction(client, i, "L'argument valeur n'a pas été renseigné")
		if err != nil {
			utils.SendAlert("config.go - Respond interaction valeur", err.Error())
		}
		return
	}

	switch data.id {
	case "xp-roles":
		data.xpRoles(client, i)
	}
}

func (data *configData) xpRoles(client *discordgo.Session, i *discordgo.InteractionCreate) {
	valid := false
	msg := ""

	for _, o := range subXpRoles {
		msg += o + ", "
		if o == data.value {
			valid = true
		}
	}
	if !valid {
		err := respondEphemeralInteraction(client, i, "L'argument value est invalide.\nValeurs possibles : "+msg[:len(msg)-2])
		if err != nil {
			utils.SendAlert("config.go - Respond interaction invalid value", err.Error())
		}
		return
	}

	if opt, ok := data.options["arg1"]; ok {
		data.arg1 = opt.StringValue()
	} else {
		err := respondEphemeralInteraction(client, i, "L'argument arg1 n'a pas été renseigné")
		if err != nil {
			utils.SendAlert("config.go - Respond interaction arg1", err.Error())
		}
		return
	}

	role, err := client.State.Role(i.GuildID, data.arg1)
	if role == nil || err != nil {
		err = respondEphemeralInteraction(client, i, "Impossible de trouver le rôle "+data.arg1)
		if err != nil {
			utils.SendAlert("config.go - Respond interaction invalid arg1", err.Error())
		}
		return
	}

	cfg := sql.Config{GuildID: i.GuildID}
	sql.DB.FirstOrCreate(&cfg)

	if data.value == "remove" {
		for id, xpr := range cfg.XpRoles {
			if xpr.Role != data.arg1 {
				continue
			}
			cfg.XpRoles = append(cfg.XpRoles[:id], cfg.XpRoles[id+1:]...)
		}
		sql.DB.Save(&cfg)
		return
	}

	if opt, ok := data.options["arg2"]; ok {
		data.arg2 = opt.StringValue()
	} else {
		err = respondEphemeralInteraction(client, i, "L'argument arg2 n'a pas été renseigné")
		if err != nil {
			utils.SendAlert("config.go - Respond interaction arg2", err.Error())
		}
		return
	}
	lvl, err := strconv.Atoi(data.arg2)
	if err != nil {
		err = respondEphemeralInteraction(client, i, "L'argument arg2 est invalide (impossible de le convertir en entier)")
		if err != nil {
			utils.SendAlert("config.go - Respond interaction invald arg2", err.Error())
		}
		return
	}
	xP := xp.CalcXpForLevel(uint(lvl))

	switch data.value {
	case "add":
		cfg.XpRoles = append(cfg.XpRoles, sql.XpRole{Role: data.arg1, XP: xP})
	case "edit":
		for id, xpr := range cfg.XpRoles {
			if xpr.Role != data.arg1 {
				continue
			}
			xpr.XP = xP
			cfg.XpRoles[id] = xpr
		}
	}
	sql.DB.Save(&cfg)
	err = respondInteraction(client, i, "Valeur enregistrée !")
	if err != nil {
		utils.SendAlert("config.go - Respond interaction value saved", err.Error())
	}
}
