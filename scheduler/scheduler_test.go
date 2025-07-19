package main

import (
	"os"
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

	b, err := json.MarshalIndent(cornstar, "", "\t")
	if err != nil {
		t.Errorf("Couldn't marshal again: %v", err)
	}

	t.Log(string(b))
}
