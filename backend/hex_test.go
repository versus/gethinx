package backend

import (
	"log"
	"testing"
)

func TestH2I(t *testing.T) {
	num, err := H2I("0xe6")
	if err != nil {
		t.Error("Error H2I", err.Error())
	}
	if num != 230 {
		t.Error("Wrong convert H2I")
	}
}

func TestI2H(t *testing.T) {
	hex := I2H(230)
	if hex != "0xe6" {
		t.Error("Error I2H")
	}
}
