package tests

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gopkg.in/olahol/melody.v1"
	"tiny-submarine/config"
	"tiny-submarine/routing"
)

type SubmarineTestSuite struct {
	suite.Suite
	Gin    *gin.Engine
	Melody *melody.Melody
}

func (s *SubmarineTestSuite) SetupTest() {
	s.Gin, s.Melody = routing.SetupRoutes()
	config.LoadConfigs()
}
