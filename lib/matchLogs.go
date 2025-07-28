package lib

import (
	"time"
	"bufio"
	"os"
	"regexp"
	"strconv"
	"fmt"
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
	Fdestroyed string
}

type matchLogRaw struct {
	StartListings []MatchLogFleetListing
	ShipListings []MatchLogShipListing
	DestructionListings []MatchLogDestructionListing
	ResultListings []MatchLogFleetListing
	SurvivalListings []MatchLogShipListing
}



func ReadMatchLogAtPath(path string) (DRRTStandardTerseMatchLog, error) {
	mlog_regex_type := regexp.MustCompile(mlog_typeRegexCaptureString)

	mlog_regex_map := map[string]*regexp.Regexp{
		mlog_start:       regexp.MustCompile(mlog_startRegexCaptureString),
		mlog_ship:        regexp.MustCompile(mlog_shipRegexCaptureString),
		mlog_destruction: regexp.MustCompile(mlog_destructionRegexCaptureString),
		mlog_result:      regexp.MustCompile(mlog_resultRegexCaptureString),
		mlog_survival:    regexp.MustCompile(mlog_survivalRegexCaptureString),
	}

	match_log, err := os.Open(path)
	if err != nil {
		return DRRTStandardTerseMatchLog{}, err
	}
	defer match_log.Close()
	// TODO: make function open buffered scanner and read line by line

	var mlog_object DRRTStandardTerseMatchLog

	var mlog_raw matchLogRaw

	match_log_scanner := bufio.NewScanner(match_log)

	mlog_RecordNumber := 0

	for match_log_scanner.Scan() {
		line := match_log_scanner.Text()

		mlog_RecordNumber++

		switch string(mlog_regex_type.Find([]byte(line))) {
		case mlog_start:
			fields := mlog_regex_map[mlog_start].FindAll([]byte(line), -1)

			if len(fields) != length_fleetListingSlice {
				return DRRTStandardTerseMatchLog{}, fmt.Errorf("Line %d: Failed parsing [START] line: %d fields detected when there should have been %d", mlog_RecordNumber, len(fields), length_fleetListingSlice)
			}

			var listing MatchLogFleetListing

			// TODO: try to make this require less statements
			listing.Faction, err = strconv.Atoi(string(fields[0]))
			if err != nil {
				return DRRTStandardTerseMatchLog{}, err
			}

			listing.Name = string(fields[1])

			listing.DamageTaken, err = strconv.Atoi(string(fields[2]))
			if err != nil {
				return DRRTStandardTerseMatchLog{}, err
			}

			listing.DamageInflicted, err = strconv.Atoi(string(fields[3]))
			if err != nil {
				return DRRTStandardTerseMatchLog{}, err
			}

			listing.Alive, err = strconv.Atoi(string(fields[4]))
			if err != nil {
				return DRRTStandardTerseMatchLog{}, err
			}

			mlog_raw.StartListings = append(mlog_raw.StartListings, listing)
		case mlog_ship:

		}
	}
}
