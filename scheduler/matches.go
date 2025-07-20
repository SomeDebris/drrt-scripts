package main

import (
	"drrt-scripts/lib"
	"path/filepath"
	"fmt"
)

type DRRTStandardMatch struct {
	MatchNumber    int
	TournamentName string
	RedAlliance    lib.Fleet
	BlueAlliance   lib.Fleet
}

func WriteMatchFleets(match DRRTStandardMatch, directory string) error {
	redpath := filepath.Join(directory, fmt.Sprintf("%s_%s.json", match.RedAlliance.Name, match.TournamentName))
	bluepath := filepath.Join(directory, fmt.Sprintf("%s_%s.json", match.BlueAlliance.Name, match.TournamentName))

	err := lib.MarshalFleetToFile(redpath, match.RedAlliance)
	if err != nil {
		return err
	}

	err = lib.MarshalFleetToFile(bluepath, match.BlueAlliance)
	if err != nil {
		return err
	}

	return nil
}
