package convcode

import (
	"fmt"
	"strconv"
)

func Hex2Dec(hexStr string) (int, error) {
	n, err := strconv.ParseUint(hexStr, 16, 64)
	if err != nil {
		return 0, err
	}
	return int(n), nil
}

func Dec2Hex(dec int) string {
	return fmt.Sprintf("%X", dec)
}
