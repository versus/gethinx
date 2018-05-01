package middle

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/versus/gethinx/buffer"
	"github.com/versus/gethinx/monitoring"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		rdr := buffer.ReadRequestBody(c.Request.Body)
		c.Request.Body = rdr.Request
		monitoring.PromRequest.Inc()
		log.Println("Req: ", rdr.Body) // Print request body
		c.Next()
	}
}
