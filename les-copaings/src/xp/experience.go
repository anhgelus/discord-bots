package xp

import (
	"github.com/anhgelus/discord-bots/les-copaings/src/db/redis"
	"github.com/anhgelus/discord-bots/les-copaings/src/db/sql"
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/bwmarrin/discordgo"
	"math"
)

// CalcExperience returns the experience gained from a message's length and diversity
func CalcExperience(length uint, diversity uint) uint {
	// f(x;y) = ((0.025 x^{1.25})/(y^{-0.5}))+1
	return uint(math.Floor(((0.025 * math.Pow(float64(length), 1.25)) / math.Pow(float64(diversity), -0.5)) + 1))
}

// CalcExperienceFromVocal returns the experience gained from a vocal session's length
func CalcExperienceFromVocal(length uint) uint {
	// f(x)=(0.01 x^{1.3})+1
	return uint(math.Floor(0.01*math.Pow(float64(length), 1.3) + 1))
}

// CalcLevel with given xp
func CalcLevel(xp uint) uint {
	// f(x)=0.2 x^0.5
	return uint(math.Floor(0.2 * math.Pow(float64(xp), 0.5)))
}

// CalcXpForLevel returns the total xp required to get the level given
func CalcXpForLevel(level uint) uint {
	// f(x)=0.2 x^0.5
	// f(x)/0.2 = x^0.5
	// (f(x)/0.2)^2 = x
	return uint(math.Floor(math.Pow(5*float64(level), 2)))
}

// CalcXpLose returns the xp to remove with the inactivity length (hour)
func CalcXpLose(inactivity uint) uint {
	// f(x)= x^2*10^{-2+ln(x/85)}
	return uint(math.Floor(
		math.Pow(float64(inactivity), 2) * math.Pow(10, -2+math.Log(float64(inactivity/85))),
	))
}

// NewXpData stores data returned by NewXp
type NewXpData struct {
	// IsNewLevel is true when the copaing get a new level
	IsNewLevel bool
	// OldLevel is the level before the NewXp call
	OldLevel uint
	// NewLevel is the level after the NewXp call
	NewLevel uint
	// LevelUp is true when the copaing get a new level
	LevelUp bool
	// LevelDown is true when the copaing lose a new level
	LevelDown bool
}

// NewXp calculate the new Xp and returns the data stored in NewXpData
func NewXp(member *discordgo.Member, copaing *sql.Copaing, exp uint, update bool) NewXpData {
	user := redis.GenerateConnectedUser(member)
	time := user.TimeSinceLastEvent()
	r := int(CalcXpLose(utils.HoursOfUnix(time))) - int(user.XpLostSaved)
	if r < 0 {
		utils.SendAlert("experience.go - NewXP calc reduce", "Reduce is bellow 0")
		return NewXpData{}
	}
	reduce := uint(r)
	if update {
		user.UpdateLastEvent()
	} else {
		user.UpdateLostXp(reduce)
	}

	oldLvl := CalcLevel(copaing.XP)
	copaing.XP += exp
	if int(copaing.XP)-r < 0 {
		copaing.XP = 0
	} else {
		copaing.XP -= reduce
	}
	sql.Save(copaing)
	lvl := CalcLevel(copaing.XP)
	data := NewXpData{
		IsNewLevel: lvl != oldLvl,
		OldLevel:   oldLvl,
		NewLevel:   lvl,
		LevelUp:    oldLvl < lvl,
		LevelDown:  oldLvl > lvl,
	}
	return data
}

// NewXpNoUpdate calculate the new Xp without updating the last event and returns the data stored in NewXpData
func NewXpNoUpdate(member *discordgo.Member, copaing *sql.Copaing, exp uint) NewXpData {
	return NewXp(member, copaing, exp, false)
}
