package lib

import (
	"github.com/SomeDebris/rsmships-go"
)

type DRRTShipStats struct {
	Name                string  `json:"name"`
	Author              string  `json:"author"`
	RankPoints          float64 `json:"rps"`
	NumberMatchesPlayed int     `json:"matchesPlayed"`
	QualsSeed           int
	Faction             int
	P                   int
	ShipData            *rsmships.Ship
}


func (m *DRRTShipStats) RankingScore() float64 {
	return m.RankPoints / float64(m.NumberMatchesPlayed)
}
