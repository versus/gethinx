package middle

import (
	"bytes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/versus/gethinx/monitoring"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// ResponseLogger logger for responce from geth
func ResponseLogger(c *gin.Context) {
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw
	monitoring.PromResponse.Inc()
	c.Next()
	statusCode := c.Writer.Status()
	log.Println("Response: ", c.ClientIP(), statusCode, blw.body.String())
}
