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
func (m *matchPerformance) scoreKill() {
	m.Destructions += 1
	m.RankPointsEarned += 1
}
func (m *matchPerformance) scoreSurvived() {
	m.Survived = true
}
func (m *matchPerformance) scoreDestructionWin() {
	m.Result = WinDestruction
	m.RankPointsEarned += 2
}
func (m *matchPerformance) scorePointsWin() {
	m.Result = WinPoints
	m.RankPointsEarned += 2
}
func (m *matchPerformance) scoreLoss() {
	m.Result = Loss
}

func (m *matchPerformance) toSheetsRow() [][]any {
	output := make([]any, 9)
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
	output[8] = int(m.Result)
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

// Create a map that connects a ship's name to its index in the schedule.
// nameCorrelator secont argument also stores the ship's faction. This is in
// case it is needed to parse a match log, and may be changed many times.
func getShipIdxFacMap(ships []*rsmships.Ship) *map[string]int {
	nametoidx := make(map[string]int)
	for i, ship := range ships {
		// NOTE: the ships' names must not use standard name format (name [by author])
		// NOTE: value is 1 less than match schedule values; schedule starts at 1 and not 0. internally, use 0 minimum. Print 1+ this value.
		nameauthor := ShipAuthorFromCommonNamefmt(ship.Data.Name)
		nametoidx[nameauthor[0]] = &nameCorrelator{i, 0}
	}
	return &nametoidx
}

// TODO: make the nametoidx variable an input argument
// TODO: This function is long and unweildly
func NewDRRTStandardMatchLogFromShips(raw *MatchLogRaw, ships []*rsmships.Ship, nametoidx *map[string]int) (*DRRTStandardMatchLog, error) {
	var mlog DRRTStandardMatchLog
	mlog.Raw = raw
	mlog.Timestamp = raw.CreatedTimestamp

	// ASSUMPTION: NOT a free-for-all (Red v Blue alliance)
	// get the match number
	// if the same for red and blue alliances: good!
	redMatchNumber := GetMatchNumberFromAllianceName(raw.StartListings[0].Name, false)
	blueMatchNumber := GetMatchNumberFromAllianceName(raw.StartListings[1].Name, true)
	if redMatchNumber != blueMatchNumber {
		return &mlog, &MatchLogAllianceMatchNumberMismatchError{
			redAllianceMatchNumber: redMatchNumber,
			blueAllianceMatchNumber: blueMatchNumber,
		}
	}
	mlog.MatchNumber = redMatchNumber // == blueMatchNumber
	slog.Debug("found match log number from filenames", "matchNumber", mlog.MatchNumber, "mlogTimestamp", mlog.Timestamp.String(), "path", raw.Path)

	// Check the length of both alliances
	redAllianceLength  := raw.StartListings[0].Alive
	blueAllianceLength := raw.StartListings[1].Alive
	if redAllianceLength != blueAllianceLength {
		return &mlog, &MatchLogAllianceLengthMismatchError{
			redAllianceLength: redAllianceLength,
			blueAllianceLength: blueAllianceLength,
		}
	}
	mlog.AllianceLength = redAllianceLength


	// FIXME: currently, there is no sanity checking of match log input. A fleet
	// line could say that an alliance has 4 members when only 3 are present,
	// and such.

	// Get the indices of each ship participating in this match log
	// The first n ships are from the Red alliance, but may not be sorted in the order they appear in Reassembly's fleet screen.
	mlog.ShipIndices = make([]int, len(raw.ShipListings))
	mlog.Ships = make([]*rsmships.Ship, len(raw.ShipListings))
	factions := make([]int, len(raw.ShipListings))
	for i, shiplsting := range raw.ShipListings {
		name := ShipAuthorFromCommonNamefmt(shiplsting.Ship)[0]
		var idx int
		idx, ok := (*nametoidx)[name]
		if !ok {
			slog.Warn("Ship index cannot be found using map.", "scoring", "SHIP", "name", name, "matchNumber", mlog.MatchNumber, "mlogTimestamp", mlog.Timestamp.String(), "path", raw.Path)
			continue
			// TODO: you may want to return an error here.. but I don't know.
		}
		mlog.ShipIndices[i] = idx
		mlog.Ships[i] = ships[idx]
		factions[i] = shiplsting.Fleet
	}

	// Map the ship's index value to its performance in the match. This contains
	// the same references as mlog.Records.
	idxtoperformance := make(map[int]*matchPerformance)
	mlog.Record = make([]*matchPerformance, mlog.AllianceLength * 2)
	// create an empty matchPerformance entry for each ship
	for i, idx := range mlog.ShipIndices {
		// get the faction of the ship
		p := &matchPerformance{Ship: ships[idx], Match: mlog.MatchNumber, Faction: factions[i]}
		mlog.Record[i] = p
		idxtoperformance[idx] = p
		slog.Debug("Add ship to match performance.", "author", ships[idx].Data.Author, "name", ships[idx].Data.Name, "idx", idx, "path", raw.Path)
	}

	// add the [DESTRUCTION] mlog information to the datatype
	for _, destruction := range raw.DestructionListings {
		destroyername := ShipAuthorFromCommonNamefmt(destruction.Ship)[0]
		idx, ok := (*nametoidx)[destroyername]
		if !ok {
			slog.Warn("Ship index of destroying ship cannot be found using map.", "scoring", "DESTRUCTION", "name", destroyername, "matchNumber", mlog.MatchNumber, "mlogTimestamp", mlog.Timestamp.String(), "path", raw.Path)
			// TODO: you may want to return an error here.. but I don't know.
			continue
		}
		// assign values to the matchPerformance of the ship whose idx was found
		var p *matchPerformance
		p, ok = idxtoperformance[idx]
		if !ok {
			slog.Warn("Cannot find performance of ship from index.", "scoring", "DESTRUCTION", "name", destroyername, "matchNumber", mlog.MatchNumber, "mlogTimestamp", mlog.Timestamp.String(), "path", raw.Path)
			continue
		}
		// SCORING
		// if a ship destroys another ship, increase ranking points earned and destructions by 1
		p.scoreKill()
	}
	
	// assertion:
	if len(raw.SurvivalListings) <= 0 {
		return &mlog, errors.New("Match log not finished: no survival lines.")
	}
	// add surviving ship information
	for _, survival := range raw.SurvivalListings {
		name := ShipAuthorFromCommonNamefmt(survival.Ship)[0]
		idx, ok := (*nametoidx)[name]
		if !ok {
			slog.Warn("Ship index cannot be found using map.", "scoring", "SURVIVAL", "name", name, "matchNumber", mlog.MatchNumber, "mlogTimestamp", mlog.Timestamp.String(), "path", raw.Path)
			// TODO: you may want to return an error here.. but I don't know.
		}
		var p *matchPerformance
		p, ok = idxtoperformance[idx]
		if !ok {
			slog.Warn("Cannot find performance of ship from index.", "scoring", "SURVIVAL", "name", name, "matchNumber", mlog.MatchNumber, "mlogTimestamp", mlog.Timestamp.String(), "path", raw.Path)
			continue
		}
		// SCORING
		p.scoreSurvived()
	}
	
	// -- SCORING --
	// Define which type of points the match shall be scored with. 
	pointsMethod := func(lst *MatchLogFleetListing) int {
		return lst.DamageInflicted
	}
	// initialize functions that will be used to send points to each alliance
	var blueResultScorer func(p *matchPerformance)
	var redResultScorer  func(p *matchPerformance)


	// assertions:
	if len(raw.ResultListings) <= 0 {
		// FIXME: make this a real error type
		return &mlog, errors.New("Match log not finished: result lines not present.")
	}
	if len(raw.ResultListings) > 2 {
		slog.Warn("More result lines than expected for standard DRRT match.", "scoring", "RESULT", "len", len(raw.ResultListings), "matchNumber", mlog.MatchNumber, "mlogTimestamp", mlog.Timestamp.String(), "path", raw.Path)
	}

	// add RESULT line infomration
	if raw.ResultListings[0].Alive <= 0 {
		// red alliance loses on destruction
		blueResultScorer = func(p *matchPerformance) { p.scoreDestructionWin() }
		redResultScorer = func(p *matchPerformance) { p.scoreLoss() }
	} else if raw.ResultListings[1].Alive <= 0 {
		// blue alliance loses on destruction
		blueResultScorer = func(p *matchPerformance) { p.scoreLoss() }
		redResultScorer = func(p *matchPerformance) { p.scoreDestructionWin() }
	} else if pointsMethod(&raw.ResultListings[0]) >= pointsMethod(&raw.ResultListings[1]) {
		// red alliance wins on points.
		// If tie, default to red alliance, just like Reassembly.
		blueResultScorer = func(p *matchPerformance) { p.scoreLoss() }
		redResultScorer = func(p *matchPerformance) { p.scorePointsWin() }
	} else {
		// blue alliance wins on points
		blueResultScorer = func(p *matchPerformance) { p.scorePointsWin() }
		redResultScorer = func(p *matchPerformance) { p.scoreLoss() }
	}
	// -- SCORING --
	for i, rec := range mlog.Record {
		if i < mlog.AllianceLength {
			// red alliance
			redResultScorer(rec)
		} else {
			// blue alliance
			blueResultScorer(rec)
		}
	}

	// fill the remaining fields
	mlog.PointsDamageInflicted = make([]int, len(raw.ResultListings))
	mlog.PointsDamageTaken     = make([]int, len(raw.ResultListings))
	for i, result := range raw.ResultListings {
		mlog.PointsDamageInflicted[i] = result.DamageInflicted
		mlog.PointsDamageTaken[i]     = result.DamageTaken
	}

	return &mlog, nil
}

