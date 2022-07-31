package tests

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gopkg.in/olahol/melody.v1"
	"vertex/config"
	"vertex/routing"
)

type VertexTestSuite struct {
	suite.Suite
	Gin    *gin.Engine
	Melody *melody.Melody
}

func (s *VertexTestSuite) SetupTest() {
	s.Gin, s.Melody = routing.SetupRoutes()
	config.LoadConfigs()
}
