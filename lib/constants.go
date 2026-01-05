package lib

import (
	"github.com/SomeDebris/rsmships-go"
)

var RED_ALLIANCE_TEMPLATE rsmships.Fleet = rsmships.Fleet{
	Blueprints: []*rsmships.Ship{},
	Color0:     0xBAA01E,
	Color1:     0x681818,
	Color2:     0x000000,
	Faction:    100,
	Name:       "The Red Alliance",
}
var BLUE_ALLIANCE_TEMPLATE rsmships.Fleet = rsmships.Fleet{
	Blueprints: []*rsmships.Ship{},
	Color0:     0x0AA879,
	Color1:     0x222D84,
	Color2:     0x000000,
	Faction:    101,
	Name:       "The Blue Alliance",
}

const (
	REASSEMBLY_FILE_TIMESTAMP_FMT = "20060102_03.04.05.PM"
	MLOG_PREFIX                   = `MLOG_`
	MLOG_EXTENSION                = `.txt`
)
