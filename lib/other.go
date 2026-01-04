package lib

import (
	"fmt"
	"strconv"
	"github.com/SomeDebris/rsmships-go"
)

func HTMLColorCodeToInt(code string) (int, error) {
	if code[0] != '#' {
		return 0, fmt.Errorf("String '%s' does not start with a '#', but it needs to.", code)
	}

	out, err := strconv.Atoi(fmt.Sprintf("0x%s", code[1:]))
	if err != nil {
		return 0, err
	}

	return out, nil
}

func Array2DToInterface[T any](in [][]T) [][]interface{} {
	var out [][]interface{} = make([][]interface{}, len(in))
	for i, row := range in {
		var outrow []interface{} = make([]interface{}, len(row))
		for j, val := range row {
			outrow[j] = val
		}
		out[i] = outrow
	}
	return out
}

// Returns an [][]any (or [][]interface{}) with the same dimensions as in, but with no assigned values.
// TODO: are these values actually empty? or do I need to assign "nil" to each index?
func Array2DToEmptyInterface[T any](in [][]T) [][]interface{} {
	var out [][]interface{} = make([][]interface{}, len(in))
	for i, row := range in {
		out[i] = make([]interface{}, len(row))
	}
	return out
}

func getShipAuthorNamePairInterface(ships []rsmships.Ship) [][]interface{} {
	var out [][]interface{} = make([][]interface{}, len(ships))
	for i, ship := range ships {
		var shipauthorpair []interface{} = make([]interface{}, 2)
		shipauthorpair[0] = ship.Data.Name
		shipauthorpair[1] = ship.Data.Author
		out[i] = shipauthorpair
	}
	return out
}
