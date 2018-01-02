package lib

import(
	"strconv"
	"fmt"
)

func H2I(hex string)  (int64, error) {
	d, err := strconv.ParseInt(hex, 0, 64)
	return d, err
}

func I2H(i int64) string {
	return fmt.Sprintf("0x%s", strconv.FormatInt(3644, 16))
}


