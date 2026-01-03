package lib

import (
	// "github.com/SomeDebris/rsmships-go"
)

type DRRTShipStats struct {
	Name                string `json:"name"`
	Author              string
	RankPoints          float64
	NumberMatchesPlayed int `json:"matchesPlayed"`
	QualsSeed           int
}

func (m *DRRTShipStats) RankingScore() float64 {
	return m.RankPoints / float64(m.NumberMatchesPlayed)
}
