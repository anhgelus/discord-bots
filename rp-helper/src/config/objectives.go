package config

import (
	"fmt"
	"github.com/anhgelus/discord-bots/rp-helper/src/utils"
	"math/rand"
	"regexp"
	"slices"
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

func (s *Secondary) GenerateGoal(players []string) string {
	if !regSimple.MatchString(s.Goal) {
		return s.Goal
	}
	numbers, g := getNumberPlaceholders(s.Goal)
	simple, g := getSimplePlaceholders(g)
	for _, n := range numbers {
		g = n.Replace(g, players)
	}
	for _, s := range simple {
		g = s.Replace(g, players)
	}
	return g
}

func (m *Main) GenerateGoals(players []string) (string, string) {
	g1, g2 := Secondary{m.Goal1}, Secondary{m.Goal2}
	return g1.GenerateGoal(players), g2.GenerateGoal(players)
}

func GenerateMainGoals(mains []Main, players []string) []string {
	var goals []string
	for _, m := range mains {
		g1, g2 := m.GenerateGoals(players)
		goals = append(goals, g1, g2)
	}
	return goals
}

func GenerateSecondaryGoals(secs []Secondary, players []string) []string {
	var goals []string
	for _, s := range secs {
		goals = append(goals, s.GenerateGoal(players))
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
		ns = strings.Replace(ns, b, fmt.Sprintf("n{$%d}", i), 1)
	}
	return numbers, ns
}

func getSimplePlaceholders(s string) ([]simplePlaceholder, string) {
	bases := regSimple.FindAllString(s, -1)
	var numbers []simplePlaceholder
	ns := s
	for i, b := range bases {
		sp := strings.ReplaceAll(b, "{", "")
		sp = strings.ReplaceAll(sp, "}", "")
		numbers = append(numbers, simplePlaceholder{
			ID:          i,
			Placeholder: sp,
		})
		ns = strings.Replace(ns, b, fmt.Sprintf("s{$%d}", i), 1)
	}
	return numbers, ns
}

func getPlaceholder(name string, players []string) []string {
	if name == "player" {
		return players
	}
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

func (p *numberPlaceholder) Replace(s string, players []string) string {
	placeholder := getPlaceholder(p.Placeholder, players)
	if len(placeholder) < p.Number {
		utils.SendAlert("objectives.go - Replacing number placeholder", "too much items for "+p.Placeholder)
		return s
	}
	ns := ""
	var randValues []int
	for i := 1; i < p.Number; i++ {
		ns += placeholderRandValue(placeholder, &randValues)
		if i != p.Number-1 {
			ns += ", "
		} else {
			var sep string
			switch Objs.Settings.Lang {
			case "en":
				sep = enSep
			case "fr":
				sep = frSep
			default:
				sep = enSep
			}
			ns += fmt.Sprintf(" %s %s", sep, placeholderRandValue(placeholder, &randValues))
		}
	}
	return strings.ReplaceAll(s, fmt.Sprintf("n{$%d}", p.ID), ns)
}

func (p *simplePlaceholder) Replace(s string, players []string) string {
	placeholder := getPlaceholder(p.Placeholder, players)
	return strings.ReplaceAll(s, fmt.Sprintf("s{$%d}", p.ID), placeholder[rand.Intn(len(placeholder))])
}

func placeholderRandValue(placeholder []string, values *[]int) string {
	n := -1
	for n == -1 || slices.Contains(*values, n) {
		n = rand.Intn(len(placeholder))
	}
	*values = append(*values, n)
	return placeholder[n]
}
