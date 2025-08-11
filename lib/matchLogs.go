package lib

import (
	"time"
	"bufio"
	"os"
	"regexp"
	"strconv"
	"fmt"
	// "slog"
	"errors"
	"sync"
)

const (
	mlog_typeRegexCaptureString = `^\[([A-Z]+)\]`
	mlog_startRegexCaptureString = `^\[START\] faction:\{([0-9]+)\} name:\{(.*)\} DT:\{([0-9]*)\} DI:\{([0-9])*\} alive:\{([0-9]*)\}$`
	mlog_shipRegexCaptureString = `^\[SHIP\] faction:\{([0-9]+)\} ship:\{(.*)\}$`
	mlog_destructionRegexCaptureString = `^\[DESTRUCTION\] ship:\{(.*)\} fship:\{([0-9]*)\} destroyed:\{(.*)\} fdestroyed:\{([0-9]*)\}$`
	mlog_resultRegexCaptureString = `^\[RESULT\] faction:\{([0-9]+)\} name:\{(.*)\} DT:\{([0-9]*)\} DI:\{([0-9])*\} alive:\{([0-9]*)\}$`
	mlog_survivalRegexCaptureString = `^\[SURVIVAL\] faction:\{([0-9]+)\} ship:\{(.*)\}$`

	mlog_start = `START`
	mlog_ship = `SHIP`
	mlog_destruction = `DESTRUCTION`
	mlog_result = `RESULT`
	mlog_survival = `SURVIVAL`

	length_fleetListingSlice = 5
	length_resultListingSlice = 5

	length_shipListingSlice = 2
	length_survivalListingSlice = 2

	length_destructionListingSlice = 4
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

var mlog_regex_map = map[string]*regexp.Regexp{
		mlog_start:       regexp.MustCompile(mlog_startRegexCaptureString),
		mlog_ship:        regexp.MustCompile(mlog_shipRegexCaptureString),
		mlog_destruction: regexp.MustCompile(mlog_destructionRegexCaptureString),
		mlog_result:      regexp.MustCompile(mlog_resultRegexCaptureString),
		mlog_survival:    regexp.MustCompile(mlog_survivalRegexCaptureString),
	}

type DRRTStandardTerseMatchLog struct {
	MatchNumber               int
	Timestamp                 time.Time
	RedAlliance               []string
	BlueAlliance              []string
	Destructions              map[string]string
	RedPointsDamageInflicted  int
	RedPointsDamageTaken      int
	BluePointsDamageInflicted int
	BluePointsDamageTaken     int
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
	Ship string
}

type MatchLogDestructionListing struct {
	Ship       string
	Fship      int
	Destroyed  string
	Fdestroyed int
}

type matchLogRaw struct {
	StartListings []MatchLogFleetListing
	ShipListings []MatchLogShipListing
	DestructionListings []MatchLogDestructionListing
	ResultListings []MatchLogFleetListing
	SurvivalListings []MatchLogShipListing
}

var (
	matchLogRawMutex_start sync.Mutex
	matchLogRawMutex_ship sync.Mutex
	matchLogRawMutex_destruction sync.Mutex
	matchLogRawMutex_result sync.Mutex
	matchLogRawMutex_survival sync.Mutex
)

// Append a START event.
func (mlograw *matchLogRaw) appendStart(record MatchLogFleetListing, mutex *sync.Mutex) {
	mutex.Lock()
	mlograw.StartListings = append(mlograw.StartListings, record)
	mutex.Unlock()
}
// Append a SHIP event.
func (mlograw *matchLogRaw) appendShip(record MatchLogShipListing, mutex *sync.Mutex) {
	mutex.Lock()
	mlograw.ShipListings = append(mlograw.ShipListings, record)
	mutex.Unlock()
}
// Append a DESTRUCTION event.
func (mlograw *matchLogRaw) appendDestruction(record MatchLogDestructionListing, mutex *sync.Mutex) {
	mutex.Lock()
	mlograw.DestructionListings = append(mlograw.DestructionListings, record)
	mutex.Unlock()
}
// Append a RESULT event.
func (mlograw *matchLogRaw) appendResult(record MatchLogFleetListing, mutex *sync.Mutex) {
	mutex.Lock()
	mlograw.ResultListings = append(mlograw.ResultListings, record)
	mutex.Unlock()
}
// Append a SURVIVAL event.
func (mlograw *matchLogRaw) appendSurvival(record MatchLogShipListing, mutex *sync.Mutex) {
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
			event:      mlog_ship,
			line:       line,
			regex:      regex.String(),
		}
	}

	listing.Ship = fields[1]
	listing.Fship, err = strconv.Atoi(fields[2])
	if err != nil {
		return listing, &MatchLogFieldError{
			message: err.Error(),
			event: mlog_destruction,
			field: "fship",
			line: line,
		}
	}
	listing.Destroyed = fields[3]
	listing.Fdestroyed, err = strconv.Atoi(fields[4])
	if err != nil {
		return listing, &MatchLogFieldError{
			message: err.Error(),
			event: mlog_destruction,
			field: "fdestroyed",
			line: line,
		}
	}

	return listing, nil
}


func ReadMatchLogAtPath(path string) (DRRTStandardTerseMatchLog, error) {
	mlog_regex_type := regexp.MustCompile(mlog_typeRegexCaptureString)

	match_log, err := os.Open(path)
	if err != nil {
		return DRRTStandardTerseMatchLog{}, err
	}
	defer match_log.Close()

	var mlog_object DRRTStandardTerseMatchLog

	var mlog_raw matchLogRaw

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
				return mlog_object, &MatchLogRegexError{
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
				return mlog_object, &MatchLogFieldError{
					field: "faction",
					event: mlog_start,
					lineNumber:  mlog_RecordNumber,
					line: line,
					path: path,
				}
			}
			listing.Name = fields[2]
			listing.DamageTaken, err = strconv.Atoi(fields[3])
			if err != nil {
				return mlog_object, &MatchLogFieldError{
					message: err.Error(),
					field: "DT",
					event: mlog_start,
					lineNumber:  mlog_RecordNumber,
					line: line,
					path: path,
				}
			}
			listing.DamageInflicted, err = strconv.Atoi(fields[4])
			if err != nil {
				return mlog_object, &MatchLogFieldError{
					message: err.Error(),
					field: "DI",
					event: mlog_start,
					lineNumber:  mlog_RecordNumber,
					line: line,
					path: path,
				}
			}
			listing.Alive, err = strconv.Atoi(fields[5])
			if err != nil {
				return mlog_object, &MatchLogFieldError{
					message: err.Error(),
					field: "alive",
					event: mlog_start,
					lineNumber:  mlog_RecordNumber,
					line: line,
					path: path,
				}
			}
			mlog_raw.appendStart(listing, &matchLogRawMutex_start)
		// If line is a [SHIP] line:
		case mlog_ship:
			fields := mlog_regex_map[mlog_ship].FindStringSubmatch(line)
			if fields == nil {
				return mlog_object, &MatchLogRegexError{
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
				return mlog_object, &MatchLogFieldError{
					message: err.Error(),
					event: mlog_ship,
					field: "faction",
					line: line,
					lineNumber: mlog_RecordNumber,
					path: path,
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
				return mlog_object, err
			}
			mlog_raw.appendDestruction(listing, &matchLogRawMutex_destruction)
		}
	}
}

