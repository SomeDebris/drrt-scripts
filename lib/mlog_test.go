package lib

import (
	"testing"
)

const MLOG_FNAME = `MLOG_20250115_04.04.10.PM.txt`



func TestShipAuthorNamefmt(t *testing.T) {
	out := ShipAuthorFromCommonNamefmt("DSF Dark-Star Dreadnought [by Dukeslayer]")
	if out[0] != "DSF Dark-Star Dreadnought" {
		t.Errorf("Failed to ship name correctly: \"%s\", by \"%s\"", out[0], out[1])
	}
	if out[1] != "Dukeslayer" {
		t.Errorf("Failed to ship name correctly: \"%s\", by \"%s\"", out[0], out[1])
	}
}

func TestGetMatchNumberRedCorrectfmt(t *testing.T) {
	out := GetMatchNumberFromAllianceName("Match 001 - ^1The Red Alliance^7", false)
	target := 1
	if out != target {
		t.Errorf("Failed to get correct match number: %d should be %d", out, target)
	}
}
func TestGetMatchNumberRedPassedBlue(t *testing.T) {
	out := GetMatchNumberFromAllianceName("Match 001 - ^1The Red Alliance^7", true)
	target := 0
	if out != 0 {
		t.Errorf("Failed to return failure (0) value: %d should be %d", out, target)
	}
}
func TestGetMatchNumberBlueCorrectfmt(t *testing.T) {
	out := GetMatchNumberFromAllianceName("Match 001 - ^4The Blue Alliance^7", true)
	target := 0
	if out != 0 {
		t.Errorf("Failed to get correct match number: %d should be %d", out, target)
	}
}
func TestGetMatchNumberBluePassedRed(t *testing.T) {
	out := GetMatchNumberFromAllianceName("Match 001 - ^4The Blue Alliance^7", false)
	target := 0
	if out != 0 {
		t.Errorf("Failed to return failure (0) value: %d should be %d", out, target)
	}
}

// func TestRawMlog(t *testing.T) {
//
// }
