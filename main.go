package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type RecordCallRequest struct {
	Uid string `json:"uid"`
}

type RateLimitValues struct {
	limits map[string]int
}

const MaxRate = 100

func main() {
	r := gin.Default()

	rateLimit := &RateLimitValues{
		limits: map[string]int{
			"example": 1,
		},
	}

	go replenishRateLimits(rateLimit)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/limits", func(c *gin.Context) {

		uid := c.Query("uid")
		limit, exists := rateLimit.limits[uid]

		if !exists {
			c.AbortWithStatus(http.StatusNotFound)
		}

		c.JSON(http.StatusOK, gin.H{
			"available": limit,
		})
	})

	r.POST("/limits/record-call", func(c *gin.Context) {
		var request RecordCallRequest

		c.BindJSON(&request)

		limit, exists := rateLimit.limits[request.Uid]

		if !exists {
			c.AbortWithStatus(http.StatusNotFound)
		}

		if limit < 1 {
			c.AbortWithStatus(http.StatusTooManyRequests)
		}

		rateLimit.limits[request.Uid]--

		c.String(http.StatusNoContent, "")
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func replenishRateLimits(limits *RateLimitValues) {
	for ; ; {
		for uid, limit := range limits.limits {
			if limit < MaxRate {
				limits.limits[uid]++
			}
		}
		time.Sleep(1 * time.Second)
	}
}
