package tests

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/stretchr/testify/suite"
	"gopkg.in/olahol/melody.v1"
	"os"
	"vertex/communication"
	"vertex/config"
	"vertex/routing"
)

type VertexTestSuite struct {
	suite.Suite
	Gin         *gin.Engine
	Melody      *melody.Melody
	Dialer      *websocket.Dialer
	SignedToken string
}

type VertexKeysTestSuite struct {
	suite.Suite
}

func (s *VertexKeysTestSuite) SetupTest() {
	config.LoadConfigs()
}

func CreateEmptyToken() (string, error) {
	var privateKey *jwk.Key
	if os.Getenv("TEST_JWK_PRIVATE") != "" {
		privateKey, _ = config.LoadKey(os.Getenv("TEST_JWK_PRIVATE"))
	} else {
		return "", errors.New("private key is not specified, please, set it up via TEST_JWK_PRIVATE env")
	}
	userId := "0000-0000-0000-0001"
	payload, _ := json.Marshal(map[string]interface{}{
		"uid": userId,
	})
	tokenBytes, err := jws.Sign(payload, jwa.RS256, *privateKey)
	if err != nil {
		return "", err
	} else {
		return string(tokenBytes), nil
	}
}

func (s *VertexTestSuite) SetupTest() {
	s.Gin, s.Melody = routing.SetupRoutes()
	config.LoadConfigs()
	config.LoadKeys()
	communication.ConnectToQueue()
	emptyToken, err := CreateEmptyToken()
	if err != nil {
		s.T().Fatal(err)
	} else {
		s.SignedToken = emptyToken
		s.Dialer = websocket.DefaultDialer
	}

}
