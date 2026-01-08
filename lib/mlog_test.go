package lib

import (
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/SomeDebris/rsmships-go"
)

const MLOG_FNAME = `MLOG_20250115_04.04.10.PM.txt`
var shiptestdir = filepath.Join("..", "test-ships")



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
	mlog, err := NewMatchLogRawFromPath(MLOG_FNAME)
	if err != nil {
		t.Logf("Encountered error: %v", err)
		t.FailNow()
	}
	expectedmlog := MatchLogRaw{
		CreatedTimestamp: time.Date(2025, 1, 15, 16, 4, 10, 0, time.Local),
		Path: MLOG_FNAME,
		StartListings: []MatchLogFleetListing{
			{Faction:100, Name:"Match 001 - ^1The Red Alliance^7", DamageTaken:0, DamageInflicted:0, Alive:3},
			{Faction:101, Name:"Match 001 - ^4The Blue Alliance^7", DamageTaken:0, DamageInflicted:0, Alive:3},
		},
		ShipListings: []MatchLogShipListing{
			{Fleet: 100, Ship:"Transcription 2025 [by joyous eighteen]"},
			{Fleet: 100, Ship:"Original Thinker [by MonsPubis]"},
			{Fleet: 100, Ship:`Muninn M6-B "LAIKA" [by Infamous YenYu]`},
			{Fleet: 101, Ship:`Spawk [by Splinter]`},
			{Fleet: 101, Ship:`Lethal K v3 [by Splinter]`},
			{Fleet: 101, Ship:`directional dismisser [by 1836 Nokia Mustang (CharredSkies)]`},
		},
		DestructionListings: []MatchLogDestructionListing{
			{Ship:`Muninn M6-B "LAIKA" [by Infamous YenYu]`, Fship:100, Destroyed:`Lethal K v3 [by Splinter]`, Fdestroyed:101},
			{Ship:`Muninn M6-B "LAIKA" [by Infamous YenYu]`, Fship:100, Destroyed:`Spawk [by Splinter]`, Fdestroyed:101},
			{Ship:`Muninn M6-B "LAIKA" [by Infamous YenYu]`, Fship:100, Destroyed:`directional dismisser [by 1836 Nokia Mustang (CharredSkies)]`, Fdestroyed:101},
		},
		ResultListings: []MatchLogFleetListing{
			{Faction:100, Name:"Match 001 - ^1The Red Alliance^7", DamageTaken:200756, DamageInflicted:185891, Alive:3},
			{Faction:101, Name:"Match 001 - ^4The Blue Alliance^7", DamageTaken:0, DamageInflicted:51176, Alive:0},
		},
		SurvivalListings: []MatchLogShipListing{
			{Fleet: 100, Ship:"Transcription 2025 [by joyous eighteen]"},
			{Fleet: 100, Ship:"Original Thinker [by MonsPubis]"},
			{Fleet: 100, Ship:`Muninn M6-B "LAIKA" [by Infamous YenYu]`},
		},
	}
	if !reflect.DeepEqual(*mlog, expectedmlog) {
		t.Errorf("Parsed match log is not identical to expectation: expected ```%v```, got ```%v```", expectedmlog, *mlog)
	}
	if !reflect.DeepEqual((*mlog).CreatedTimestamp, expectedmlog.CreatedTimestamp) {
		t.Errorf("Parsed match log Timestamp is not identical to expectation: expected `%s`, got `%s`", expectedmlog.CreatedTimestamp.String(), mlog.CreatedTimestamp.String())
	}
	if !reflect.DeepEqual((*mlog).Path, expectedmlog.Path) {
		t.Errorf("Parsed match log path is not identical to expectation: expected `%s`, got `%s`", expectedmlog.Path, mlog.Path)
	}

	if !reflect.DeepEqual((*mlog).ShipListings, expectedmlog.ShipListings) {
		t.Errorf("Parsed match log ship listing is not identical to expectation: expected `%v`, got `%v`", expectedmlog.ShipListings, mlog.ShipListings)
	}
	if !reflect.DeepEqual((*mlog).DestructionListings, expectedmlog.DestructionListings) {
		t.Errorf("Parsed match log DestructionListings is not identical to expectation: expected `%v`, got `%v`", expectedmlog.DestructionListings, mlog.DestructionListings)
	}
	if !reflect.DeepEqual((*mlog).ResultListings, expectedmlog.ResultListings) {
		t.Errorf("Parsed match log ResultListings is not identical to expectation: expected `%v`, got `%v`", expectedmlog.ResultListings, mlog.ResultListings)
	}
	if !reflect.DeepEqual((*mlog).SurvivalListings, expectedmlog.SurvivalListings) {
		t.Errorf("Parsed match log SurvivalListings is not identical to expectation: expected `%v`, got `%v`", expectedmlog.SurvivalListings, mlog.SurvivalListings)
	}
}


func TestNewDRRTStandardMatchLog(t *testing.T) {
	ship_paths, err := GetJSONFilesSortedByModTime(shiptestdir)
	if err != nil {
		t.Logf("Cannot get inspected ship paths: %v", err)
		t.FailNow()
	}
	ships := make([]*rsmships.Ship, len(ship_paths))

	// unmarshal ship files
	var unmarshal_wait_group sync.WaitGroup
	GoUnmarshalAllShipsFromPaths(&ships, ship_paths, &unmarshal_wait_group)
	unmarshal_wait_group.Wait()

	raw, err := NewMatchLogRawFromPath(MLOG_FNAME)
	if err != nil {
		t.Logf("Encountered error parsing matchlog: %v", err)
		t.FailNow()
	}
	
	idxfac := getShipIdxFacMap(ships)
	_, err = NewDRRTStandardMatchLogFromShips(raw, ships, idxfac)
	if err != nil {
		t.Errorf("Encountered error while producing match log object: %v", err)
	}
}

// func TestRawMlog(t *testing.T) {
//
// }
