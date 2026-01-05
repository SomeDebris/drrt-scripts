package lib

import (
	"github.com/SomeDebris/rsmships-go"
	"path/filepath"
	"fmt"
	"strconv"
	"strings"
	"os"
	"encoding/csv"
)

type ScheduleLengthMismatch struct {
	SchedulesLength  int // length of schedules slice
	SurrogatesLength int // length of surrogates slice
}
type MatchLengthMismatch struct {
	SchedulesLength  int // length of schedules slice
	SurrogatesLength int // length of surrogates slice
}
func (m *ScheduleLengthMismatch) Error() string {
	return "Length of schedule indices slice (" + strconv.Itoa(m.SchedulesLength) + ") and surrogates slice (" + strconv.Itoa(m.SurrogatesLength) + ") should be equivalent, but are not."
}
func (m *MatchLengthMismatch) Error() string {
	return "Length of match's schedule indices slice (" + strconv.Itoa(m.SchedulesLength) + ") and match's surrogates slice (" + strconv.Itoa(m.SurrogatesLength) + ") should be equivalent, but are not."
}

type MatchSchedule struct {
	Schedule       [][]int
	Surrogates     [][]bool
	Length         int
	AllianceLength int
}

type DRRTStandardMatch struct {
	MatchNumber    int
	TournamentName string
	RedAlliance    rsmships.Fleet
	BlueAlliance   rsmships.Fleet
}

// Create 2d string slice of schedule (without surrogate-indicating asterisks).
func (m *MatchSchedule) GetRecordsNoSurrogates() ([][]string, error) {
	return Int2dSliceToString(m.Schedule)
}

// Create 2d string slice of schedule. An asterisk is printed after the ship
// index if the ship is to play the match as a surrogate.
func (m *MatchSchedule) GetRecordsSurrogates() ([][]string, error) {
	records := make([][]string, m.Length)
	var builder strings.Builder
	for j, match := range m.Schedule {
		record := make([]string, len(match))
		for i, ship := range match {
			_, err := builder.WriteString(strconv.Itoa(ship))
			if err != nil {
				return records, err
			}
			if m.Surrogates[j][i] {
				_, err := builder.WriteString("*")
				if err != nil {
					return records, err
				}
			}
			record[i] = builder.String()
			builder.Reset()
		}
		records[j] = record
	}
	return records, nil
}

// Write schedule (no surrogate-indicating asterisks) to csv file.
func (m *MatchSchedule) WriteScheduleToFileNoSurrogates(path string) error {
	records, err := m.GetRecordsNoSurrogates()
	if err != nil {
		return nil
	}
	return WriteCSVRecordsToFile(path, records)
}

// Write schedule (with surrogate-indicating asterisks) to csv file.
func (m *MatchSchedule) WriteScheduleToFileSurrogates(path string) error {
	records, err := m.GetRecordsSurrogates()
	if err != nil {
		return nil
	}
	return WriteCSVRecordsToFile(path, records)
}


func WriteMatchFleets(match DRRTStandardMatch, directory string) error {
	redpath := filepath.Join(directory, fmt.Sprintf("%s_%s.json.gz", match.RedAlliance.Name, match.TournamentName))
	bluepath := filepath.Join(directory, fmt.Sprintf("%s_%s.json.gz", match.BlueAlliance.Name, match.TournamentName))

	err := rsmships.MarshalFleetToFileGzip(redpath, match.RedAlliance)
	if err != nil {
		return err
	}

	err = rsmships.MarshalFleetToFileGzip(bluepath, match.BlueAlliance)
	if err != nil {
		return err
	}

	return nil
}

/**
* schedule array
* whether ship is participating as surrogate
* any error recieved
 */
func ReadScheduleAtPath(path string) (MatchSchedule, [][]string, error) {
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
