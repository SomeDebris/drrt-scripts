package main

import (
	// "context"
	// "encoding/json"
	// "fmt"
	"log"
	"log/slog"
	// "net/http"
	"flag"
	"os"
	// "bufio"
	"path/filepath"
	"drrt-scripts/lib"
	"github.com/SomeDebris/rsmships-go"
	"time"
	"sync"
	// "cmp"
	"errors"
	"slices"
	"regexp"

	// "golang.org/x/oauth2"
	// "golang.org/x/oauth2/google"
	// "google.golang.org/api/option"
	// "google.golang.org/api/sheets/v4"
)
type matchResult int
const (
	WinDestruction matchResult = iota
	WinPoints
	Loss
)


type matchPerformance struct {
	Match            uint
	Ship             *rsmships.Ship
	Destructions     uint
	RankPointsEarned uint
	Result           matchResult
	Survived         bool
}
func (m *matchPerformance) toSheetsRow() [][]any {
	output := make([]any, 8)
	output[0] = m.Ship.Data.Name
	output[1] = m.Match
	output[2] = m.Destructions
	output[3] = m.RankPointsEarned

	output[4] = 0
	output[5] = 0
	output[6] = 0
	switch m.Result {
	case WinDestruction:
		output[4] = 1
	case WinPoints:
		output[5] = 1
	case Loss:
		output[6] = 1
	}

	if m.Survived {
		output[7] = 1
	} else {
		output[7] = 0
	}
	return [][]any{output}
}

type DRRTStandardMatchLog struct {
	MatchNumber           int
	Timestamp             time.Time
	Ships                 []*rsmships.Ship
	AllianceLength        int
	Record                []*matchPerformance
	ShipIndices           []int
	PointsDamageInflicted []int
	PointsDamageTaken     []int
	Raw                   *lib.MatchLogRaw
}
// TODO: make the nametoidx variable an input argument
func NewDRRTStandardMatchLogFromShips(raw *lib.MatchLogRaw, ships []*rsmships.Ship) (*DRRTStandardMatchLog, error) {
	var mlog DRRTStandardMatchLog
	mlog.Raw = raw
	mlog.Timestamp = raw.CreatedTimestamp

	// collect array of shipidxs
	nametoidx := make(map[string]int)
	for i, ship := range ships {
		// NOTE: the ships' names must not use standard name format (name [by author])
		// +1 because the match schedule index starts at 1 and not 0. GO loops start at 0.
		nametoidx[ship.Data.Name] = i + 1
	}

	mlog.AllianceLength = len(raw.ShipListings) / 2
	// Get the indices of each ship participating in this match log
	// The first n ships are from the Red alliance, but may not be sorted in the order they appear in Reassembly's fleet screen.
	mlog.ShipIndices = make([]int, len(raw.ShipListings))
	mlog.Ships = make([]*rsmships.Ship, len(raw.ShipListings))
	for i, shiplsting := range raw.ShipListings {
		name := lib.ShipAuthorFromCommonNamefmt(shiplsting.Ship)[0]
		idx, ok := nametoidx[name]
		if !ok {
			slog.Warn("Ship index cannot be found using map.", "name", name)
			// TODO: you may want to return an error here.. but I don't know.
		}
		mlog.ShipIndices[i] = idx
		mlog.Ships[i] = ships[idx]
	}
	
	// ASSUMPTION: NOT a free-for-all (Red v Blue alliance)
	// get the match number
	// if the same for red and blue alliances: good!
	redMatchNumber := lib.GetMatchNumberFromAllianceName(raw.StartListings[0].Name, false)
	blueMatchNumber := lib.GetMatchNumberFromAllianceName(raw.StartListings[1].Name, true)
	if redMatchNumber != blueMatchNumber {
		return &mlog, errors.New("Red and Blue Alliance match numbers are different. Bad match log!")
	}
	mlog.MatchNumber = redMatchNumber // == blueMatchNumber

	// TODO next time: 
	// You need to assign the remaining fields to things.
	// Start by finding a way to generate MatchPerformances for ships (create a
	// hash map, I think).
	// Take the Points field from the ResultListings.
	// Create a function that generates a [][]any from this for passing into
	// google sheets.
	// Create a function for

	return &mlog, nil
}

