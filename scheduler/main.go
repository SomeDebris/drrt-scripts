package main

import (
    "drrt-scripts/lib"
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "sort"
    "log"
    //"json"
    "encoding/csv"
    "strconv"
    "strings"
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

/**
* schedule array
* whether ship is participating as surrogate
* any error recieved
*/
func get_schedule_from_path(path string) ([][]int, [][]bool, error) {
    schedule_bytes, err := os.ReadFile(path)
    if err != nil {
        log.Fatal("Cannot find schedule file: ", err)
    }

    schedule_string := string(schedule_bytes)

    r := csv.NewReader(strings.NewReader(schedule_string))

    records, err := r.ReadAll()
    if err != nil {
        log.Fatal(err)
    }

    schedule    := make([][]int, 0)
    surrogates  := make([][]bool, 0)

    // Parse out the strings in the schedules to integers
    for _, match := range records {
        ships_in_match      := make([]int, len(match))
        surrogates_in_match := make([]bool, len(match))

        for i, ship := range match {
            if strings.ContainsAny(ship, "*") {
                surrogates_in_match[i] = true
            } else {
                surrogates_in_match[i] = false
            }

            ship_noasterisk := strings.ReplaceAll(ship, "*", "")

            ships_in_match[i], err = strconv.Atoi(ship_noasterisk)
            if err != nil {
                return nil, nil, err
            }
        }
        schedule    = append(schedule, ships_in_match)
        surrogates  = append(surrogates, surrogates_in_match)
    }
    
    return schedule, surrogates, nil
}

// assemble alliance
// func assemble_alliance(

func main() {
    drrt_directory_arg := flag.String("drrt-directory", ".", "Set the directory the DRRT will be run in.")
    ships_per_alliance_arg := flag.Int("n", 3, "Set the number of ships per alliance.")

    flag.Parse()

    fmt.Println("drrt_directory:     ", *drrt_directory_arg)
    fmt.Println("ships_per_alliance: ", *ships_per_alliance_arg)

    ships_directory := filepath.Join(*drrt_directory_arg, "Ships")
    quals_directory := filepath.Join(*drrt_directory_arg, "Qualifications")
    stags_directory := filepath.Join(*drrt_directory_arg, "Staging")
    playf_directory := filepath.Join(*drrt_directory_arg, "Playoffs")
    scrpt_directory := filepath.Join(*drrt_directory_arg, "drrt-scripts")

    drrt_subdirectories := []string{ships_directory, quals_directory,
            stags_directory, playf_directory, scrpt_directory}
    for _, dir := range drrt_subdirectories {
        err := os.MkdirAll(dir, os.ModePerm)
        if err != nil {
            if os.IsExist(err) {
                fmt.Printf("Directory '%s' already exists.\n", dir)
            } else {
                fmt.Printf("Error making directory '%s': %s\n", dir, err)
                os.Exit(1)
            }
        }
    }


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
    // sch_out_filepath := filepath.Join(scrpt_directory, "selected_schedule.csv")
    // sch_out_filepath_no_asterisks := filepath.Join(scrpt_directory, ".no_asterisks.csv")

    schedule, surrogates, err := get_schedule_from_path(sch_in_filepath)  
    fmt.Printf("Schedule has %d matches.\n", len(schedule))
    fmt.Printf("Assembling Alliances.")

    for i, match := range schedule {
        fmt.Printf("match %d: ", i + 1)
        for _, ship := range match {
            fmt.Printf("%d ", ship)
        }
        fmt.Printf("\n")
    }

    for i, match := range surrogates {
        fmt.Printf("match %d: ", i + 1)
        for _, surrogate := range match {
            fmt.Printf("%d ", surrogate)
        }
        fmt.Printf("\n")
    }
}
