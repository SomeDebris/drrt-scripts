package main

import (
	"os"
	"fmt"
	"path/filepath"
	"testing"
	"drrt-scripts/lib"
	"encoding/json"
)

func TestShipUnmarshal(t *testing.T) {
	ship_target := filepath.Join("..", "Ships", "Chernobyl_[by_Kaboom]_2023W.json")
	
	content, err := os.ReadFile(ship_target)
	if err != nil {
		t.Errorf("Couldn't read file: %v", err)
	}

	var cornstar lib.Ship
	if err := json.Unmarshal([]byte(content), &cornstar); err != nil {
		t.Errorf("Couldn't unmarshal: %v", err)
	}
	
	t.Logf("name: %s\n", cornstar.Data.Name)

	b, err := json.Marshal(cornstar)
	if err != nil {
		t.Errorf("Couldn't marshal again: %v", err)
	}
	
	filename := fmt.Sprintf("%s_[by %s]_re-marshalled.json", cornstar.Data.Name, cornstar.Data.Author)
	if err := os.WriteFile(filename, b, 0666); err != nil {
		t.Errorf("Couldn't save file: %v", err)
	}
}

func TestFleetUnmarshal(t *testing.T) {
	ship_target := filepath.Join("..", "Ships", "Reassembly_Point_filler_2.0_20250719_12.48.44.PM_530P.json")
	
	content, err := os.ReadFile(ship_target)
	if err != nil {
		t.Errorf("Couldn't read file: %v", err)
	}

	var debsonder lib.Fleet
	if err := json.Unmarshal([]byte(content), &debsonder); err != nil {
		t.Errorf("Couldn't unmarshal: %v", err)
	}
	
	t.Logf("name: %s\n", debsonder.Name)

	b, err := json.Marshal(debsonder)
	if err != nil {
		t.Errorf("Couldn't marshal again: %v", err)
	}

	filename := fmt.Sprintf("%s_[by %s]_re-marshalled.json", debsonder.Blueprints[0].Data.Name, debsonder.Blueprints[0].Data.Author)
	if err := os.WriteFile(filename, b, 0666); err != nil {
		t.Errorf("Couldn't save file: %v", err)
	}
}
