package lib

import (
	"fmt"
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
