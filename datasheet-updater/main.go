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
	"sync"
	// "cmp"
	"slices"

	// "golang.org/x/oauth2"
	// "golang.org/x/oauth2/google"
	// "google.golang.org/api/option"
	// "google.golang.org/api/sheets/v4"
)

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

	ships := make([]*rsmships.Ship, len(ship_paths))
	// unmarshal ship files
	var unmarshal_wait_group sync.WaitGroup
	lib.GoUnmarshalAllShipsFromPaths(&ships, fullshippaths, &unmarshal_wait_group)
	unmarshal_wait_group.Wait()
	
	// TODO: open named pipe
	// when I recieve a string on the named pipe, decide what I need to do.

// we now have enough information to put match logs in context.
}

