package gethinx

import (
	"log"
	"testing"
)

func TestKey(t *testing.T) {
	rnd := Key(32)
	log.Println("key = ", rnd)
	if len(rnd) != 64 {
		t.Error("Error len rnd", len(rnd))
	}

}
