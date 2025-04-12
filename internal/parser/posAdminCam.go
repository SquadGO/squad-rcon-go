package parser

import (
	"regexp"
	"strings"

	"github.com/SquadGO/squad-rcon-go/v2/rconEvents"
)

type PosAdminCam struct {
	Raw       string
	EosID     string
	SteamID   string
	AdminName string
}

func posAdminCam(line string) (event string, data interface{}) {
	var re *regexp.Regexp
	var matches []string

	re = regexp.MustCompile(`\[Online Ids:EOS: ([0-9a-f]{32}) steam: (\d{17})\] (.+) has possessed admin camera\.`)
	matches = re.FindStringSubmatch(line)

	if matches != nil {
		return rconEvents.POSSESSED_ADMIN_CAMERA, PosAdminCam{
			Raw:       line,
			EosID:     matches[1],
			SteamID:   matches[2],
			AdminName: strings.TrimSpace(matches[3]),
		}
	}

	return rconEvents.POSSESSED_ADMIN_CAMERA, nil
}
