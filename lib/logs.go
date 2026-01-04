package lib

import (
	"log/slog"
	"bufio"
	"os"
	"io"
)

// opens a log file at the specified location. Selects standard error if 
func DRRTLoggerPreferences(path string, log_lvl slog.Level) (*os.File, *bufio.Writer, error) {
	var handler *slog.TextHandler
	// Set the default log tech to the logger we set before
	defer slog.SetDefault(slog.New(handler))

	var log_file os.File
	var writer io.Writer
	var bufferedwriter *bufio.Writer

	handler_options := &slog.HandlerOptions{Level: log_lvl}
	if path != "" {
		log_file, err := os.Create(path)
		if err != nil {
			// log.Fatalf("Could not open log file '%s': %v", path, err)
			return nil, nil, err
		}
		bufferedwriter = bufio.NewWriter(log_file)
		writer = bufferedwriter
	} else {
		log_file = *os.Stderr
		writer = os.Stderr
	}
	handler = slog.NewTextHandler(writer, handler_options)

	return &log_file, bufferedwriter, nil
}
