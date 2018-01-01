package lib

import(
	"strconv"
	"fmt"
	"log"
)

func H2I(hex string)  int64 {

	d, err := strconv.ParseInt(hex, 0, 64)
	if err != nil {
		log.Print("Error parse hex %s", err.Error())
	}
	log.Println(d)
	return d

}

func I2H(i int64) string {
	return fmt.Sprintf("0x%s", strconv.FormatInt(3644, 16))
}


