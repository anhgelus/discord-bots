package config

import (
	"fmt"
	"github.com/anhgelus/discord-bots/rp-helper/src/utils"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

type Objectives struct {
	Settings     Settings
	Mains        []Main
	Secondaries  []Secondary
	Placeholders []Placeholder
}

type Settings struct {
	Lang                string
	NumberOfSecondaries int
}

type Main struct {
	Name  string
	Goal1 string
	Goal2 string
}

type Secondary struct {
	Goal string
}

type Placeholder struct {
	Name string
	List []string
}

type numberPlaceholder struct {
	ID          int
	Placeholder string
	Number      int
}

type simplePlaceholder struct {
	ID          int
	Placeholder string
}

var (
	Objs  = Objectives{}
	mObjs = map[string][]string{}

	regSimple  = regexp.MustCompile("{[a-z_\\-0-9]+}")
	regNumbers = regexp.MustCompile("{[a-z_\\-0-9]+}\\[[0-9]+]")
)

const (
	enSep = "and"
	frSep = "et"

	UnsetGoal = "#-- UNSET --#"
)

func (s *Secondary) GenerateGoal() string {
	if !regSimple.MatchString(s.Goal) {
		return s.Goal
	}
	numbers, g := getNumberPlaceholders(s.Goal)
	simple, g := getSimplePlaceholders(g)
	for _, n := range numbers {
		g = n.Replace(g)
	}
	for _, s := range simple {
		g = s.Replace(g)
	}
	return g
}

func (m *Main) GenerateGoals() (string, string) {
	g1, g2 := Secondary{m.Goal1}, Secondary{m.Goal2}
	return g1.Goal, g2.Goal
}

func GenerateMainGoals(mains []Main) []string {
	var goals []string
	for _, m := range mains {
		g1, g2 := m.GenerateGoals()
		goals = append(goals, g1, g2)
	}
	return goals
}

func GenerateSecondaryGoals(secs []Secondary) []string {
	var goals []string
	for _, s := range secs {
		goals = append(goals, s.GenerateGoal())
	}
	return goals
}

func getNumberPlaceholders(s string) ([]numberPlaceholder, string) {
	bases := regNumbers.FindAllString(s, -1)
	var numbers []numberPlaceholder
	ns := s
	for i, b := range bases {
		sp := strings.Split(b, "}[")
		sp[0] = strings.ReplaceAll(sp[0], "{", "")
		sp[1] = strings.ReplaceAll(sp[1], "]", "")
		n, err := strconv.Atoi(sp[1])
		if err != nil {
			utils.SendAlert("objectives.go - Parsing number placeholder", err.Error())
			continue
		}
		numbers = append(numbers, numberPlaceholder{
			ID:          i,
			Placeholder: sp[0],
			Number:      n,
		})
		ns = strings.Replace(ns, b, fmt.Sprintf("$n{%d}", i), 1)
	}
	return numbers, ns
}

func getSimplePlaceholders(s string) ([]simplePlaceholder, string) {
	bases := regNumbers.FindAllString(s, -1)
	var numbers []simplePlaceholder
	ns := s
	for i, b := range bases {
		b = strings.ReplaceAll(b, "{", "")
		b = strings.ReplaceAll(b, "}", "")
		numbers = append(numbers, simplePlaceholder{
			ID:          i,
			Placeholder: b,
		})
		ns = strings.Replace(ns, b, fmt.Sprintf("$s{%d}", i), 1)
	}
	return numbers, ns
}

func getPlaceholder(name string) []string {
	if len(mObjs) == 0 {
		for _, p := range Objs.Placeholders {
			mObjs[p.Name] = p.List
		}
	}
	v, ok := mObjs[name]
	if !ok {
		utils.SendAlert("objectives.go - Getting placeholder", "impossible to get placeholder with the name "+name)
		return make([]string, 0)
	}
	return v
}

func (p *numberPlaceholder) Replace(s string) string {
	placeholder := getPlaceholder(p.Placeholder)
	if len(placeholder) < p.Number {
		utils.SendAlert("objectives.go - Replacing number placeholder", "too much items for "+p.Placeholder)
		return s
	}
	ns := ""
	for i := 1; i < p.Number; i++ {
		if i != p.Number-1 {
			ns += placeholder[rand.Intn(len(placeholder))] + ", "
		} else {
			switch Objs.Settings.Lang {
			case "en":
				ns += enSep + " " + placeholder[rand.Intn(len(placeholder))]
			case "fr":
				ns += frSep + " " + placeholder[rand.Intn(len(placeholder))]
			default:
				ns += enSep + " " + placeholder[rand.Intn(len(placeholder))]
			}
		}
	}
	return strings.ReplaceAll(s, fmt.Sprintf("$p{%d}", p.ID), ns)
}

func (p *simplePlaceholder) Replace(s string) string {
	placeholder := getPlaceholder(p.Placeholder)
	return strings.ReplaceAll(s, fmt.Sprintf("$p{%d}", p.ID), placeholder[rand.Intn(len(placeholder))])
}
