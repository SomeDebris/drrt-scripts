package lib

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
	"errors"
	"path/filepath"
	"strings"
	"sync"
	"github.com/SomeDebris/rsmships-go"
	"log/slog"
)

/** example
MLOG OPEN: 793.504312
[START] faction:{100} name:{Interceptor} DT:{0} DI:{0} alive:{2}
[SHIP] fleet:{100} ship:{Street Sweeper v0.1.5 [by Debris]}
[SHIP] fleet:{100} ship:{Octal v1.1 [by Debris]}
[START] faction:{101} name:{WhatTheHell} DT:{0} DI:{0} alive:{3}
[SHIP] fleet:{101} ship:{DesmondTheDiamond}
[SHIP] fleet:{101} ship:{001-0013}
[SHIP] fleet:{101} ship:{WhatTheHell}
[DESTRUCTION] ship:{Street Sweeper v0.1.5 [by Debris]} fship:{100} destroyed:{WhatTheHell} fdestroyed:{101}
[DESTRUCTION] ship:{Octal v1.1 [by Debris]} fship:{100} destroyed:{DesmondTheDiamond} fdestroyed:{101}
[RESULT] faction:{100} name:{Interceptor} DT:{12093} DI:{6668} alive:{2}
[SURVIVAL] fleet:{100} ship:{Street Sweeper v0.1.5 [by Debris]}
[SURVIVAL] fleet:{100} ship:{Octal v1.1 [by Debris]}
[RESULT] faction:{101} name:{WhatTheHell} DT:{6337} DI:{30650} alive:{1}
[SURVIVAL] fleet:{101} ship:{001-0013}
MLOG CLOSE: 855.581397
*/

const (
	mlog_typeRegexCaptureString        = `^\[([A-Z]+)\]`
	mlog_startRegexCaptureString       = `^\[START\] faction:\{([0-9]+)\} name:\{(.*)\} DT:\{([0-9]*)\} DI:\{([0-9])*\} alive:\{([0-9]*)\}$`
	mlog_shipRegexCaptureString        = `^\[SHIP\] faction:\{([0-9]+)\} ship:\{(.*)\}$`
	mlog_destructionRegexCaptureString = `^\[DESTRUCTION\] ship:\{(.*)\} fship:\{([0-9]*)\} destroyed:\{(.*)\} fdestroyed:\{([0-9]*)\}$`
	mlog_resultRegexCaptureString      = `^\[RESULT\] faction:\{([0-9]+)\} name:\{(.*)\} DT:\{([0-9]*)\} DI:\{([0-9])*\} alive:\{([0-9]*)\}$`
	mlog_survivalRegexCaptureString    = `^\[SURVIVAL\] faction:\{([0-9]+)\} ship:\{(.*)\}$`
	mlog_shipauthorRegexCaptureString  = `^(.+) \[by (.+)\]$`

	mlog_qualsRedCaptureString  = `^Match ([0-9]+) - \^1The Red Alliance\^7$`
	mlog_qualsBlueCaptureString = `^Match ([0-9]+) - \41The Blue Alliance\^7$`
	MLOG_RED_FACTION            = 100
	MLOG_BLUE_FACTION           = 101

	mlog_start       = `START`
	mlog_ship        = `SHIP`
	mlog_destruction = `DESTRUCTION`
	mlog_result      = `RESULT`
	mlog_survival    = `SURVIVAL`

	length_fleetListingSlice  = 5
	length_resultListingSlice = 5

	length_shipListingSlice     = 2
	length_survivalListingSlice = 2

	length_destructionListingSlice = 4
)

var (
	mlog_regex_matchnumberCaptureRed  = regexp.MustCompile(mlog_qualsRedCaptureString)
	mlog_regex_matchnumberCaptureBlue = regexp.MustCompile(mlog_qualsBlueCaptureString)
	mlog_regex_type                   = regexp.MustCompile(mlog_typeRegexCaptureString)
	mlog_regex_shipauthor             = regexp.MustCompile(mlog_shipRegexCaptureString)
	mlog_regex_map                    = map[string]*regexp.Regexp{
		mlog_start:       regexp.MustCompile(mlog_startRegexCaptureString),
		mlog_ship:        regexp.MustCompile(mlog_shipRegexCaptureString),
		mlog_destruction: regexp.MustCompile(mlog_destructionRegexCaptureString),
		mlog_result:      regexp.MustCompile(mlog_resultRegexCaptureString),
		mlog_survival:    regexp.MustCompile(mlog_survivalRegexCaptureString),
	}
)

