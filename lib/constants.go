package lib

import (
	"encoding/json"
)

type CommandData struct {
	Flags   json.RawMessage `json:"flags,omitempty"`
	Faction int             `json:"faction,omitempty"`
}

type Block struct {
	Id      any          `json:"ident"`
	Offset  [2]float64   `json:"offset"`
	Angle   float64      `json:"angle"`
	Command *CommandData `json:"command,omitempty"`
	BindingId int `json:"bindingId,omitempty"`
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
	Blueprints []Ship          `json:"blueprints"`
	Color0     json.RawMessage `json:"color0,omitempty"`
	Color1     json.RawMessage `json:"color1,omitempty"`
	Color2     json.RawMessage `json:"color2,omitempty"`
	Faction    int             `json:"faction"`
	Name       string          `json:"name"`
}
