package main

import (
	"fmt"
	"os"

    "path/filepath"
	"drrt-scripts/lib"
	"flag"
    "sort"
)

func get_inspected_ship_paths(dir string) ([]string, error) {
    var ship_files []string

    f, err := os.Open(dir)
    if err != nil {
        return ship_files, err
    }

    file_info, err := f.Readdir(-1)
    f.Close()
    if err != nil {
        return ship_files, err
    }

    sort.Slice(file_info, func(i,j int) bool {
        return file_info[i].ModTime().Before(file_info[j].ModTime())
    })

    for _, file := range file_info {
        ship_files = append(ship_files, file.Name())
    }

    return ship_files, nil
}

func main() {
    drrt_directory_arg := flag.String("drrt-directory", ".", "Set the directory the DRRT will be run in.")
    ships_per_alliance_arg := flag.Int("n", 3, "Set the number of ships per alliance.")

    flag.Parse()

    fmt.Println("drrt_directory:     ", *drrt_directory_arg)
    fmt.Println("ships_per_alliance: ", *ships_per_alliance_arg)

    drrt_subdirectories := []string{"Playoffs", "Qualifications", "Ships", "Staging"}
    for _, dir := range drrt_subdirectories {
        err := os.MkdirAll(filepath.Join(*drrt_directory_arg, dir), os.ModePerm)
        if err != nil {
            if os.IsExist(err) {
                fmt.Printf("Directory '%s' already exists.\n", dir)
            } else {
                fmt.Printf("Error making directory '%s'.\n", dir)
            }
        }
    }

    // Empty the contents of the Qualifications directory
    err := lib.Remove_directory_contents(filepath.Join(*drrt_directory_arg, "Qualifications"))
    if err != nil {
        fmt.Printf("cannot remove contents of the 'Qualifications' directory.\n")
    }

        
}
