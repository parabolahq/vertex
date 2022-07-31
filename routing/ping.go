package routing

import "github.com/gin-gonic/gin"

func setupPingRoute(g *gin.Engine) *gin.Engine {
	g.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"ok": true,
			"data": gin.H{
				"pong": "pong",
			},
		})
	})
	return g
}