// func getShipIndexFromName(name string, ships []*rsmships.Ship) int {
// 	for i, ship := range ships {
// 		ship_nameauthor := lib.ShipAuthorFromCommonNamefmt(ship.Data.Name)
// 		if name == 
// 	}
// 	return -1
// }
func getMatchesOnAlliance(shipidx int, schedule [][]int, isBlue bool) ([]int, []bool) {
	out := make([]int, 0)
	participated := make([]bool, len(schedule))
	n := len(schedule[0]) / 2
	for i, match := range schedule {
		contains := false
		if isBlue {
			contains = slices.Contains(match[n:(2*n-1)], shipidx)
		} else {
			contains = slices.Contains(match[0:(n-1)], shipidx)
		}
		if contains {
			out = append(out, i+1)
		}
		participated[i] = contains
	}
	return out, participated
}

// Checks whether the ship indices in shipidxs are indeed members of the same
// alliance in the Match Schedule.
// If the ships did play together 
func matchesPlayedTogether(shipidxs []int, schedule [][]int, isBlue bool) ([]int, bool) {
	// contains columns of bools determining whether the ship was a memeber of an alliance.
	isOnAlliance := make([][]bool, len(shipidxs))
	alliancePresence := make([][]int, len(shipidxs))
	for i, idx := range shipidxs {
		alliancePresence[i], isOnAlliance[i] = getMatchesOnAlliance(idx, schedule, isBlue)
	}

	// matches all ships played as members of the same alliance
	togetherMatches := make([]int, 0)
	
	// loop over the matches on the alliance played by ship 0
	for _, matchShip0 := range alliancePresence[0] {
		// loop over the remaining ships' alliance presences
		allMembers := true
		for _, presencesj := range alliancePresence[1:] {
			// if any of the other ships do not play in this match, set
			// allmembers to false
			if !slices.Contains(presencesj, matchShip0) {
				allMembers = false
			}
		}
		// if all members play in this match, append it to the list of matches
		// all ships play together
		if allMembers {
			togetherMatches = append(togetherMatches, matchShip0)
		}
	}
	// return the list of matches that all ships play together and whether there
	// are more than 0 matches that all ships play together.
	return togetherMatches, len(togetherMatches) > 0
}


func main() {
// boilerplate stuff for exit codes
	exit_code := 0
	defer os.Exit(exit_code)

	log_lvl := slog.LevelInfo

	drrt_directory_arg := flag.String("drrt-directory", "/home/magnus/Documents/reassembly_ships/tournaments/DRRT/2026 Winter DRRT", "Set the directory the DRRT will be run in.")
	log_file_name := flag.String("log-filename", "", "Send log messages to a file. If not set, log to standard error.")
	flag.Parse()
	ships_directory := filepath.Join(*drrt_directory_arg, "Ships")

	// same tech used in scheduler/main.go
	log_ref, log_writer_ref, err := lib.DRRTLoggerPreferences(*log_file_name, log_lvl)
	if err != nil {
		log.Fatalf("Could not open log file '%s': %v", *log_file_name, err)
	}
	logfile := *log_ref
	defer logfile.Close()
	if log_writer_ref != nil {
		logwriter := *log_writer_ref
		defer logwriter.Flush()
	}

	drrtdatasheet := lib.NewDRRTDatasheetDefaults()
	if drrtdatasheet.Service == nil {
		slog.Error("Failed to get google sheets service.", "err", err)
		exit_code = 1
		return
	}

	// TODO: use the lib.MatchSchedule type
	// get ship indices of the match schedule from google sheets. This is uploaded by the Scheduler.
	scheduleindices, err := drrtdatasheet.GetMatchSheduleValues()
	if err != nil {
		slog.Error("Failed to get match schdule.", "err", err)
		exit_code = 1
		return
	}

	// TODO: slap this in a function to reduce repeated code
	// get a slice comprising paths to all ships
	ship_paths, err := lib.GetJSONFilesSortedByModTime(ships_directory)
	if err != nil {
		slog.Error("Cannot get inspected ship paths.", "err", err)
		exit_code = 1
		return
	}
	// get full path of each ship
	fullshippaths := make([]string, len(ship_paths))
	for i, path := range ship_paths {
		fullshippaths[i] = filepath.Join(ships_directory, path)
	}
	slog.Info("Found paths for ship files.", "count", len(ship_paths))
	for i, path := range ship_paths {
		slog.Debug("Ship path", "path", path)
		slog.Debug("Full Ship path", "path", fullshippaths[i])
	}

	ships := make([]rsmships.Ship, len(ship_paths))
	// unmarshal ship files
	var unmarshal_wait_group sync.WaitGroup
	lib.GoUnmarshalAllShipsFromPaths(&ships, fullshippaths, &unmarshal_wait_group)
	unmarshal_wait_group.Wait()
	

// we now have enough information to put match logs in context.
}

