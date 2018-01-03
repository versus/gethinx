package middle

import (
	"github.com/gin-gonic/gin"
	"log"
	"github.com/versus/gethinx/lib"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		rdr := lib.ReadRequestBody(c.Request.Body)
		c.Request.Body = rdr.Request
		log.Println("Req: ",rdr.Body) // Print request body
		c.Next()
	}
}


