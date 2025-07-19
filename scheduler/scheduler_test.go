package main

import (
	"os"
	"log"
	"path/filepath"
	"testing"
	"drrt-scripts/lib"
	"encoding/json"
)

func TestShipUnmarshal(t *testing.T) {
	ship_target := filepath.Join("Ships", "CornStar_AF-75_[by_Kepler]_2023W.json")

	content, err := os.ReadFile(ship_target)
	if err != nil {
		t.Errorf("Couldn't read file: %v", err)
	}

	var cornstar lib.Ship
	if err := json.Unmarshal([]byte(content), &cornstar); err != nil {
		t.Errorf("Couldn't unmarshal: %v", err)
	}
	
	log.Printf("name: %s\n", cornstar.Data.Name)
}
