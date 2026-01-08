//go:build !windows
package lib

import (
	"bufio"
	"errors"
	"io"
	"log/slog"
	"os"
)

const (
	DRRT_MLOG_SIGNAL_PIPE_PATH = `/tmp/drrt_mlog_signal_pipe`
)

// Make function for creating a named pipe in the /tmp directory with a global filename

// make function for returning value from named pipe
func ReadDRRTMlogPipe(path string) (string, error) {
	slog.Info("Trying to open mlog signal pipe.", "path", path)
	pipe, err := os.OpenFile(path, os.O_RDONLY, os.ModeNamedPipe|0600)
	if err != nil {
		return "", err
	}
	defer pipe.Close()

	var line []byte

	reader := bufio.NewReader(pipe)
	for {
		line, err = reader.ReadBytes('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				slog.Info("Done reading from pipe.")
				break
			}
			slog.Error("Failed to read from pipe.", "err", err)
			return "", err
		}
	}
	return string(line), nil
}
