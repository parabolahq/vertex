package tests

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/suite"
	"gopkg.in/olahol/melody.v1"
	"vertex/config"
	"vertex/routing"
)

type VertexTestSuite struct {
	suite.Suite
	Gin    *gin.Engine
	Melody *melody.Melody
	Dialer *websocket.Dialer
}

func (s *VertexTestSuite) SetupTest() {
	s.Gin, s.Melody = routing.SetupRoutes()
	config.LoadConfigs()
	s.Dialer = websocket.DefaultDialer
}
