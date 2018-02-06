package middle

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/versus/gethinx/lib"
	"github.com/versus/gethinx/monitoring"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		rdr := lib.ReadRequestBody(c.Request.Body)
		c.Request.Body = rdr.Request
		monitoring.PromRequest.Inc()
		log.Println("Req: ", rdr.Body) // Print request body
		c.Next()
	}
}
