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
	ANSI_RESET                    = "\033[0m"
	ANSI_RED                      = "\033[31m"
	ANSI_GREEN                    = "\033[32m"
	ANSI_YELLOW                   = "\033[33m"
	ANSI_BLUE                     = "\033[34m"
	ANSI_MAGENTA                  = "\033[35m"
	ANSI_CYAN                     = "\033[36m"
	ANSI_GRAY                     = "\033[37m"
	ANSI_WHITE                    = "\033[97m"

	ANSI_HEADER    = "\033[95m"
	ANSI_OKBLUE    = "\033[94m"
	ANSI_OKCYAN    = "\033[96m"
	ANSI_OKGREEN   = "\033[92m"
	ANSI_WARNING   = "\033[93m"
	ANSI_FAIL      = "\033[91m"
	ANSI_ENDC      = "\033[0m"
	ANSI_BOLD      = "\033[1m"
	ANSI_UNDERLINE = "\033[4m"
)
