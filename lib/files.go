package lib

import (
	"os"
	"path/filepath"
)

func Path_exists(path string) bool {
    _, err := os.Stat(path)

    if err != nil {
        if os.IsNotExist(err) {
            return false
        } else {
            // It may not exist!
            return false
        }
    }

    return true
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
