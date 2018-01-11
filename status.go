package main

import (
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

func getStatus(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": atomic.LoadInt64(&lastBlock),
	})
}
