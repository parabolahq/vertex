package tests

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/stretchr/testify/suite"
	"gopkg.in/olahol/melody.v1"
	"os"
	"vertex/config"
	"vertex/routing"
)

type VertexTestSuite struct {
	suite.Suite
	Gin         *gin.Engine
	Melody      *melody.Melody
	Dialer      *websocket.Dialer
	PrivateKey  *jwk.Key
	SignedToken string
}

type VertexKeysTestSuite struct {
	suite.Suite
}

func (s *VertexKeysTestSuite) SetupTest() {
	config.LoadConfigs()
}

func (s *VertexTestSuite) SetupTest() {
	s.Gin, s.Melody = routing.SetupRoutes()
	config.LoadConfigs()
	config.LoadKeys()
	if os.Getenv("TEST_JWK_PRIVATE") != "" {
		s.PrivateKey, _ = config.LoadKey(os.Getenv("TEST_JWK_PRIVATE"))
	} else {
		s.T().Skip("Private key is not specified, please, set it up via TEST_JWK_PRIVATE env")
	}
	userId := "0000-0000-0000-0001"
	payload, _ := json.Marshal(map[string]interface{}{
		"uid": userId,
	})
	tokenBytes, err := jws.Sign(payload, jwa.RS256, *s.PrivateKey)
	if err != nil {
		s.T().Fail()
	} else {
		s.SignedToken = string(tokenBytes)
		s.Dialer = websocket.DefaultDialer
	}

}
