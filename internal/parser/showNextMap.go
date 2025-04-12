package parser

import (
	"fmt"
	"regexp"

	"github.com/SquadGO/squad-rcon-go/v2/rconEvents"
)

type NextMap struct {
	Raw      string
	Level    string
	Layer    string
	Factions []string
}

func showNextMap(line, command string) (event string, data interface{}) {
	if command == rconEvents.SHOW_NEXT_MAP {
		var re *regexp.Regexp
		var matches []string

		re = regexp.MustCompile(`^Next level is (.*), layer is (.*)`)
		matches = re.FindStringSubmatch(line)

		fmt.Println(line)

		if matches != nil {
			return rconEvents.SHOW_NEXT_MAP, NextMap{
				Raw:   line,
				Level: matches[1],
				Layer: matches[2],
			}
		}
	}

	return rconEvents.SHOW_NEXT_MAP, nil
}
