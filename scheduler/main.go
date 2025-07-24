package main

import (
	"drrt-scripts/lib"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
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
func get_schedule_from_path(path string) ([][]int, [][]bool, error) {
	schedule_bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	schedule_string := string(schedule_bytes)

	r := csv.NewReader(strings.NewReader(schedule_string))

	records, err := r.ReadAll()
	if err != nil {
		return nil, nil, err
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
				return schedule, surrogates, err
			}
		}
		schedule[j] = ships_in_match
		surrogates[j] = surrogates_in_match
	}
	return schedule, surrogates, nil
}

// assemble alliance
// func assemble_alliance(ship_filenames []string, red_name string, blue_name string)

func main() {
	drrt_directory_arg := flag.String("drrt-directory", ".", "Set the directory the DRRT will be run in.")
	ships_per_alliance_arg := flag.Int("n", 3, "Set the number of ships per alliance.")
	tournament_name_arg := flag.String("tournament-name", "DRRT", "Set the number of ships per alliance.")

	flag.Parse()

	log.Printf("drrt_directory: %s\n", *drrt_directory_arg)
	log.Printf("ships_per_alliance: %d\n", *ships_per_alliance_arg)

	ships_directory := filepath.Join(*drrt_directory_arg, "Ships")
	quals_directory := filepath.Join(*drrt_directory_arg, "Qualifications")
	stags_directory := filepath.Join(*drrt_directory_arg, "Staging")
	playf_directory := filepath.Join(*drrt_directory_arg, "Playoffs")
	schej_directory := filepath.Join(*drrt_directory_arg, "schedules")

	drrt_subdirectories := []string{ships_directory, quals_directory,
		stags_directory, playf_directory, schej_directory}
	for _, dir := range drrt_subdirectories {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			if os.IsExist(err) {
				log.Printf("Directory '%s' already exists.\n", dir)
			} else {
				log.Fatalf("Error making directory '%s': %s\n", dir, err)
			}
		}
	}

	// Empty the contents of the Qualifications directory
	err := lib.Remove_directory_contents(quals_directory)
	if err != nil {
		log.Fatalf("Cannot remove contents of Qualifications directory: %s\n", err)
	}

	ship_paths, err := get_inspected_ship_paths(ships_directory)
	if err != nil {
		log.Fatalf("Cannot get inspected ship paths: %s\n", err)
	}

	log.Printf("Found %d ship files!\n", len(ship_paths))
	for _, path := range ship_paths {
		log.Printf("SHIP: %s\n", path)
	}

	if len(ship_paths) < (*ships_per_alliance_arg * 2) {
		log.Fatalf("Error: %d is lesser than the minimum number of ships (%d).\n",
			len(ship_paths), *ships_per_alliance_arg)
	}

	sch_in_filename := fmt.Sprintf("%d_%dv%d.csv", len(ship_paths), *ships_per_alliance_arg, *ships_per_alliance_arg)
	sch_in_filepath := filepath.Join(schej_directory, "out", sch_in_filename)
	// sch_out_filepath := filepath.Join(scrpt_directory, "selected_schedule.csv")
	// sch_out_filepath_no_asterisks := filepath.Join(scrpt_directory, ".no_asterisks.csv")

	schedule_idxs, _, err := get_schedule_from_path(sch_in_filepath)
	if err != nil {
		log.Fatalf("Could not get scheduling information: %v\n", err)
	}
	log.Printf("Used schedule file: %s\n", sch_in_filepath)
	log.Printf("Schedule has %d matches.\n", len(schedule_idxs))
	log.Printf("Unmarshalling ships.\n")

	ships := make([]lib.Ship, len(ship_paths))

	var unmarshal_wait_group sync.WaitGroup

	for i, path := range ship_paths {
		unmarshal_wait_group.Add(1)

		go func(i int, path string) {
			defer unmarshal_wait_group.Done()
			fullpath := filepath.Join(ships_directory, path)

			isfleet, err := lib.IsReassemblyJSONFileFleet(fullpath)
			if err != nil {
				log.Fatalf("Could not unmarshal ship file '%s': %v", fullpath, err)
			}

			if isfleet {
				fleet, err := lib.UnmarshalFleetFromFile(fullpath)
				if err != nil {
					log.Fatalf("Could not unmarshal fleet file '%s': %v", fullpath, err)
				}
				// Use the first blueprint in the fleet file
				ships[i] = fleet.Blueprints[0]
				log.Printf("Unmarshalled ship '%s' with author '%s' with schedule index %d from fleet '%s'", ships[i].Data.Name, ships[i].Data.Author, i + 1, fleet.Name)
			} else {
				ships[i], err = lib.UnmarshalShipFromFile(fullpath)
				if err != nil { log.Fatalf("Could not unmarshal ship file '%s': %v", fullpath, err) }
				log.Printf("Unmarshalled ship '%s' with author '%s' with schedule index %d", ships[i].Data.Name, ships[i].Data.Author, i + 1)
			}
		}(i, path)
	}

	schedule := make([]DRRTStandardMatch, len(schedule_idxs))
	
	unmarshal_wait_group.Wait()

	log.Printf("Saving alliance fleet files to '%s'.", ships_directory)
	var save_wait_group sync.WaitGroup
	for i, match := range schedule_idxs {
		save_wait_group.Add(1)

		go func(i int, match []int) {
			defer save_wait_group.Done()

			schedule[i].TournamentName = *tournament_name_arg
			schedule[i].MatchNumber = i

			red  := make([]lib.Ship, *ships_per_alliance_arg)
			blue := make([]lib.Ship, *ships_per_alliance_arg)

			for j, ship := range match {
				if j >= *ships_per_alliance_arg {
					blue[j - *ships_per_alliance_arg] = ships[ship - 1]
				} else {
					red[j] = ships[ship - 1]
				}
			}

			schedule[i].RedAlliance  = lib.AssembleAlliance(lib.RED_ALLIANCE_TEMPLATE, red)
			schedule[i].BlueAlliance = lib.AssembleAlliance(lib.BLUE_ALLIANCE_TEMPLATE, blue)

			schedule[i].RedAlliance.Name = fmt.Sprintf("Match %03d - ^1The Red Alliance^7", i+1)
			schedule[i].BlueAlliance.Name = fmt.Sprintf("Match %03d - ^4The Blue Alliance^7", i+1)

			err = WriteMatchFleets(schedule[i], quals_directory)
			if err != nil {
				log.Fatalf("Failed to save fleets for match %d: %v", i + 1, err)
			}

			log.Printf("Saved fleets for match %d", i + 1)
		}(i, match)
	}

	save_wait_group.Wait()

	log.Printf("Done. Have a great tournament!")
}

