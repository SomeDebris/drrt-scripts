package lib

import (
	"time"
	"bufio"
	"os"
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


func ReadMatchLogAtPath(path string) (DRRTStandardTerseMatchLog, error) {
	match_log, err := os.Open(path)
	if err != nil {
		return DRRTStandardTerseMatchLog{}, err
	}
	defer match_log.Close()
	// TODO: make function open buffered scanner and read line by line
}
