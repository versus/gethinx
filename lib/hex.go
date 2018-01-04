package lib

import (
	"fmt"
	"strconv"
)

//_ = lib.H2I("0xe6")
//_ = lib.I2H(230)

func H2I(hex string) (int64, error) {
	d, err := strconv.ParseInt(hex, 0, 64)
	return d, err
}

func I2H(i int64) string {
	return fmt.Sprintf("0x%s", strconv.FormatInt(3644, 16))
}