func ShipAuthorFromCommonNamefmt(name string) [2]string {
	fields := mlog_regex_shipauthor.FindStringSubmatch(name)
	if fields == nil {
		return [2]string{name, ""}
	}
	return [2]string{fields[1], fields[2]}
}

// Get the match number from the filename of the match log. This works for Red
// Alliance fleets. Returns 0 when an error is encountered or no match can be
// found.
func GetMatchNumberFromAllianceName(allianceName string, isBlue bool) int {
	var fields []string
	if isBlue {
		fields = mlog_regex_matchnumberCaptureBlue.FindStringSubmatch(allianceName)
	} else {
		fields = mlog_regex_matchnumberCaptureRed.FindStringSubmatch(allianceName)
	}
	if fields == nil {
		return 0
	}
	out, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0
	}
	return out
}

type MatchLogFleetListing struct {
	Faction         int
	Name            string
	DamageTaken     int
	DamageInflicted int
	Alive           int
}

type MatchLogShipListing struct {
	Fleet int
	Ship  string
}

type MatchLogDestructionListing struct {
	Ship       string
	Fship      int
	Destroyed  string
	Fdestroyed int
}

// Raw match log type. The data has been parsed out of the log file, but nothing has been done to the data yet.
type MatchLogRaw struct {
	CreatedTimestamp    time.Time
	Path                string
	StartListings       []MatchLogFleetListing
	ShipListings        []MatchLogShipListing
	DestructionListings []MatchLogDestructionListing
	ResultListings      []MatchLogFleetListing
	SurvivalListings    []MatchLogShipListing
}

var (
	matchLogRawMutex_start       sync.Mutex
	matchLogRawMutex_ship        sync.Mutex
	matchLogRawMutex_destruction sync.Mutex
	matchLogRawMutex_result      sync.Mutex
	matchLogRawMutex_survival    sync.Mutex
)

// Append a START event.
func (mlograw *MatchLogRaw) appendStart(record MatchLogFleetListing, mutex *sync.Mutex) {
	mutex.Lock()
	mlograw.StartListings = append(mlograw.StartListings, record)
	mutex.Unlock()
}

// Append a SHIP event.
func (mlograw *MatchLogRaw) appendShip(record MatchLogShipListing, mutex *sync.Mutex) {
	mutex.Lock()
	mlograw.ShipListings = append(mlograw.ShipListings, record)
	mutex.Unlock()
}

// Append a DESTRUCTION event.
func (mlograw *MatchLogRaw) appendDestruction(record MatchLogDestructionListing, mutex *sync.Mutex) {
	mutex.Lock()
	mlograw.DestructionListings = append(mlograw.DestructionListings, record)
	mutex.Unlock()
}

// Append a RESULT event.
func (mlograw *MatchLogRaw) appendResult(record MatchLogFleetListing, mutex *sync.Mutex) {
	mutex.Lock()
	mlograw.ResultListings = append(mlograw.ResultListings, record)
	mutex.Unlock()
}

// Append a SURVIVAL event.
func (mlograw *MatchLogRaw) appendSurvival(record MatchLogShipListing, mutex *sync.Mutex) {
	mutex.Lock()
	mlograw.SurvivalListings = append(mlograw.SurvivalListings, record)
	mutex.Unlock()
}

func parseDestructionLine(line string) (MatchLogDestructionListing, error) {
	var listing MatchLogDestructionListing
	var err error
	regex := mlog_regex_map[mlog_destruction]

	fields := regex.FindStringSubmatch(line)

	if fields == nil {
		return listing, &MatchLogRegexError{
			event: mlog_ship,
			line:  line,
			regex: regex.String(),
		}
	}

	listing.Ship = fields[1]
	listing.Fship, err = strconv.Atoi(fields[2])
	if err != nil {
		return listing, &MatchLogFieldError{
			message: err.Error(),
			event:   mlog_destruction,
			field:   "fship",
			line:    line,
		}
	}
	listing.Destroyed = fields[3]
	listing.Fdestroyed, err = strconv.Atoi(fields[4])
	if err != nil {
		return listing, &MatchLogFieldError{
			message: err.Error(),
			event:   mlog_destruction,
			field:   "fdestroyed",
			line:    line,
		}
	}

	return listing, nil
}

