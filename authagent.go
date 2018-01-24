package main

import "github.com/gin-gonic/gin"

func AuthAgent(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "auth",
	})
}
