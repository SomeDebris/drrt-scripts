package lib

import (
	"path/filepath"
	"os"
)

var REASSEMBLY_DATA_DIR = filepath.Join(Must(os.UserHomeDir()), ".local", "share", "Reassembly", "data")
