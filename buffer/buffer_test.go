package buffer

import (
	"io"
	"log"
	"os"
	"testing"
)

type eReader struct {
	r io.Reader
}

func TestReadRequestBody(t *testing.T) {
	f, err := os.Open("buffer.go")
	if err != nil {
		log.Fatal(err)
	}
	r := &eReader{r: f}
	req := ReadRequestBody(r.r)
	if !(len(req.Body) > 0) {
		t.Error("Error ReadRequestBody")
	}

}
