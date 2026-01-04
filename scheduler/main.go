package main

import (
	"drrt-scripts/lib"
	"encoding/csv"
	"flag"
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"github.com/SomeDebris/rsmships-go"
	"errors"
)

const (
	VERSION                       = "0.0.0"
	PROGRAM_NAME                  = "drrt-scheduler"
	SELECTED_SCHEDULE_FNAME       = "selected_schedule.csv"
	SELECTED_SCHEDULE_NOAST_FNAME = ".no_asterisks.csv"
)

type ShipDataError struct{}
type MultipleShipsInFleetError struct{}

func (m *ShipDataError) Error() string {
	return "Ship data not found in file or formatted incorrectly."
}

func (m *MultipleShipsInFleetError) Error() string {
	return "Fleet file has multiple ships defined. Only one ship should be defined in the file."
}

func get_inspected_ship_paths(dir string) ([]string, error) {
	var ship_files []string

	f, err := os.Open(dir)
	if err != nil {
		return ship_files, err
	}
	defer f.Close()

	file_info, err := f.Readdir(-1)
	if err != nil {
		return ship_files, err
	}

	sort.Slice(file_info, func(i, j int) bool {
		return file_info[i].ModTime().Before(file_info[j].ModTime())
	})

	for _, file := range file_info {
		if !strings.EqualFold(filepath.Ext(file.Name()), ".json") {
			continue
		}
		ship_files = append(ship_files, file.Name())
	}

	return ship_files, nil
}


/**
* schedule array
* whether ship is participating as surrogate
* any error recieved
 */
func readScheduleAtPath(path string) (MatchSchedule, [][]string, error) {
	schedule_bytes, err := os.ReadFile(path)
	if err != nil {
		return MatchSchedule{}, nil, err
	}

	schedule_string := string(schedule_bytes)

	r := csv.NewReader(strings.NewReader(schedule_string))

	records, err := r.ReadAll()
	if err != nil {
		return MatchSchedule{}, nil, err
	}

	schedule := make([][]int, len(records))
	surrogates := make([][]bool, len(records))

	// Parse out the strings in the schedules to integers
	for j, match := range records {
		ships_in_match := make([]int, len(match))
		surrogates_in_match := make([]bool, len(match))

		for i, ship := range match {
			if strings.ContainsAny(ship, "*") {
				surrogates_in_match[i] = true
			} else {
				surrogates_in_match[i] = false
			}

			ship_noasterisk := strings.ReplaceAll(ship, "*", "")

			ships_in_match[i], err = strconv.Atoi(ship_noasterisk)
			if err != nil {
				return MatchSchedule{}, records, err
			}
		}
		schedule[j] = ships_in_match
		surrogates[j] = surrogates_in_match
	}
	return MatchSchedule{
		Schedule:       schedule,
		Surrogates:     surrogates,
		Length:         len(records),
		AllianceLength: len(records[0]),
	}, records, nil
}

func int2dSliceToString(ints [][]int) ([][]string, error) {
	count := len(ints)
	records := make([][]string, count)
	for j, row := range ints {
		record := make([]string, len(row))
		for i, val := range row {
			record[i] = strconv.Itoa(val)
		}
		records[j] = record
	}
	return records, nil
}

func writeCSVRecordsToFile(path string, records [][]string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	err = writer.WriteAll(records)
	if err != nil {
		return err
	}
	return nil
}


// assemble alliance
// func assemble_alliance(ship_filenames []string, red_name string, blue_name string)

