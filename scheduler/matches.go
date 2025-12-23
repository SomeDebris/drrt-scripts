package main

import (
	"github.com/SomeDebris/rsmships-go"
	"path/filepath"
	"fmt"
)

type DRRTStandardMatch struct {
	MatchNumber    int
	TournamentName string
	RedAlliance    rsmships.Fleet
	BlueAlliance   rsmships.Fleet
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
