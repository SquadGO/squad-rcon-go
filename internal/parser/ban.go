package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/SquadGO/squad-rcon-go/v2/rconEvents"
)

type Ban struct {
	Raw        string
	PlayerID   string
	SteamID    string
	PlayerName string
	Interval   int
}

func ban(line string) (event string, data interface{}) {
	var re *regexp.Regexp
	var matches []string

	re = regexp.MustCompile(`Banned player ([0-9]+)\. \[steamid=(.*?)\] (.*) for interval (.*)`)
	matches = re.FindStringSubmatch(line)

	if matches != nil {
		interval, err := strconv.Atoi(matches[4])
		if err != nil {
			return rconEvents.PLAYER_BANNED, nil
		}

		return rconEvents.PLAYER_BANNED, Ban{
			Raw:        line,
			PlayerID:   matches[1],
			SteamID:    matches[2],
			PlayerName: strings.TrimSpace(matches[3]),
			Interval:   interval,
		}
	}

	return rconEvents.PLAYER_BANNED, nil
}
