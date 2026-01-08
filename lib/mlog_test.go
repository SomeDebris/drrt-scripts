package lib

import (
	"testing"
	"time"
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

func TestGetTimeOfMlogFname(t *testing.T) {
	correcttime := time.Date(2025, 1, 15, 16, 4, 10, 0, time.Local)
	outtime, err := GetTimeOfMatchLogFilename(MLOG_FNAME)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}
	switch correcttime.Compare(outtime) {
	case -1:
		t.Errorf("Failed to parse time: Correct time (%s) is before time specified in %s (%s).", correcttime.String(), MLOG_FNAME, outtime.String())
	case 0:
		return
	case 1:
		t.Errorf("Failed to parse time: Correct time (%s) is after time specified in %s (%s).", correcttime.String(), MLOG_FNAME, outtime.String())
	default:
		t.Errorf("Failed to parse time: Correct time (%s) was not compared in the expected way to the time specified in %s (%s).", correcttime.String(), MLOG_FNAME, outtime.String())
	}
}

func TestParseMlog(t *testing.T) {
	_, err := NewMatchLogRawFromPath(MLOG_FNAME)
	if err != nil {
		t.Errorf("Encountered error: %v", err)
	}
	expectedmlog := MatchLogRaw{
		CreatedTimestamp: time.Date(2025, 1, 15, 16, 4, 10, 0, time.Local),
		Path: MLOG_FNAME,
		StartListings: []MatchLogFleetListing{
			MatchLogFleetListing{Faction:100, Name:"Match 001 - ^1The Red Alliance^7", DamageTaken:0, DamageInflicted:0, Alive:3},
			MatchLogFleetListing{Faction:101, Name:"Match 001 - ^4The Blue Alliance^7", DamageTaken:0, DamageInflicted:0, Alive:3},
		},
		ShipListings: []MatchLogShipListing{
			MatchLogShipListing{Fleet: 100, Ship:"Transcription 2025 [by joyous eighteen]"},
			MatchLogShipListing{Fleet: 100, Ship:"Original Thinker [by MonsPubis]"},
			MatchLogShipListing{Fleet: 100, Ship:`Muninn M6-B "LAIKA" [by Infamous YenYu]`},
			MatchLogShipListing{Fleet: 101, Ship:`Spawk [by Splinter]`},
			MatchLogShipListing{Fleet: 101, Ship:`Lethal K v3 [by Splinter]`},
			MatchLogShipListing{Fleet: 101, Ship:`directional dismisser [by 1836 Nokia Mustang (CharredSkies)]`},
		},
		DestructionListings: []MatchLogDestructionListing{
			MatchLogDestructionListing{Ship:`Muninn M6-B "LAIKA" [by Infamous YenYu]`, Fship:100, Destroyed:`Lethal K v3 [by Splinter]`, Fdestroyed:101},
			MatchLogDestructionListing{Ship:`Muninn M6-B "LAIKA" [by Infamous YenYu]`, Fship:100, Destroyed:`Spawk [by Splinter]`, Fdestroyed:101},
			MatchLogDestructionListing{Ship:`Muninn M6-B "LAIKA" [by Infamous YenYu]`, Fship:100, Destroyed:`directional dismisser [by 1836 Nokia Mustang (CharredSkies)]`, Fdestroyed:101},
		},
		ResultListings: []MatchLogFleetListing{
			MatchLogFleetListing{Faction:100, Name:"Match 001 - ^1The Red Alliance^7", DamageTaken:200756, DamageInflicted:185891, Alive:3},
			MatchLogFleetListing{Faction:101, Name:"Match 001 - ^4The Blue Alliance^7", DamageTaken:0, DamageInflicted:51176, Alive:0},
		},
		SurvivalListings: []MatchLogShipListing{
			MatchLogShipListing{Fleet: 100, Ship:"Transcription 2025 [by joyous eighteen]"},
			MatchLogShipListing{Fleet: 100, Ship:"Original Thinker [by MonsPubis]"},
			MatchLogShipListing{Fleet: 100, Ship:`Muninn M6-B "LAIKA" [by Infamous YenYu]`},
		},
	}
}

// func TestRawMlog(t *testing.T) {
//
// }
