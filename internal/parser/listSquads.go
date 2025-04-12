package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/SquadGO/squad-rcon-go/v2/rconEvents"
)

type Squad struct {
	SquadID        string
	SquadName      string
	Size           string
	CreatorName    string
	CreatorEOSID   string
	CreatorSteamID string
	TeamID         int
	TeamName       string
	Locked         bool
}

type Squads []Squad

func listSquads(line, command string) (event string, data interface{}) {
	strs := strings.Split(line, "\n")
	squads := make(Squads, 0)
	teamID := 0
	teamName := ""

	if command == rconEvents.LIST_SQUADS {
		for _, v := range strs {
			re := regexp.MustCompile(`ID: ([0-9]+) \| Name: (.+) \| Size: ([0-9]+) \| Locked: (True|False) \| Creator Name: (.+) \| Creator Online IDs: EOS: ([0-9a-f]{32}) steam: (\d{17})`)
			matches := re.FindStringSubmatch(v)

			teamMatches := regexp.MustCompile(`Team ID: (1|2) \((.+)\)/`).FindStringSubmatch(v)

			if teamMatches != nil {
				if id, err := strconv.Atoi(teamMatches[1]); err == nil {
					teamID = id
				}

				teamName = teamMatches[2]
			}

			if matches == nil {
				continue
			}

			squads = append(squads, Squad{
				SquadID:        matches[1],
				SquadName:      matches[2],
				Size:           matches[3],
				Locked:         matches[4] == "True",
				CreatorName:    matches[5],
				CreatorEOSID:   matches[6],
				CreatorSteamID: matches[7],
				TeamID:         teamID,
				TeamName:       teamName,
			})
		}

		return rconEvents.LIST_SQUADS, squads
	}

	return rconEvents.LIST_SQUADS, nil
}
