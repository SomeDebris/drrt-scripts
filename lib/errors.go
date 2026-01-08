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
	return fmt.Sprintf("Could not match line %d against regex: %s", e.lineNumber, e.line)
}
func (e *MatchLogRegexError) LogError(logger *slog.Logger) {
	logger.Error(e.Error(), "line", e.line, "event", e.event, "lineNumber", e.lineNumber, "path", e.path, "regex", e.regex)
}
func (e *MatchLogRegexError) AddContext(lineNumber int, path string) {
	e.lineNumber = lineNumber
	e.path = path
}

type MatchLogAllianceLengthMismatchError struct {
	redAllianceLength  int
	blueAllianceLength int
}
func (e *MatchLogAllianceLengthMismatchError) Error() string {
	return "Red and blue alliances have different lengths."
}
func (e *MatchLogAllianceLengthMismatchError) LogError(logger *slog.Logger) {
	logger.Error(e.Error(), "redAllianceLength", e.redAllianceLength, "blueAllianceLength", e.blueAllianceLength)
}
func (e *MatchLogAllianceLengthMismatchError) AddContext(redAllianceLength, blueAllianceLength int) {
	e.redAllianceLength = redAllianceLength
	e.blueAllianceLength = blueAllianceLength
}

type MatchLogAllianceMatchNumberMismatchError struct {
	redAllianceMatchNumber  int
	blueAllianceMatchNumber int
}
func (e *MatchLogAllianceMatchNumberMismatchError) Error() string {
	return fmt.Sprintf("Red and blue alliances have different match numbers: red=\"%d\" blue=\"%d\"", e.redAllianceMatchNumber, e.blueAllianceMatchNumber)
}
func (e *MatchLogAllianceMatchNumberMismatchError) LogError(logger *slog.Logger) {
	logger.Error(e.Error(), "redAllianceMatchNumber", e.redAllianceMatchNumber, "blueAllianceMatchNumber", e.blueAllianceMatchNumber)
}
func (e *MatchLogAllianceMatchNumberMismatchError) AddContext(redAllianceMatchNumber, blueAllianceMatchNumber int) {
	e.redAllianceMatchNumber = redAllianceMatchNumber
	e.blueAllianceMatchNumber = blueAllianceMatchNumber
}

type MatchLogIncompleteError struct {
	message     string
	path        string
	matchNumber int
}
func (e *MatchLogIncompleteError) Error() string {
	return "Match log is incomplete: " + e.message
}
func (e *MatchLogIncompleteError) LogError(logger *slog.Logger) {
	logger.Error(e.Error(), "msg", e.message, "path", e.path, "matchNumber", e.matchNumber)
}
func (e *MatchLogIncompleteError) AddContext(message, path string, matchNumber int) {
	e.message = message
	e.path = path
	e.matchNumber = matchNumber
}

