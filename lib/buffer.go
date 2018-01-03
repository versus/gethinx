package lib

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
)

type MyRequestBody struct {
	Body string
	Request io.ReadCloser
}

func ReadRequestBody(reader io.Reader) (MyRequestBody)  {
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal("Error parse c.Request.Body  %s", err.Error())
	}
	buf2 := new(bytes.Buffer)
	buf2.ReadFrom(ioutil.NopCloser(bytes.NewBuffer(buf)))
	return MyRequestBody{buf2.String(),ioutil.NopCloser(bytes.NewBuffer(buf))}
}

func TrimQuote(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}