func main() {
	exit_code := 0
	defer os.Exit(exit_code)

	log_lvl := slog.LevelInfo

	// Define arguments
	drrt_directory_arg := flag.String("drrt-directory", ".", "Set the directory the DRRT will be run in.")
	ships_per_alliance_arg := flag.Int("n", 3, "Set the number of ships per alliance.")
	tournament_name_arg := flag.String("tournament-name", "DRRT", "Set the name of the tournament. Red and Blue Alliance fleet files are suffixed with this.")
	log_file_name := flag.String("log-filename", "", "Send log messages to a file. If not set, log to standard error.")

	flag.Parse()

	// Use the slog.TextHandler for the log format
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

	// log the input arguments
	slog.Info("Starting DRRT Scheduler.", "exec", os.Args[0], "version", VERSION)
	slog.Info("Arguments", "drrt-directory", *drrt_directory_arg, "n", *ships_per_alliance_arg, "tournament-name", *tournament_name_arg, "log-filename", *log_file_name)

	// set some path variables to be used later
	ships_directory := filepath.Join(*drrt_directory_arg, "Ships")
	quals_directory := filepath.Join(*drrt_directory_arg, "Qualifications")
	stags_directory := filepath.Join(*drrt_directory_arg, "Staging")
	playf_directory := filepath.Join(*drrt_directory_arg, "Playoffs")
	schej_directory := filepath.Join(*drrt_directory_arg, "schedules")

	slog.Debug("Directories",
		"ships",          ships_directory,
		"qualifications", quals_directory,
		"staging",        stags_directory,
		"playoffs",       playf_directory,
		"schedules",      schej_directory)

	drrt_subdirectories := []string{ships_directory, quals_directory,
		stags_directory, playf_directory, schej_directory}
	for _, dir := range drrt_subdirectories {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			if os.IsExist(err) {
				slog.Warn("Tried to create directory, but it already exists. Continuing.", "dir", dir)
			} else {
				slog.Error("Failed to create directory. Cannot recover.", "dir", dir, "err", err)
				exit_code = 1
				return
			}
		}
	}

	// Empty the contents of the Qualifications directory
	err := lib.Remove_directory_contents(quals_directory)
	if err != nil {
		slog.Error("Cannot remove contents of Qualifications directory.", "err", err)
		exit_code = 1
		return
	}

	// get a slice comprising paths to all ships
	ship_paths, err := get_inspected_ship_paths(ships_directory)
	if err != nil {
		slog.Error("Cannot get inspected ship paths.", "err", err)
		exit_code = 1
		return
	}

	slog.Info("Found paths for ship files.", "count", len(ship_paths))
	for _, path := range ship_paths {
		slog.Debug("Ship path", "path", path)
	}

	// check if there are less ships than can participate in a single match. Fail if this is true.
	if len(ship_paths) < (*ships_per_alliance_arg * 2) {
		slog.Error("Number of participating shps is lower than the minimum number of ships.", "min", *ships_per_alliance_arg * 2, "count", len(ship_paths))
		exit_code = 1
		return
	}

	sch_in_filename := fmt.Sprintf("%d_%dv%d.csv", len(ship_paths), *ships_per_alliance_arg, *ships_per_alliance_arg)
	sch_in_filepath := filepath.Join(schej_directory, "out", sch_in_filename)
	sch_out_filepath := filepath.Join(*drrt_directory_arg, SELECTED_SCHEDULE_FNAME)
	sch_out_filepath_no_asterisks := filepath.Join(*drrt_directory_arg, SELECTED_SCHEDULE_NOAST_FNAME)

	matchschedule, _, err := readScheduleAtPath(sch_in_filepath)
	if err != nil {
		slog.Error("Could not get information from schedule file.", "path", sch_in_filepath, "err", err)
		exit_code = 1
		return
	}
	slog.Info("Schedule information", "path", sch_in_filepath, "matches", matchschedule.Length)
	

	slog.Debug("Starting unmarshalling ships.")

	ships := make([]rsmships.Ship, len(ship_paths))

	var unmarshal_wait_group sync.WaitGroup

	for i, path := range ship_paths {
		unmarshal_wait_group.Add(1)

		go func(i int, path string) {
			defer unmarshal_wait_group.Done()
			fullpath := filepath.Join(ships_directory, path)

			isfleet, err := rsmships.IsReassemblyJSONFileFleet(fullpath)
			if err != nil {
				slog.Error("Failed preparation for unmarshalling ship", "path", fullpath, "err", err)
				exit_code = 1
				return
			}

			if isfleet {
				fleet, err := rsmships.UnmarshalFleetFromFile(fullpath)
				if err != nil {
					slog.Error("Failed unmarshalling fleet", "path", fullpath, "err", err)
					exit_code = 1
					return
				}
				// Use the first blueprint in the fleet file
				ships[i] = *fleet.Blueprints[0]
				slog.Info("Unmarshalled ship from fleet.", "name", ships[i].Data.Name, "author", ships[i].Data.Author, "idx", i + 1, "fleet.Name", fleet.Name)
			} else {
				ships[i], err = rsmships.UnmarshalShipFromFile(fullpath)
				if err != nil {
					slog.Error("Failed unmarshalling ship", "path", fullpath, "err", err)
					exit_code = 1
					return
				}
				slog.Info("Unmarshalled ship", "name", ships[i].Data.Name, "author", ships[i].Data.Author, "idx", i + 1)
			}
		}(i, path)
	}

	schedule := make([]DRRTStandardMatch, matchschedule.Length)
	
	unmarshal_wait_group.Wait()

	slog.Info("Saving alliance fleet files", "dir", ships_directory)
	var save_wait_group sync.WaitGroup
	for i, match := range matchschedule.Schedule {
		save_wait_group.Add(1)

		go func(i int, match []int) {
			defer save_wait_group.Done()

			schedule[i].TournamentName = *tournament_name_arg
			schedule[i].MatchNumber = i

			red  := make([]*rsmships.Ship, *ships_per_alliance_arg)
			blue := make([]*rsmships.Ship, *ships_per_alliance_arg)

			for j, ship := range match {
				if j >= *ships_per_alliance_arg {
					blue[j - *ships_per_alliance_arg] = &ships[ship - 1]
				} else {
					red[j] = &ships[ship - 1]
				}
			}

			schedule[i].RedAlliance  = lib.RED_ALLIANCE_TEMPLATE.CopyUsingShips(red)
			schedule[i].BlueAlliance = lib.BLUE_ALLIANCE_TEMPLATE.CopyUsingShips(blue)

			schedule[i].RedAlliance.Name = fmt.Sprintf("Match %03d - ^1The Red Alliance^7", i+1)
			schedule[i].BlueAlliance.Name = fmt.Sprintf("Match %03d - ^4The Blue Alliance^7", i+1)

			err = WriteMatchFleets(schedule[i], quals_directory)
			if err != nil {
				slog.Error("Failed to save fleets for match", "match", i + 1, "err", err)
				exit_code = 1
				return
			}

			slog.Info("Saved fleets for match", "match", i + 1)
		}(i, match)
	}

	// during marshalling of fleet files, write the match log to a file
	// delete the selected schedule files if they exist
	for _, schedpath := range [2]string{sch_out_filepath, sch_out_filepath_no_asterisks} {
		if err := os.Remove(schedpath); err == nil {
			slog.Info("removed old schedule file.", "path", schedpath)
		} else if errors.Is(err, os.ErrNotExist) {
			slog.Info("No old schedule file exists.", "path", schedpath, "err", err)
		} else {
			slog.Error("Error removing old schedule selection file.", "path", schedpath, "err", err)
			exit_code = 1
			return
		}
	}
	// write the schedule to a file
	err = matchschedule.WriteScheduleToFileSurrogates(sch_out_filepath)
	if err != nil {
		slog.Error("Could not write schedule to file.", "path", sch_out_filepath, "err", err)
		exit_code = 1
		return
	}
	slog.Info("Wrote schedule to file.", "path", sch_out_filepath)
	// write the no-asterisks version of the schedule to a file
	err = matchschedule.WriteScheduleToFileNoSurrogates(sch_out_filepath_no_asterisks)
	if err != nil {
		slog.Error("Could not write no-surrogate schedule to file.", "path", sch_out_filepath_no_asterisks, "err", err)
		exit_code = 1
		return
	}
	slog.Info("Wrote no-surrogates schedule to file.", "path", sch_out_filepath_no_asterisks)

	save_wait_group.Wait()

	slog.Info("Scheduler finished. Have a great tournament!")
}

