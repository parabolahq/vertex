package routing

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwt"
	"gopkg.in/olahol/melody.v1"
	"log"
	"net/http"
	"strings"
	"vertex/communication"
	"vertex/config"
	apiErr "vertex/error"
)

func setupWebSocketRoute(g *gin.Engine, m *melody.Melody) *gin.Engine {
	g.GET("/ws", func(c *gin.Context) {
		_, err := AuthorizeRequest(c.Request)
		if err != nil {
			c.JSON(403, apiErr.ApiError{
				ErrorCode: 1,
				Data:      "Authentication error",
			}.AsInternalEvent())
		} else {
			err = m.HandleRequest(c.Writer, c.Request)
			if err != nil {
				c.JSON(500, apiErr.ApiError{
					ErrorCode: 0,
					Data:      "Internal error occurred",
				}.AsInternalEvent())
			}
		}
	})
	return g
}

func AuthorizeRequest(req *http.Request) (data map[string]interface{}, err error) {
	authorizationHeader := strings.TrimSpace(req.Header.Get("Authorization"))
	log.Println(authorizationHeader)
	if authorizationHeader == "" || !strings.HasPrefix(authorizationHeader, "Bearer") {
		return nil, errors.New("token unspecified")
	}
	tokenRaw := strings.TrimSpace(strings.TrimPrefix(authorizationHeader, "Bearer"))
	token, err := jwt.Parse(
		[]byte(tokenRaw),
		jwt.WithKeySet(config.KeySet),
		jwt.WithValidate(true),
	)
	if err != nil {
		return nil, err
	}
	data, _ = token.AsMap(context.Background())
	_, uidExists := data["uid"]
	if !uidExists {
		return nil, errors.New("uid not found in token")
	}
	return data, nil
}

func HandleMessage(s *melody.Session, data []byte) {
	userRequest := communication.UserRequest{}
	jsonErr := json.Unmarshal(data, &userRequest)
	if jsonErr != nil {
		validationErr := apiErr.ApiError{
			ErrorCode: 2,
			Data:      "Parse of JSON failed",
		}.AsInternalEvent()
		s.Write(validationErr.AsBytes())
	} else {
		if userRequest.ServiceAlias == "vertex" {
			switch userRequest.MethodName {
			case "disconnect":
				{
					closeErr := s.Close()
					if closeErr != nil {
						log.Print(closeErr)
					} else {
						if s.IsClosed() {
							log.Printf("Session %s succesfuly closed", s.Keys["session_id"])
						}
					}
				}
			}
		}
	}
}

func HandleConnection(s *melody.Session) {
	sessionId, _ := uuid.NewUUID()
	s.Keys = map[string]interface{}{}
	s.Keys["session_id"] = sessionId.String()
	log.Printf("Connected user. Added new session %s", sessionId.String())
}
