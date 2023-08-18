package commands

import (
	"fmt"
	"github.com/anhgelus/discord-bots/les-copaings/src/db/sql"
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/anhgelus/discord-bots/les-copaings/src/xp"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"time"
)

var (
	configs   = []string{"xp-roles", "show", "set-broadcast", "disabled-xp-channels"}
	subManage = []string{"add", "edit", "remove"}
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

	if data.id == "show" {
		data.showConfig(client, i)
		return
	}

	if opt, ok := options["value"]; ok {
		data.value = opt.StringValue()
	} else {
		err := respondEphemeralInteraction(client, i, "L'argument valeur n'a pas été renseigné")
		if err != nil {
			utils.SendAlert("config.go - Respond interaction value", err.Error())
		}
		return
	}

	switch data.id {
	case "xp-roles":
		data.xpRoles(client, i)
	case "set-broadcast":
		data.setBroadcast(client, i)
	case "disabled-xp-channels":
		data.disabledXpChannels(client, i)
	default:
		utils.SendAlert("config.go - Switch id", "not handled "+data.id)
	}
}

func (data *configData) xpRoles(client *discordgo.Session, i *discordgo.InteractionCreate) {
	valid := false
	msg := ""

	for _, o := range subManage {
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
	loadConfig(&cfg)

	if data.value == "remove" {
		for id, xpr := range cfg.XpRoles {
			if xpr.Role != data.arg1 {
				continue
			}
			cfg.XpRoles = append(cfg.XpRoles[:id], cfg.XpRoles[id+1:]...)
		}
		sql.Save(&cfg)
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
			utils.SendAlert("config.go - Respond interaction invalid arg2", err.Error())
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
	sql.Save(&cfg)
	err = respondInteraction(client, i, "Valeur enregistrée !")
	if err != nil {
		utils.SendAlert("config.go - Respond interaction value saved", err.Error())
	}
}

func (data *configData) disabledXpChannels(client *discordgo.Session, i *discordgo.InteractionCreate) {
	valid := false
	msg := ""

	for _, o := range subManage {
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

	role, err := client.Channel(data.arg1)
	if role == nil || err != nil {
		err = respondEphemeralInteraction(client, i, "Impossible de trouver le salon "+data.arg1)
		if err != nil {
			utils.SendAlert("config.go - Respond interaction invalid arg1", err.Error())
		}
		return
	}

	cfg := sql.Config{GuildID: i.GuildID}
	loadConfig(&cfg)

	if data.value == "remove" {
		for id, dxp := range cfg.DisabledXpChannel {
			if dxp != data.arg1 {
				continue
			}
			cfg.XpRoles = append(cfg.XpRoles[:id], cfg.XpRoles[id+1:]...)
		}
		sql.Save(&cfg)
		return
	}

	switch data.value {
	case "add":
		cfg.DisabledXpChannel = append(cfg.DisabledXpChannel, data.arg1)
	case "edit":
		err = respondInteraction(client, i, "Edit n'est pas supporté !")
		if err != nil {
			utils.SendAlert("config.go - Respond interaction value saved", err.Error())
		}
		return
	}
	sql.Save(&cfg)
	err = respondInteraction(client, i, "Valeur enregistrée !")
	if err != nil {
		utils.SendAlert("config.go - Respond interaction value saved", err.Error())
	}
}

func (data *configData) setBroadcast(client *discordgo.Session, i *discordgo.InteractionCreate) {
	cfg := sql.Config{GuildID: i.GuildID}
	loadConfig(&cfg)

	_, err := client.Channel(data.value)
	if err != nil {
		err = respondEphemeralInteraction(client, i, "Impossible de récupérer le salon avec l'id "+data.value)
		if err != nil {
			utils.SendAlert("config.go - Respond interaction invalid value", err.Error())
		}
		return
	}
	cfg.BroadcastChannel = data.value
	sql.Save(&cfg)

	err = respondInteraction(client, i, "Changement effectué !")
	if err != nil {
		utils.SendAlert("config.go - Respond interaction invalid value", err.Error())
	}
}

func (data *configData) showConfig(client *discordgo.Session, i *discordgo.InteractionCreate) {
	cfg := sql.Config{GuildID: i.GuildID}
	loadConfig(&cfg)

	var embeds []*discordgo.MessageEmbed

	main := discordgo.MessageEmbed{
		Title:       "Config",
		Description: "Configuration du serveur\n",
		Author:      &discordgo.MessageEmbedAuthor{Name: i.Member.User.Username},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "© 2023 - Les Copaings",
		},
		Color:     utils.Success,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	var brcChan string
	if cfg.BroadcastChannel == "" {
		brcChan = "Pas de salon :(\nUtilisez la commande `/config` pour le setup !"
	} else {
		brcChan = fmt.Sprintf("<#%s>", cfg.BroadcastChannel)
	}
	fields := []*discordgo.MessageEmbedField{
		{
			Name:   "Salon annonces",
			Value:  brcChan,
			Inline: false,
		},
	}
	main.Fields = fields

	xpRoles := discordgo.MessageEmbed{
		Title:       "Rôles liés aux niveaux",
		Description: "Liste des rôles:\n",
		Author:      &discordgo.MessageEmbedAuthor{Name: i.Member.User.Username},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "© 2023 - Les Copaings",
		},
		Color:     utils.Success,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	for _, xpr := range cfg.XpRoles {
		field := discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("Niveau %d", xp.CalcLevel(xpr.XP)),
			Value:  fmt.Sprintf("<@&%s>", xpr.Role),
			Inline: false,
		}
		xpRoles.Fields = append(xpRoles.Fields, &field)
	}

	embeds = append(embeds, &main, &xpRoles)

	err := client.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embeds,
		},
	})

	if err != nil {
		utils.SendAlert("config.go - Respond interaction show", err.Error())
	}
}

func loadConfig(cfg *sql.Config) {
	sql.DB.Where("guild_id = ?", cfg.GuildID).Preload("XpRoles").FirstOrCreate(cfg)
}
