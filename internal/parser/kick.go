package parser

import (
	"regexp"
	"strings"

	"github.com/SquadGO/squad-rcon-go/v2/rconEvents"
)

type Kick struct {
	Raw        string
	PlayerID   string
	EosID      string
	SteamID    string
	PlayerName string
}

func kick(line string) (event string, data interface{}) {
	var re *regexp.Regexp
	var matches []string

	re = regexp.MustCompile(`Kicked player ([0-9]+)\. \[Online IDs= EOS: ([0-9a-f]{32}) steam: (\d{17})] (.*)`)
	matches = re.FindStringSubmatch(line)

	if matches != nil {
		return rconEvents.PLAYER_KICKED, Kick{
			Raw:        line,
			PlayerID:   matches[1],
			EosID:      matches[2],
			SteamID:    matches[3],
			PlayerName: strings.TrimSpace(matches[4]),
		}
	}

	return rconEvents.PLAYER_KICKED, nil
}
