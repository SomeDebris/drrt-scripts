package lib

import (
	"time"
	"bufio"
	"os"
	"regexp"
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

	match_log_scanner := bufio.NewScanner(match_log)
	for match_log_scanner.Scan() {
		line := match_log_scanner.Text()

		switch string(mlog_regex_type.Find([]byte(line))) {
		case mlog_start:
			// TODO
		}
	}
}
