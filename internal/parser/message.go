package parser

import (
	"regexp"
	"strings"

	"github.com/SquadGO/squad-rcon-go/v2/rconEvents"
)

type Message struct {
	Raw        string
	ChatType   string
	EosID      string
	SteamID    string
	PlayerName string
	Message    string
}

func message(line string) (event string, data interface{}) {
	var re *regexp.Regexp
	var matches []string

	re = regexp.MustCompile(`\[(ChatAll|ChatTeam|ChatSquad|ChatAdmin)] \[Online IDs:EOS: ([0-9a-f]{32}) steam: (\d{17})\] (.+?) : (.*)`)
	matches = re.FindStringSubmatch(line)

	if matches != nil {
		return rconEvents.CHAT_MESSAGE, Message{
			Raw:        line,
			ChatType:   matches[1],
			EosID:      matches[2],
			SteamID:    matches[3],
			PlayerName: strings.TrimSpace(matches[4]),
			Message:    matches[5],
		}
	}

	return rconEvents.CHAT_MESSAGE, nil
}
