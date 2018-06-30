package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(headerRequired())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message_from_original": "hello"})
	})

	r.GET("/heavy", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message_from_original": "hello hello"})
	})

	http.Handle("/", r)
}

func headerRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: avoid direct access
		miningResult := c.GetHeader("X-MINING-GATEWAY-KEY")
		fmt.Println(miningResult)
		if miningResult == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid access"})
		}

		c.Next()
	}
}
