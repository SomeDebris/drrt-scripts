package lib

import (
	"log/slog"
	"fmt"
)


// If an error occurs during the parse of a match log, throw this error.
type MatchLogFieldError struct {
	message    string
	field      string
	event      string
	line       string
	lineNumber int
	path       string
}

func (e *MatchLogFieldError) Error() string {
	return fmt.Sprintf("Could not parse field: %v", e.message)
}
func (e *MatchLogFieldError) LogError(logger *slog.Logger) {
	logger.Error(e.message, "line", e.line, "field", e.field, "event", e.event, "lineNumber", e.lineNumber, "path", e.path)
}
func (e *MatchLogFieldError) AddContext(lineNumber int, path string) {
	e.lineNumber = lineNumber
	e.path = path
}

type MatchLogRegexError struct {
	event      string
	line       string
	lineNumber int
	path       string
	regex      string
}

func (e *MatchLogRegexError) Error() string {
	return "Could not match line against regex."
}
func (e *MatchLogRegexError) LogError(logger *slog.Logger) {
	logger.Error(e.Error(), "line", e.line, "event", e.event, "lineNumber", e.lineNumber, "path", e.path, "regex", e.regex)
}
func (e *MatchLogRegexError) AddContext(lineNumber int, path string) {
	e.lineNumber = lineNumber
	e.path = path
}

type MatchLogAllianceLengthMismatch struct {
	redAllianceLength  int
	blueAllianceLength int
}
func (e *MatchLogAllianceLengthMismatch) Error() string {
	return "Red and blue alliances have different lengths."
}
func (e *MatchLogAllianceLengthMismatch) LogError(logger *slog.Logger) {
	logger.Error(e.Error(), "redAllianceLength", e.redAllianceLength, "blueAllianceLength", e.blueAllianceLength)
}
func (e *MatchLogAllianceLengthMismatch) AddContext(redAllianceLength, blueAllianceLength int) {
	e.redAllianceLength = redAllianceLength
	e.blueAllianceLength = blueAllianceLength
}

