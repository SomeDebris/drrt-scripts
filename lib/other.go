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

func Array2DToInterface[T any](in [][]T) [][]any {
	var out [][]any = make([][]any, len(in))
	for i, row := range in {
		var outrow []any = make([]any, len(row))
		for j, val := range row {
			outrow[j] = val
		}
		out[i] = outrow
	}
	return out
}

// Returns an [][]any (or [][]any) with the same dimensions as in, but with no assigned values.
// TODO: are these values actually empty? or do I need to assign "nil" to each index?
func Array2DToEmptyInterface[T any](in [][]T) [][]any {
	var out [][]any = make([][]any, len(in))
	for i, row := range in {
		out[i] = make([]any, len(row))
	}
	return out
}

func getShipAuthorNamePairInterface(ships []rsmships.Ship) [][]any {
	var out [][]any = make([][]any, len(ships))
	for i, ship := range ships {
		var shipauthorpair []any = make([]any, 2)
		shipauthorpair[0] = ship.Data.Name
		shipauthorpair[1] = ship.Data.Author
		out[i] = shipauthorpair
	}
	return out
}
