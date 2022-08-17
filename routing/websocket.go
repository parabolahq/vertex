package routing

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
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
			c.JSON(403, apiErr.InvalidToken().AsInternalEvent())
		} else {
			err = m.HandleRequest(c.Writer, c.Request)
			if err != nil {
				c.JSON(500, apiErr.InternalError().AsInternalEvent())
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
		validationErr := apiErr.BadRequest().AsInternalEvent()
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
							log.Printf("Session %s succesfuly closed", s.Keys["user_id"])
						}
					}
				}
			case "userInfo":
				{
					asJson, _ := json.Marshal(communication.Event{
						ServiceAlias: "vertex",
						EventType:    "infoUpdate",
						Data:         &map[string]interface{}{"userId": s.Keys["user_id"]},
					})
					s.Write(asJson)
				}
			}
		} else {
			serviceRequest, err := userRequest.ToServiceRequest(s.Keys["user_id"].(string))
			if err != nil {
				errMsg := apiErr.BadRequest().AsInternalEvent()
				s.Write(errMsg.AsBytes())
			} else {
				err := communication.SendMessageToService(serviceRequest)
				if err != nil {
					errMsg := apiErr.SendingMQError().AsInternalEvent()
					s.Write(errMsg.AsBytes())
				}
			}
		}
	}
}

func HandleConnection(s *melody.Session) {
	data, _ := AuthorizeRequest(s.Request)
	s.Keys = map[string]interface{}{}
	UserId := data["uid"].(string)
	s.Keys["user_id"] = UserId
	for _, handlerName := range config.Config.Strings("handlers.connect") {
		// If message is not delivered to handler, and error occurs, then just ignore it
		_ = communication.SendMessageToService(communication.PoolActionServiceRequest(handlerName, UserId, "connect"))
	}
	log.Printf("Connected user. Added new session for user %s", data["uid"])
}

func HandleDisconnection(s *melody.Session) {
	UserId := s.Keys["user_id"].(string)
	for _, handlerName := range config.Config.Strings("handlers.disconnect") {
		// If message is not delivered to handler, and error occurs, then just ignore it
		_ = communication.SendMessageToService(communication.PoolActionServiceRequest(handlerName, UserId, "disconnect"))
	}
	log.Printf("Disconnected user. Removed session for user %s", s.Keys["user_id"])
}
