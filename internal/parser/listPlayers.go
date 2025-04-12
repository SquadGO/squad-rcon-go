package parser

import (
	"regexp"
	"strings"

	"github.com/SquadGO/squad-rcon-go/v2/rconEvents"
)

type Player struct {
	PlayerID   string
	EosID      string
	SteamID    string
	PlayerName string
	TeamID     string
	SquadID    string
	Role       string
	IsLeader   bool
	IsInSquad  bool
}

type Players []Player

func listPlayers(line, command string) (event string, data interface{}) {
	strs := strings.Split(line, "\n")
	players := make(Players, 0)

	if command == rconEvents.LIST_PLAYERS {
		for _, v := range strs {
			re := regexp.MustCompile(`ID: ([0-9]+) \| Online IDs: EOS: ([0-9a-f]{32}) steam: (\d{17}) \| Name: (.+) \| Team ID: ([0-9]+) \| Squad ID: ([0-9]+|N\/A) \| Is Leader: (True|False) \| Role: ([A-Za-z0-9_]*)\b`)
			matches := re.FindStringSubmatch(v)

			if matches == nil {
				continue
			}

			players = append(players, Player{
				PlayerID:   matches[1],
				EosID:      matches[2],
				SteamID:    matches[3],
				PlayerName: strings.TrimSpace(matches[4]),
				TeamID:     matches[5],
				SquadID:    matches[6],
				IsInSquad:  matches[6] != "N/A",
				IsLeader:   matches[7] == "True",
				Role:       matches[8],
			})
		}

		return rconEvents.LIST_PLAYERS, players
	}

	return rconEvents.LIST_PLAYERS, nil
}
