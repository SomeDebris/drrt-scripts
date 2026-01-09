package main

import (
	// "context"
	// "encoding/json"
	// "fmt"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"strconv"

	// "net/http"
	"flag"
	"os"

	// "bufio"
	"drrt-scripts/lib"
	"path/filepath"
	"sync"

	"github.com/SomeDebris/rsmships-go"

	// "cmp"
	"slices"
	// "golang.org/x/oauth2"
	// "golang.org/x/oauth2/google"
	// "google.golang.org/api/option"
	// "google.golang.org/api/sheets/v4"
)

const (
	pipecmd_reload = `reload`
	pipecmd_stop   = `stop`
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
	pipewatch_arg := flag.String("pipe", lib.DRRT_MLOG_SIGNAL_PIPE_PATH, "Send log messages to a file. If not set, log to standard error.")
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
		logwriter := log_writer_ref
		fmt.Printf("Yeah, uh, you are flushed.")
		defer (*logwriter).Flush()
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
//
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

	nametoidx := lib.GetShipIdxFacMap(ships)
	// mlogs, err := updateMatchLogs(drrtdatasheet, ships, nametoidx)
	// if err != nil {
	// 	slog.Error("Failed to update match logs for the first time.", "err", err)
	// }
		

	var mlogs []*lib.DRRTStandardMatchLog
	ranks := make(map[string]int)
	
	// ASSUMPTION: you're running quals
	// and you've already run the scheduler
	OuterLoop:
	for {
		fmt.Println(lib.ANSI_BOLD + lib.ANSI_OKGREEN + "Waiting on pipe." + lib.ANSI_RESET)
		data, err := lib.ReadDRRTMlogPipe(*pipewatch_arg)
		if err != nil {
			slog.Error("Encountered error reading mlog signal pipe.", "err", err)
			exit_code = 1
			return
		}
		fmt.Println(lib.ANSI_BOLD + lib.ANSI_OKGREEN + "Got command: \"" + data + "\"" + lib.ANSI_RESET)
		switch data {
		case pipecmd_reload:
			mlogs, err = updateMatchLogs(drrtdatasheet, ships, nametoidx)
			if err != nil {
				slog.Error("Failed to update match logs.", "err", err)
			}
			ranks, err = drrtdatasheet.GetRanks()
			if err != nil {
				slog.Error("Failed to get ranks.", "err", err)
			}
			lib.UpdateNextUp(ships, shipidxsfromSchedule(mlogs[len(mlogs)-1].MatchNumber, scheduleindices), mlogs, ranks)
			lib.UpdateGame(ships, shipidxsfromSchedule(mlogs[len(mlogs)-1].MatchNumber, scheduleindices), mlogs, ranks)
			lib.UpdateVictory(ships, mlogs, ranks)
		case pipecmd_stop:
			fmt.Println(lib.ANSI_BOLD + lib.ANSI_WHITE + "Stopping." + lib.ANSI_RESET)
			break OuterLoop
		default:
			fmt.Println(lib.ANSI_BOLD + lib.ANSI_FAIL + "Command \"" + data + "\" is not known." + lib.ANSI_RESET)
		}
		
	}

	slog.Info("Done!")
}

// Returns the latest list of match logs and uploads them to the sheet
func updateMatchLogs(sheet *lib.DRRTDatasheet, ships []*rsmships.Ship, nametoidx *map[string]int) ([]*lib.DRRTStandardMatchLog, error) {
	mlogsraw, err := lib.ReadMlogRawsFromPath("/home/magnus/.local/share/Reassembly/data")
	if err != nil {
		slog.Error("Error while reading match logs.", "err", err)
		return nil, err
	}
	// get DRRTStandardMatchLogs
	var stdMatchLogs []*lib.DRRTStandardMatchLog
	for _, mlograw := range mlogsraw {
		mlogparsed, err := lib.NewDRRTStandardMatchLogFromShips(mlograw, ships, nametoidx)
		if err != nil {
			// check the error type. Many of these do not cause problems; they just should be ignored.:
			var mlogincomplete *lib.MatchLogIncompleteError
			if errors.As(err, &mlogincomplete) {
				mlogincomplete.LogError(slog.Default())
				continue
			}
			var mlogbadmatchnumbers *lib.MatchLogAllianceMatchNumberMismatchError
			if errors.As(err, &mlogbadmatchnumbers) {
				mlogbadmatchnumbers.LogError(slog.Default())
				continue
			}
			var mlogbadlength *lib.MatchLogAllianceLengthMismatchError
			if errors.As(err, &mlogbadlength) {
				mlogbadlength.LogError(slog.Default())
				continue
			}
			slog.Error("Error while scoring match log.", "err", err, "path", mlograw.Path)
			continue
		}
		slog.Info("Successfully scored match log.", "path", mlograw.Path, "matchNumber", mlogparsed.MatchNumber)
		stdMatchLogs = append(stdMatchLogs, mlogparsed)
	}
	// clean up match logs
	lib.SortStandardMlogs(&stdMatchLogs)
	lib.DeleteDuplicateMlogs(&stdMatchLogs)
	_, err = sheet.UpdateMatchLogs(stdMatchLogs)
	if err != nil {
		slog.Error("Error updating match logs in spreadsheet.", "err", err)
		return stdMatchLogs, err
	}

	return stdMatchLogs, nil
}

func shipidxsfromSchedule(idx int, schedule [][]any) []int {
	out := make([]int, len(schedule[idx]))
	for i, val := range schedule[idx] {
		pstr, ok := val.(string)
		if !ok {
			slog.Error("Rank from sheets cannot be interpreted as string.", "pstr", pstr)
			out[i] = 0
		}
		p, err := strconv.Atoi(pstr)
		if err != nil {
			slog.Error("Rank from sheets cannot be parsed INTO an integer.", "err", err)
			continue
		}
		out[i] = p - 1
	}
	return out
}
