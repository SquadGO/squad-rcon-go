package parser

import (
	"regexp"
	"strings"

	"github.com/SquadGO/squad-rcon-go/v2/rconEvents"
	"github.com/SquadGO/squad-rcon-go/v2/rconTypes"
)

func showCurrentMap(line, command string) (event string, data interface{}) {
	if command == rconEvents.SHOW_CURRENT_MAP {
		var re *regexp.Regexp
		var matches []string

		re = regexp.MustCompile(`^Current level is (.*), layer is (.*), factions (.*)`)
		matches = re.FindStringSubmatch(line)

		if matches != nil {
			return rconEvents.SHOW_CURRENT_MAP, rconTypes.CurrentMap{
				Raw:      line,
				Level:    matches[1],
				Layer:    matches[2],
				Factions: strings.Split(matches[3], " "),
			}
		}
	}

	return rconEvents.SHOW_CURRENT_MAP, nil
}
