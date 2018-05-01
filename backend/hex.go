package backend

import (
	"fmt"
	"strconv"
)

// H2I function return int64 number from hex string of number
func H2I(hex string) (int64, error) {
	d, err := strconv.ParseInt(hex, 0, 64)
	return d, err
}

// I2H function return hex string of int64 base number
func I2H(i int64) string {
	return fmt.Sprintf("0x%s", strconv.FormatInt(i, 16))
}
