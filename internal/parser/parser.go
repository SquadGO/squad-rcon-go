package parser

import (
	"github.com/SquadGO/squad-rcon-go/v2/rconEvents"
	"github.com/iamalone98/eventEmitter"
)

type chatParser func(string) (event string, data interface{})
type commandParser func(string, string) (event string, data interface{})

var parsers = []chatParser{
	ban,
	kick,
	message,
	posAdminCam,
	unposAdminCam,
	squadCreated,
	warn,
}
var commandParsers = []commandParser{
	listPlayers,
	listSquads,
	showCurrentMap,
	showNextMap,
	showServerInfo,
}

func RconParser(line, command string, emitter eventEmitter.EventEmitter) {
	if len(command) > 0 {
		for _, fn := range commandParsers {
			event, data := fn(line, command)

			if data != nil {
				emitter.Emit(rconEvents.DATA, data)
				emitter.Emit(event, data)
				break
			}
		}
	} else {
		for _, fn := range parsers {
			event, data := fn(line)

			if data != nil {
				emitter.Emit(rconEvents.DATA, data)
				emitter.Emit(event, data)
				break
			}
		}
	}

}
