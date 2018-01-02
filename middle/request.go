package middle

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"bytes"
	"log"
	"github.com/versus/gethinx/lib"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, _ := ioutil.ReadAll(c.Request.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.
		log.Println("Req: ",lib.ReadBody(rdr1)) // Print request body
		c.Request.Body = rdr2
		c.Next()
	}
}


