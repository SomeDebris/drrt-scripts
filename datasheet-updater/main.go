package main

import (
	// "context"
	// "encoding/json"
	// "fmt"
	"log"
	"log/slog"
	// "net/http"
	"flag"
	"os"
	// "bufio"
	"drrt-scripts/lib"
	"github.com/SomeDebris/rsmships-go"
	"time"

	// "golang.org/x/oauth2"
	// "golang.org/x/oauth2/google"
	// "google.golang.org/api/option"
	// "google.golang.org/api/sheets/v4"
)
type matchResult int
const (
	WinDestruction matchResult = iota
	WinPoints
	Loss
)

type matchPerformance struct {
	Ship             *rsmships.Ship
	Destructions     uint
	RankPointsEarned uint
	Result           matchResult
	Survived         bool
}
func (m *matchPerformance) toSheetsRow() []any {
	output := make([]any, 7)
	output[0] = m.Ship.Data.Name
	output[1] = m.Destructions
	output[2] = m.RankPointsEarned

	output[3] = 0
	output[4] = 0
	output[5] = 0
	switch m.Result {
	case WinDestruction:
		output[3] = 1
	case WinPoints:
		output[4] = 1
	case Loss:
		output[5] = 1
	}

	if m.Survived {
		output[6] = 1
	} else {
		output[6] = 0
	}
	return output
}

type DRRTStandardTerseMatchLog struct {
	MatchNumber               int
	Timestamp                 time.Time
	RedAlliance               []*rsmships.Ship
	BlueAlliance              []*rsmships.Ship
	Record                    []*matchPerformance
	RedPointsDamageInflicted  int
	RedPointsDamageTaken      int
	BluePointsDamageInflicted int
	BluePointsDamageTaken     int
	Raw                       *lib.MatchLogRaw
}
func NewDRRTStandardTerseMatchLog(raw *lib.MatchLogRaw) (*DRRTStandardTerseMatchLog, error) {
	var mlog DRRTStandardTerseMatchLog
	mlog.Raw = raw
	mlog.Timestamp = raw.CreatedTimestamp
	return &mlog, nil
}

func main() {
// boilerplate stuff for exit codes
	exit_code := 0
	defer os.Exit(exit_code)

	log_lvl := slog.LevelInfo

	log_file_name := flag.String("log-filename", "", "Send log messages to a file. If not set, log to standard error.")
	flag.Parse()

	// same tech used in scheduler/main.go
	log_ref, log_writer_ref, err := lib.DRRTLoggerPreferences(*log_file_name, log_lvl)
	if err != nil {
		log.Fatalf("Could not open log file '%s': %v", *log_file_name, err)
	}
	logfile := *log_ref
	defer logfile.Close()
	if log_writer_ref != nil {
		logwriter := *log_writer_ref
		defer logwriter.Flush()
	}

	drrtdatasheet := lib.NewDRRTDatasheetDefaults()
	if drrtdatasheet.Service == nil {
		slog.Error("Failed to get google sheets service.", "err", err)
		exit_code = 1
		return
	}
	// TODO: use the lib.MatchSchedule type
	scheduleindices, err := drrtdatasheet.GetMatchSheduleValues()
	if err != nil {
		slog.Error("Failed to get match schdule.", "err", err)
		exit_code = 1
		return
	}

	
}
