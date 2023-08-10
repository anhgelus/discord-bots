package xp

import (
	"github.com/anhgelus/discord-bots/les-copaings/src/db/redis"
	"github.com/anhgelus/discord-bots/les-copaings/src/db/sql"
	"github.com/anhgelus/discord-bots/les-copaings/src/utils"
	"github.com/bwmarrin/discordgo"
	"math"
)

func CalcExperience(length uint, diversity uint) uint {
	// f(x;y) = ((0.025 x^{1.25})/(y^{-0.5}))+1
	return uint(math.Floor(((0.025 * math.Pow(float64(length), 1.25)) / math.Pow(float64(diversity), -0.5)) + 1))
}

func CalcExperienceFromVocal(length uint) uint {
	// f(x)=((0.25 x^{1.3})/(60000))+1
	return uint(math.Floor(((0.25 * math.Pow(float64(length), 1.3)) / 60000) + 1))
}

func CalcLevel(xp uint) uint {
	// f(x)=0.2 x^0.5
	return uint(math.Floor(0.2 * math.Pow(float64(xp), 0.5)))
}

func CalcXpForLevel(level uint) uint {
	// f(x)=0.2 x^0.5
	// f(x)/0.2 = x^0.5
	// (f(x)/0.2)^2 = x
	return uint(math.Floor(math.Pow(20*float64(level), 2)))
}

func CalcXpLose(inactivity uint) uint {
	// f(x)= 0.01x^2
	return uint(math.Floor(0.01 * math.Pow(float64(inactivity), 2)))
}

type NewXpData struct {
	IsNewLevel bool
	OldLevel   uint
	NewLevel   uint
	LevelUp    bool
	LevelDown  bool
}

func NewXp(member *discordgo.Member, copaing *sql.Copaing, exp uint) NewXpData {
	user := redis.GenerateConnectedUser(member)
	time := user.TimeSinceLastEvent()
	reduce := CalcXpLose(utils.HoursOfUnix(time))
	user.UpdateLastEvent()

	oldLvl := CalcLevel(copaing.XP)
	if int(copaing.XP)-int(reduce) < 0 {
		copaing.XP = 0
	} else {
		copaing.XP -= reduce
	}
	copaing.XP += exp
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
