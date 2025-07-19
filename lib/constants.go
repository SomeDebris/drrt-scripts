package lib

// import (
//     "os"
//     "path/filepath"
//     "regexp"
// )

type CommandData struct {
	Flags   []string `json:"flags"`
	Faction int      `json:"faction"`
}

type Block struct {
	Id      any         `json:"ident"`
	Offset  [2]float64  `json:"offset"`
	Angle   float64     `json:"angle"`
	Command CommandData `json:"command"`
}

type ShipData struct {
	Name   string `json:"name"`
	Author string `json:"author"`
	Color0 any    `json:"color0"`
	Color1 any    `json:"color1"`
	Color2 any    `json:"color2"`
	Wgroup [4]int `json:"wgroup"`
}

type Ship struct {
	Angle    float64    `json:"angle"`
	Position [2]float64 `json:"position"`
	Data     ShipData   `json:"data"`
	Blocks   []Block    `json:"blocks"`
}

type Fleet struct {
	Blueprints []Ship `json:"blueprints"`
	Color0     any    `json:"color0"`
	Color1     any    `json:"color1"`
	Color2     any    `json:"color2"`
	Faction    int    `json:"faction"`
	Name       string `json:"name"`
}
