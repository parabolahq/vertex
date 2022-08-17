package routing

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
	"vertex/communication"
)

func SetupRoutes() (g *gin.Engine, m *melody.Melody) {
	g, m = gin.Default(), melody.New()
	g = setupPingRoute(g)
	g = setupWebSocketRoute(g, m)
	m = setupMelodyEvents(m)
	return
}

func setupMelodyEvents(m *melody.Melody) *melody.Melody {
	m.HandleMessage(HandleMessage)
	m.HandleConnect(HandleConnection)
	m.HandleDisconnect(HandleDisconnection)
	go communication.ListenForEvents(m, HandleAmqpMessage)
	return m
}
