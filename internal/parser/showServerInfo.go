package parser

import (
	"encoding/json"

	"github.com/SquadGO/squad-rcon-go/v2/rconEvents"
)

type ServerInfo struct {
	Raw              string
	ServerName       string `json:"ServerName_s"`
	MaxPlayers       int
	PublicQueueLimit int    `json:"PublicQueueLimit_I,string"`
	ReserveSlots     int    `json:"PlayerReserveCount_I,string"`
	PlayerCount      int    `json:"PlayerCount_I,string"`
	PublicQueue      int    `json:"PublicQueue_I,string"`
	ReserveQueue     int    `json:"ReservedQueue_I,string"`
	MatchTimeout     int    `json:"MatchTimeout_d"`
	MatchStartTime   int    `json:"PLAYTIME_I,string"`
	CurrentLayer     string `json:"MapName_s"`
	NextLayer        string `json:"NextLayer_s"`
	TeamOne          string `json:"TeamOne_s"`
	TeamTwo          string `json:"TeamTwo_s"`
	GameMode         string `json:"GameMode_s"`
	GameVersion      string `json:"GameVersion_s"`
}

func showServerInfo(line, command string) (event string, data interface{}) {
	if command == rconEvents.SHOW_SERVER_INFO {
		if len(line) > 10 {
			info := ServerInfo{}

			err := json.Unmarshal([]byte(line), &info)
			if err != nil {
				return rconEvents.SHOW_SERVER_INFO, nil
			}

			return rconEvents.SHOW_SERVER_INFO, info
		}
	}

	return rconEvents.SHOW_SERVER_INFO, nil
}
