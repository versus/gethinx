package main

import (
	"github.com/gin-gonic/gin"
	"sync/atomic"
)

func getStatus(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": atomic.LoadInt64(&numBlocks),
	})
}
