package lib

import (
	"bytes"
	"io"
)

func ReadBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	s := buf.String()
	return s
}