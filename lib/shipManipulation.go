package lib

import (
	"encoding/json"
	"os"
)

type CommandData struct {
	Flags   json.RawMessage `json:"flags,omitempty"`
	Faction int             `json:"faction,omitempty"`
}

type Block struct {
	Id        any          `json:"ident"`
	Offset    [2]float64   `json:"offset"`
	Angle     float64      `json:"angle"`
	Command   *CommandData `json:"command,omitempty"`
	BindingId int          `json:"bindingId,omitempty"`
}

type ShipData struct {
	Name   string          `json:"name"`
	Author string          `json:"author"`
	Color0 json.RawMessage `json:"color0,omitempty"`
	Color1 json.RawMessage `json:"color1,omitempty"`
	Color2 json.RawMessage `json:"color2,omitempty"`
	Wgroup [4]int          `json:"wgroup,omitempty"`
}

type Ship struct {
	Angle    float64    `json:"angle,omitempty"`
	Position [2]float64 `json:"position,omitempty"`
	Data     ShipData   `json:"data"`
	Blocks   []Block    `json:"blocks"`
}

type Fleet struct {
	Blueprints []Ship `json:"blueprints"`
	Color0     any    `json:"color0,omitempty"`
	Color1     any    `json:"color1,omitempty"`
	Color2     any    `json:"color2,omitempty"`
	Faction    int    `json:"faction"`
	Name       string `json:"name"`
}

type UnprocessedShip struct {
	Name json.RawMessage `json:"name"`
}

func IsReassemblyJSONFileFleet(path string) (bool, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}

	var idk UnprocessedShip

	if err := json.Unmarshal([]byte(content), &idk); err != nil {
		return false, err
	}

	if idk.Name == nil {
		return false, nil
	} else {
		return true, nil
	}
}

func UnmarshalShipFromFile(path string) (Ship, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Ship{}, err
	}

	var ship Ship

	if err := json.Unmarshal([]byte(content), &ship); err != nil {
		return Ship{}, err
	}

	return ship, nil
}

func MarshalShipToFile(path string, ship Ship) error {
	b, err := json.Marshal(ship)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, b, 0666); err != nil {
		return err
	}

	return nil
}

func UnmarshalFleetFromFile(path string) (Fleet, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Fleet{}, err
	}

	var fleet Fleet

	if err := json.Unmarshal([]byte(content), &fleet); err != nil {
		return Fleet{}, err
	}

	return fleet, nil
}

func MarshalFleetToFile(path string, ship Fleet) error {
	b, err := json.Marshal(ship)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, b, 0666); err != nil {
		return err
	}

	return nil
}


func FleetFromShips(template Fleet, ships ...Ship) Fleet {
	template.Blueprints = ships

	return template
}

func AssembleAlliance(template Fleet, ships []Ship) Fleet {
	template.Blueprints = ships

	return template
}
