package lib

import (
	"errors"
	"log/slog"
	"time"

	"github.com/SomeDebris/rsmships-go"
)

type DRRTShipStats struct {
	Name                string  `json:"name"`
	Author              string  `json:"author"`
	RankPoints          float64 `json:"rps"`
	NumberMatchesPlayed int     `json:"matchesPlayed"`
	QualsSeed           int
	Faction             int
	P                   int
	ShipData            *rsmships.Ship
}


func (m *DRRTShipStats) RankingScore() float64 {
	return m.RankPoints / float64(m.NumberMatchesPlayed)
}

type matchResult int
const (
	WinDestruction matchResult = iota
	WinPoints
	Loss
)


type matchPerformance struct {
	Match            int
	Ship             *rsmships.Ship
	Faction          int
	Destructions     int
	RankPointsEarned int
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

// Tiny structure for storing the index of a ship and the ship's faction. This
// is so that we need less lookups in the hashmap
type nameCorrelator struct {
	Idx     int
	Faction int
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
	Raw                   *MatchLogRaw
}

// TODO: make the nametoidx variable an input argument
// TODO: This function is long and unweildly
func NewDRRTStandardMatchLogFromShips(raw *MatchLogRaw, ships []*rsmships.Ship) (*DRRTStandardMatchLog, error) {
	var mlog DRRTStandardMatchLog
	mlog.Raw = raw
	mlog.Timestamp = raw.CreatedTimestamp

	// ASSUMPTION: NOT a free-for-all (Red v Blue alliance)
	// get the match number
	// if the same for red and blue alliances: good!
	redMatchNumber := GetMatchNumberFromAllianceName(raw.StartListings[0].Name, false)
	blueMatchNumber := GetMatchNumberFromAllianceName(raw.StartListings[1].Name, true)
	if redMatchNumber != blueMatchNumber {
		return &mlog, errors.New("Red and Blue Alliance match numbers are different. Bad match log!")
	}
	mlog.MatchNumber = redMatchNumber // == blueMatchNumber
	slog.Debug("found match log number from filenames", "matchNumber", mlog.MatchNumber, "mlogTimestamp", mlog.Timestamp.String(), "path", raw.Path)

	// collect array of shipidxs
	// TODO: possibly good idea to accept ship list from the spreadsheet and stop using actual ship datatypes?
	nametoidx := make(map[string]*nameCorrelator)
	for i, ship := range ships {
		// NOTE: the ships' names must not use standard name format (name [by author])
		// NOTE: value is 1 less than match schedule values; schedule starts at 1 and not 0. internally, use 0 minimum. Print 1+ this value.
		nametoidx[ship.Data.Name] = &nameCorrelator{i, 0}
		slog.Debug("Ship in match", "matchNumber", mlog.MatchNumber, "name", ship.Data.Name, "mlogTimestamp", mlog.Timestamp.String(), "path", raw.Path)
	}

	mlog.AllianceLength = len(raw.ShipListings) / 2
	// Get the indices of each ship participating in this match log
	// The first n ships are from the Red alliance, but may not be sorted in the order they appear in Reassembly's fleet screen.
	mlog.ShipIndices = make([]int, len(raw.ShipListings))
	mlog.Ships = make([]*rsmships.Ship, len(raw.ShipListings))
	for i, shiplsting := range raw.ShipListings {
		name := ShipAuthorFromCommonNamefmt(shiplsting.Ship)[0]
		nameCorrelator, ok := nametoidx[name]
		if !ok {
			slog.Warn("Ship index cannot be found using map.", "name", name, "matchNumber", mlog.MatchNumber, "mlogTimestamp", mlog.Timestamp.String(), "path", raw.Path)
			// TODO: you may want to return an error here.. but I don't know.
		}
		mlog.ShipIndices[i] = nameCorrelator.Idx
		mlog.Ships[i] = ships[nameCorrelator.Idx]
		nameCorrelator.Faction = shiplsting.Fleet
	}
	


	// Map the ship's index value to its performance in the match
	idxtoperformance := make(map[int]*matchPerformance)
	// create an empty matchPerformance entry for each ship
	for _, idx := range mlog.ShipIndices {
		// get the faction of the ship
		idxtoperformance[idx] = &matchPerformance{Ship: ships[idx], Match: mlog.MatchNumber}
		slog.Debug("Add ship to match performance.", "author", ships[idx].Data.Author, "name", ships[idx].Data.Name, "idx", idx, "path", raw.Path)
	}
	// add the [DESTRUCTION] mlog information to the datatype
	for _, destruction := range raw.DestructionListings {
		destroyername := ShipAuthorFromCommonNamefmt(destruction.Ship)[0]
		destroyernameCorrelator, ok := nametoidx[destroyername]
		if !ok {
			slog.Warn("Ship index of destroying ship cannot be found using map.", "name", destroyername, "matchNumber", mlog.MatchNumber, "mlogTimestamp", mlog.Timestamp.String(), "path", raw.Path)
			// TODO: you may want to return an error here.. but I don't know.
			continue
		}
		// assign values to the matchPerformance of the ship whose idx was found
		var p *matchPerformance
		p, ok = idxtoperformance[destroyernameCorrelator.Idx]
		if !ok {
			slog.Warn("Cannot find performance of ship from index.", "name", destroyername, "matchNumber", mlog.MatchNumber, "mlogTimestamp", mlog.Timestamp.String(), "path", raw.Path)
			continue
		}
		// SCORING
		// if a ship destroys another ship, increase ranking points earned and destructions by 1
		p.Destructions += 1
		p.RankPointsEarned += 1
	}

	// add surviving ship information
	for _, survival := range raw.SurvivalListings {
		// SCORING
		// If score is to be added due to a survival, do it here. Do a switch
		// statement on the faction name et cetera.
		// switch survival.Fleet {
		// case MLOG_RED_FACTION:

	}


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
