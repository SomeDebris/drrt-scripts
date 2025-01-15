package main

import (
	"drrt-scripts/lib"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
    "log"
    "json"
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
                fmt.Printf("Error making directory '%s': %s\n", dir, err)
                os.Exit(1)
            }
        }
    }

    ships_directory := filepath.Join(*drrt_directory_arg, "Ships")
    quals_directory := filepath.Join(*drrt_directory_arg, "Qualifications")
    stags_directory := filepath.Join(*drrt_directory_arg, "Staging")
    playf_directory := filepath.Join(*drrt_directory_arg, "Playoffs")
    scrpt_directory := filepath.Join(*drrt_directory_arg, "drrt-scripts")

    // Empty the contents of the Qualifications directory
    err := lib.Remove_directory_contents(quals_directory)
    if err != nil {
        log.Fatalf("Cannot remove contents of Qualifications directory: %s\n", err)
    }

    ships, err := get_inspected_ship_paths(ships_directory) 
    if err != nil {
        log.Fatalf("Cannot get inspected ship paths: %s\n", err)
    }
    
    fmt.Printf("Found %s ship files!\n", len(ships))

    if len(ships) < (*ships_per_alliance_arg * 2) {
        log.Fatalf("%d is lesser than the minimum number of ships (%d).\n",
                   len(ships), *ships_per_alliance_arg)
    }

    sch_in_filename := fmt.Sprintf("%d_%dv%d.csv", len(ships), *ships_per_alliance_arg, *ships_per_alliance_arg)
    sch_in_filepath := filepath.Join(scrpt_directory, "schedules", "out", sch_in_filename)
}
