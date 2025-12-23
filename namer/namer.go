package namer

import (
	// "drrt-scripts/lib"
	"flag"
	// "log"
	"fmt"
	"log/slog"
	"os"
	// "path/filepath"
	"github.com/SomeDebris/rsmships-go"
	// "bufio"
)

// Goal is to make a command line tool for
// - setting the ship and author name properly
// - saving the ship to the right directory
// - 

func main() {
	exit_code := 0
	defer os.Exit(exit_code)

	author_arg := flag.String("author", "", "Declare the name of the ship's author.")
	name_arg := flag.String("name", "Unnamed Spaceship", "Declare the name of the ship.")
	copy_arg := flag.Bool("copy", true, "Set whether the file should be copied with new name and location")
	suffix_arg := flag.String("suffix", "2026W", "Declare a suffix for filenames. Occurs before file extension.")
	
	// ParseFlags()
	flag.Parse()

	if flag.NArg() < 1 {
		slog.Error("Expected filename of ship as first positional argument, but no positional arguments specified.")
		exit_code = 1
		return
	} else if flag.NArg() > 1 {
		slog.Error("Expected single positional argument (ship filename), but multiple positinal arguments specified.")
		exit_code = 1
		return
	}

	filename := flag.Arg(0)
	var ship rsmships.Ship

	// determine if ship file is fleet
	isfleet, err := rsmships.IsReassemblyJSONFileFleet(filename)
	if err != nil {
		slog.Error("Failed preparation for unmarshalling ship", "path", filename, "err", err)
		exit_code = 1
		return
	}

	// extract the ship in the correct method
	if isfleet {
		fleet, err := rsmships.UnmarshalFleetFromFile(filename)
		if err != nil {
			slog.Error("Failed unmarshalling fleet", "path", filename, "err", err)
			exit_code = 1
			return
		}

		ship = *fleet.Blueprints[0]
		slog.Info("Umarshalled ship from fleet.", "name", ship.Data.Name, "author", ship.Data.Author, "fleet.Name", fleet.Name)
	} else {
		ship, err = rsmships.UnmarshalShipFromFile(filename)
		if err != nil {
			slog.Error("Failed preparation for unmarshalling ship", "path", filename, "err", err)
			exit_code = 1
			return
		}

		slog.Info("Unmarshalled ship", "name", ship.Data.Name, "author", ship.Data.Author)
	}

	// the ship now exists. set its author and name
	ship.Data.Name = *name_arg
	ship.Data.Author = *author_arg
	
	out_filename = fmt.Sprintf("%s_[by_%s]_%s.json", *name_arg, *author_arg, *suffix_arg)
}
