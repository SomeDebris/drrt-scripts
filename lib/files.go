package lib

import (
	"os"
	"strings"
	"path/filepath"
)

var Replacer_Out_Filename = strings.NewReplacer(
	` `, `_`,
	`/`, `-`,
	`\`, `-`,
	`?`, `-`,
	`*`, `⋆`,
	`:`, `꞉`,
	`%`, `-`,
	`|`, `∣`,
	`"`, `''`,
	`<`, `lt`,
	`>`, `gt`,
	`.`, `p`,
	`=`, `-`)

func Path_exists(path string) (bool, error) {
    _, err := os.Stat(path)

    if err != nil {
        if os.IsNotExist(err) {
            return false, err
        } else {
            // It may not exist!
            return false, err
        }
    }

    return true, nil
}

func Remove_directory_contents(directory string) error {
    d, err := os.Open(directory)
    if err != nil {
        return err
    }

    defer d.Close()

    names, err := d.Readdirnames(-1)

    if err != nil {
        return err
    }

    for _, name := range names {
        err = os.RemoveAll(filepath.Join(directory, name))
        if err != nil {
            return err
        }
    }

    return nil
}
