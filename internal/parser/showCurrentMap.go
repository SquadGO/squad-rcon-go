package parser

import (
	"regexp"
	"strings"

	"github.com/SquadGO/squad-rcon-go/v2/rconEvents"
)

type CurrentMap struct {
	Raw      string
	Level    string
	Layer    string
	Factions []string
}

func showCurrentMap(line, command string) (event string, data interface{}) {
	if command == rconEvents.SHOW_CURRENT_MAP {
		var re *regexp.Regexp
		var matches []string

		re = regexp.MustCompile(`^Current level is (.*), layer is (.*), factions (.*)`)
		matches = re.FindStringSubmatch(line)

		if matches != nil {
			return rconEvents.SHOW_CURRENT_MAP, CurrentMap{
				Raw:      line,
				Level:    matches[1],
				Layer:    matches[2],
				Factions: strings.Split(matches[3], " "),
			}
		}
	}

	return rconEvents.SHOW_CURRENT_MAP, nil
}
