package main

import (
	"os"
	"fmt"
	"path/filepath"
	"testing"
	"encoding/json"
	"github.com/SomeDebris/rsmships-go"
)

func TestShipUnmarshal(t *testing.T) {
	ship_target := filepath.Join("..", "test-ships", "Chernobyl_[by_Kaboom]_2023W.json")
	content, err := os.ReadFile(ship_target)
	if err != nil {
		t.Errorf("Couldn't read file: %v", err)
	}

	var cornstar rsmships.Ship
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

// func TestFleetUnmarshal(t *testing.T) {
// 	ship_target := filepath.Join("test-ships", "Reassembly_Point_filler_2.0_20250719_12.48.44.PM_530P.json")
//
// 	content, err := os.ReadFile(ship_target)
// 	if err != nil {
// 		t.Errorf("Couldn't read file: %v", err)
// 	}
//
// 	var debsonder rsmships.Fleet
// 	if err := json.Unmarshal([]byte(content), &debsonder); err != nil {
// 		t.Errorf("Couldn't unmarshal: %v", err)
// 	}
//
// 	t.Logf("name: %s\n", debsonder.Name)
//
// 	b, err := json.Marshal(debsonder)
// 	if err != nil {
// 		t.Errorf("Couldn't marshal again: %v", err)
// 	}
//
// 	filename := fmt.Sprintf("%s_[by %s]_re-marshalled.json", debsonder.Blueprints[0].Data.Name, debsonder.Blueprints[0].Data.Author)
// 	if err := os.WriteFile(filename, b, 0666); err != nil {
// 		t.Errorf("Couldn't save file: %v", err)
// 	}
// }
