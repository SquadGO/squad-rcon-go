package parser

import (
	"regexp"
	"strings"

	"github.com/SquadGO/squad-rcon-go/v2/rconEvents"
)

type SquadCreated struct {
	Raw        string
	PlayerName string
	EosID      string
	SteamID    string
	SquadID    string
	SquadName  string
	TeamName   string
}

func squadCreated(line string) (event string, data interface{}) {
	var re *regexp.Regexp
	var matches []string

	re = regexp.MustCompile(`(.+) \(Online IDs: EOS: ([0-9a-f]{32}) steam: (\d{17})\) has created Squad (\d+) \(Squad Name: (.+)\) on (.+)`)
	matches = re.FindStringSubmatch(line)

	if matches != nil {
		return rconEvents.SQUAD_CREATED, SquadCreated{
			Raw:        line,
			PlayerName: strings.TrimSpace(matches[1]),
			EosID:      matches[2],
			SteamID:    matches[3],
			SquadID:    matches[4],
			SquadName:  matches[5],
			TeamName:   matches[6],
		}
	}

	return rconEvents.SQUAD_CREATED, nil
}
