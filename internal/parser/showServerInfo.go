package parser

import (
	"encoding/json"

	"github.com/SquadGO/squad-rcon-go/v2/rconEvents"
	"github.com/SquadGO/squad-rcon-go/v2/rconTypes"
)

func showServerInfo(line, command string) (event string, data interface{}) {
	if command == rconEvents.SHOW_SERVER_INFO {
		if len(line) > 10 {
			info := rconTypes.ServerInfo{}

			err := json.Unmarshal([]byte(line), &info)
			if err != nil {
				return rconEvents.SHOW_SERVER_INFO, nil
			}

			return rconEvents.SHOW_SERVER_INFO, info
		}
	}

	return rconEvents.SHOW_SERVER_INFO, nil
}
