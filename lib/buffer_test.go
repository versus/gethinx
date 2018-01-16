package lib

import (
	"io"
	"log"
	"os"
	"testing"
)

func TestTrimQuote(t *testing.T) {
	str := TrimQuote("\"test\"")
	if str != "test" {
		t.Error("Error TestTrimQuote ", str)
	}
}

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