func GetTimeOfMatchLogFilename(path string) (time.Time, error) {
	basename := filepath.Base(path)
	trimmed, found := strings.CutPrefix(basename, MLOG_PREFIX)
	if !found {
		return time.Time{}, fmt.Errorf("Cannot find prefix `%s`.", MLOG_PREFIX)
	}
	timestamp, found := strings.CutSuffix(trimmed, MLOG_EXTENSION)
	if !found {
		return time.Time{}, fmt.Errorf("Cannot find extension `%s`.", MLOG_EXTENSION)
	}
	return time.Parse(REASSEMBLY_FILE_TIMESTAMP_FMT, timestamp)
}

// *
func NewMatchLogRawFromPath(path string) (*MatchLogRaw, error) {

	match_log, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer match_log.Close()

	var mlog_raw MatchLogRaw
	mlog_raw.Path = path

	/*
		// get the last modified time of the file:
		mlog_fileinfo, err := match_log.Stat()
		if err != nil {
			return mlog_raw, err
		}
		// Set the CreatedTimestamp value to the last modified time of the file
		mlog_raw.CreatedTimestamp = mlog_fileinfo.ModTime()
		// */
	mlog_raw.CreatedTimestamp, err = GetTimeOfMatchLogFilename(path)
	if err != nil {
		return &mlog_raw, err
	}

	match_log_scanner := bufio.NewScanner(match_log)

	mlog_RecordNumber := 0

	for match_log_scanner.Scan() {
		line := match_log_scanner.Text()

		mlog_RecordNumber++

		switch string(mlog_regex_type.FindString(line)) {
		// If line is a [START] Line
		case mlog_start:
			fields := mlog_regex_map[mlog_start].FindStringSubmatch(line)

			if fields == nil {
				return &mlog_raw, &MatchLogRegexError{
					event:      mlog_start,
					line:       line,
					lineNumber: mlog_RecordNumber,
					path:       path,
					regex:      mlog_regex_map[mlog_start].String(),
				}
			}
			var listing MatchLogFleetListing
			listing.Faction, err = strconv.Atoi(fields[1])
			if err != nil {
				return &mlog_raw, &MatchLogFieldError{
					field:      "faction",
					event:      mlog_start,
					lineNumber: mlog_RecordNumber,
					line:       line,
					path:       path,
				}
			}
			listing.Name = fields[2]
			listing.DamageTaken, err = strconv.Atoi(fields[3])
			if err != nil {
				return &mlog_raw, &MatchLogFieldError{
					message:    err.Error(),
					field:      "DT",
					event:      mlog_start,
					lineNumber: mlog_RecordNumber,
					line:       line,
					path:       path,
				}
			}
			listing.DamageInflicted, err = strconv.Atoi(fields[4])
			if err != nil {
				return &mlog_raw, &MatchLogFieldError{
					message:    err.Error(),
					field:      "DI",
					event:      mlog_start,
					lineNumber: mlog_RecordNumber,
					line:       line,
					path:       path,
				}
			}
			listing.Alive, err = strconv.Atoi(fields[5])
			if err != nil {
				return &mlog_raw, &MatchLogFieldError{
					message:    err.Error(),
					field:      "alive",
					event:      mlog_start,
					lineNumber: mlog_RecordNumber,
					line:       line,
					path:       path,
				}
			}
			mlog_raw.appendStart(listing, &matchLogRawMutex_start)
		// If line is a [SHIP] line:
		case mlog_ship:
			fields := mlog_regex_map[mlog_ship].FindStringSubmatch(line)
			if fields == nil {
				return &mlog_raw, &MatchLogRegexError{
					event:      mlog_ship,
					line:       line,
					lineNumber: mlog_RecordNumber,
					path:       path,
					regex:      mlog_regex_map[mlog_ship].String(),
				}
			}
			var listing MatchLogShipListing
			listing.Fleet, err = strconv.Atoi(fields[1])
			if err != nil {
				return &mlog_raw, &MatchLogFieldError{
					message:    err.Error(),
					event:      mlog_ship,
					field:      "faction",
					line:       line,
					lineNumber: mlog_RecordNumber,
					path:       path,
				}
			}
			listing.Ship = fields[2]
			mlog_raw.appendShip(listing, &matchLogRawMutex_ship)

		case mlog_destruction:
			listing, err := parseDestructionLine(line)
			if err != nil {
				regexerr := &MatchLogRegexError{}
				fielderr := &MatchLogFieldError{}
				switch {
				case errors.As(err, &regexerr):
					regexerr.AddContext(mlog_RecordNumber, path)
					err = regexerr
				case errors.As(err, &fielderr):
					fielderr.AddContext(mlog_RecordNumber, path)
					err = fielderr
				}
				return &mlog_raw, err
			}
			mlog_raw.appendDestruction(listing, &matchLogRawMutex_destruction)

		case mlog_result:
			fields := mlog_regex_map[mlog_result].FindStringSubmatch(line)

			if fields == nil {
				return &mlog_raw, &MatchLogRegexError{
					event:      mlog_result,
					line:       line,
					lineNumber: mlog_RecordNumber,
					path:       path,
					regex:      mlog_regex_map[mlog_result].String(),
				}
			}
			var listing MatchLogFleetListing
			listing.Faction, err = strconv.Atoi(fields[1])
			if err != nil {
				return &mlog_raw, &MatchLogFieldError{
					field:      "faction",
					event:      mlog_result,
					lineNumber: mlog_RecordNumber,
					line:       line,
					path:       path,
				}
			}
			listing.Name = fields[2]
			listing.DamageTaken, err = strconv.Atoi(fields[3])
			if err != nil {
				return &mlog_raw, &MatchLogFieldError{
					message:    err.Error(),
					field:      "DT",
					event:      mlog_result,
					lineNumber: mlog_RecordNumber,
					line:       line,
					path:       path,
				}
			}
			listing.DamageInflicted, err = strconv.Atoi(fields[4])
			if err != nil {
				return &mlog_raw, &MatchLogFieldError{
					message:    err.Error(),
					field:      "DI",
					event:      mlog_result,
					lineNumber: mlog_RecordNumber,
					line:       line,
					path:       path,
				}
			}
			listing.Alive, err = strconv.Atoi(fields[5])
			if err != nil {
				return &mlog_raw, &MatchLogFieldError{
					message:    err.Error(),
					field:      "alive",
					event:      mlog_result,
					lineNumber: mlog_RecordNumber,
					line:       line,
					path:       path,
				}
			}
			mlog_raw.appendResult(listing, &matchLogRawMutex_result)

		case mlog_survival:
			fields := mlog_regex_map[mlog_survival].FindStringSubmatch(line)
			if fields == nil {
				return &mlog_raw, &MatchLogRegexError{
					event:      mlog_survival,
					line:       line,
					lineNumber: mlog_RecordNumber,
					path:       path,
					regex:      mlog_regex_map[mlog_survival].String(),
				}
			}
			var listing MatchLogShipListing
			listing.Fleet, err = strconv.Atoi(fields[1])
			if err != nil {
				return &mlog_raw, &MatchLogFieldError{
					message:    err.Error(),
					event:      mlog_survival,
					field:      "faction",
					line:       line,
					lineNumber: mlog_RecordNumber,
					path:       path,
				}
			}
			listing.Ship = fields[2]
			mlog_raw.appendShip(listing, &matchLogRawMutex_survival)
		}
	}

	// We have successfully assembled a mlog_raw. from this, we shall make a DRRTStandardTerseMatchLog.
	// TODO how does one get the timestamp at which a file was created?
	return &mlog_raw, nil
}

// */

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
	Raw                   *lib.MatchLogRaw
}

// TODO: make the nametoidx variable an input argument
// TODO: This function is long and unweildly
func NewDRRTStandardMatchLogFromShips(raw *lib.MatchLogRaw, ships []*rsmships.Ship) (*DRRTStandardMatchLog, error) {
	var mlog DRRTStandardMatchLog
	mlog.Raw = raw
	mlog.Timestamp = raw.CreatedTimestamp

	// ASSUMPTION: NOT a free-for-all (Red v Blue alliance)
	// get the match number
	// if the same for red and blue alliances: good!
	redMatchNumber := lib.GetMatchNumberFromAllianceName(raw.StartListings[0].Name, false)
	blueMatchNumber := lib.GetMatchNumberFromAllianceName(raw.StartListings[1].Name, true)
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
		name := lib.ShipAuthorFromCommonNamefmt(shiplsting.Ship)[0]
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
		destroyername := lib.ShipAuthorFromCommonNamefmt(destruction.Ship)[0]
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
		// case lib.MLOG_RED_FACTION:

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
