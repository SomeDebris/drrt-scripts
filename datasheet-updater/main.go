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

	// "golang.org/x/oauth2"
	// "golang.org/x/oauth2/google"
	// "google.golang.org/api/option"
	// "google.golang.org/api/sheets/v4"
)

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
}
