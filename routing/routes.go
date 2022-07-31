package routing

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func SetupRoutes() (g *gin.Engine, m *melody.Melody) {
	g, m = gin.Default(), melody.New()
	g = setupPingRoute(g)
	return
}
