package buffer

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
)

type MyRequestBody struct {
	Body    string
	Request io.ReadCloser
}

func ReadRequestBody(reader io.Reader) MyRequestBody {
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal("Error parse c.Request.Body  ", err.Error())
	}
	buf2 := new(bytes.Buffer)
	buf2.ReadFrom(ioutil.NopCloser(bytes.NewBuffer(buf)))
	return MyRequestBody{buf2.String(), ioutil.NopCloser(bytes.NewBuffer(buf))}
}
