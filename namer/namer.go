package namer

import (
	"drrt-scripts/lib"
	"encoding/flag"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"github.com/SomeDebris/rsmships-go"
	"bufio"
)

// Goal is to make a command line tool for
// - setting the ship and author name properly
// - saving the ship to the right directory
// - 

func main() {
	author_arg := flag.String("author", "", "Declare the name of the ship's author.")
	name_arg := flag.String("name", "Unnamed Spaceship", "Declare the name of the ship.")
	
}
