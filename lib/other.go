package lib

import (
	"fmt"
	"github.com/SomeDebris/rsmships-go"
	"strconv"
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

func Array2DToInt(in [][]any) ([][]int, error) {
	out := make([][]int, len(in))
	var ok bool
	for i, row := range in {
		outrow := make([]int, len(row))
		for j, val := range row {
			outrow[j], ok = val.(int)
			if !ok {
				return out, fmt.Errorf("Failed to convert [][]any to [][]int at index out[%d][%d].", i, j)
			}
		}
		out[i] = outrow
	}
	return out, nil
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

func Int2dSliceToString(ints [][]int) ([][]string, error) {
	count := len(ints)
	records := make([][]string, count)
	for j, row := range ints {
		record := make([]string, len(row))
		for i, val := range row {
			record[i] = strconv.Itoa(val)
		}
		records[j] = record
	}
	return records, nil
}

