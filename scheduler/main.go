package main

import (
    "os"
    "fmt"
    "path/filepath"
    "flag"
)

func main() {
    drrt_directory_arg := flag.String("drrt-directory", ".", "Set the directory the DRRT will be run in.")
    ships_per_alliance_arg := flag.Int("n", 3, "Set the number of ships per alliance.")

    flag.Parse()

    fmt.Println("drrt_directory:     ", *drrt_directory_arg)
    fmt.Println("ships_per_alliance: ", *ships_per_alliance_arg)

    drrt_subdirectories := []string{"Playoffs", "Qualifications", "Ships", "Staging"}
    for _, dir := range drrt_subdirectories {
        err := os.MkdirAll(dir, os.ModePerm)
        if err != nil {
            if os.IsExist(err) {
                
            } else {
                fmt.Printf("Error making directory '%s'.\n", dir)
            }
            continue
        }


    }

    


    
}
