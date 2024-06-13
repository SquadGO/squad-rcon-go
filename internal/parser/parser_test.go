package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_warn(t *testing.T) {
	str := "Remote admin has warned player Player1. Message was \"Message from RCON\""
	expected := Warn{
		PlayerName: "Player1",
		Message:    "Message from RCON",
	}

	result := ChatParser(str)

	if diff := cmp.Diff(expected, result); diff != "" {
		t.Error(diff)
	}
}

func Test_kick(t *testing.T) {
	str := "Kicked player 1. [Online IDs= EOS: 00011111111111111111111111111111 steam: 11111111111111111] Player1"
	expected := Kick{
		PlayerID:   "1",
		EosID:      "00011111111111111111111111111111",
		SteamID:    "11111111111111111",
		PlayerName: "Player1",
	}

	result := ChatParser(str)

	if diff := cmp.Diff(expected, result); diff != "" {
		t.Error(diff)
	}
}

func Test_message(t *testing.T) {
	str := "[ChatAll] [Online IDs:EOS: 00011111111111111111111111111111 steam: 11111111111111111] Player1 : Message from RCON"
	expected := Message{
		ChatType:   "ChatAll",
		EosID:      "00011111111111111111111111111111",
		SteamID:    "11111111111111111",
		PlayerName: "Player1",
		Message:    "Message from RCON",
	}

	result := ChatParser(str)

	if diff := cmp.Diff(expected, result); diff != "" {
		t.Error(diff)
	}
}

func Test_squadCreated(t *testing.T) {
	str := "Player1 (Online IDs: EOS: 00011111111111111111111111111111 steam: 11111111111111111) has created Squad 1 (Squad Name: Squad 1) on United States Marine Corps"
	expected := SquadCreated{
		PlayerName: "Player1",
		EosID:      "00011111111111111111111111111111",
		SteamID:    "11111111111111111",
		SquadID:    "1",
		SquadName:  "Squad 1",
		TeamName:   "United States Marine Corps",
	}

	result := ChatParser(str)

	if diff := cmp.Diff(expected, result); diff != "" {
		t.Error(diff)
	}
}

func Test_posAdminCam(t *testing.T) {
	str := "[Online Ids:EOS: 00011111111111111111111111111111 steam: 11111111111111111] Player1 has possessed admin camera."
	expected := PosAdminCam{
		EosID:     "00011111111111111111111111111111",
		SteamID:   "11111111111111111",
		AdminName: "Player1",
	}

	result := ChatParser(str)

	if diff := cmp.Diff(expected, result); diff != "" {
		t.Error(diff)
	}
}

func Test_unposAdminCam(t *testing.T) {
	str := "[Online IDs:EOS: 00011111111111111111111111111111 steam: 11111111111111111] Player1 has unpossessed admin camera."
	expected := UnposAdminCam{
		EosID:     "00011111111111111111111111111111",
		SteamID:   "11111111111111111",
		AdminName: "Player1",
	}

	result := ChatParser(str)

	if diff := cmp.Diff(expected, result); diff != "" {
		t.Error(diff)
	}
}
