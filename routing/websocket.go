package routing

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwt"
	amqp "github.com/rabbitmq/amqp091-go"
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

func TokenFromRequest(req *http.Request) (string, error) {
	var tokenRaw string
	if req.URL.Query().Get("token") != "" {
		tokenRaw = req.URL.Query().Get("token")
	} else {
		authorizationHeader := strings.TrimSpace(req.Header.Get("Authorization"))
		if authorizationHeader == "" || !strings.HasPrefix(authorizationHeader, "Bearer") {
			return "", errors.New("token unspecified")
		}
		tokenRaw = strings.TrimSpace(strings.TrimPrefix(authorizationHeader, "Bearer"))
	}
	return tokenRaw, nil
}

func AuthorizeRequest(req *http.Request) (data map[string]interface{}, err error) {
	var tokenRaw string
	tokenRaw, err = TokenFromRequest(req)
	if err != nil {
		return nil, err
	}
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
		if userRequest.Service == "vertex" {
			switch userRequest.Method {
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
						Service: "vertex",
						Event:   "infoUpdate",
						Data:    &map[string]interface{}{"userId": s.Keys["user_id"]},
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
	communication.StoreSession(UserId, s)
	log.Printf("Connected user. Added new session for user %s", data["uid"])
}

func HandleDisconnection(s *melody.Session) {
	UserId := s.Keys["user_id"].(string)
	for _, handlerName := range config.Config.Strings("handlers.disconnect") {
		// If message is not delivered to handler, and error occurs, then just ignore it
		_ = communication.SendMessageToService(communication.PoolActionServiceRequest(handlerName, UserId, "disconnect"))
	}
	communication.RemoveSession(UserId, s)
	log.Printf("Disconnected user. Removed session for user %s", s.Keys["user_id"])
}

func HandleAmqpMessage(m *melody.Melody, d amqp.Delivery) {
	receivedEvent := new(communication.Event)
	err := json.Unmarshal(d.Body, &receivedEvent)
	log.Printf("Recieved message from %s\n", receivedEvent.Service)
	if err != nil {
		log.Println("Failed to decode message:", err)
	} else {
		recipients := receivedEvent.Recipients
		receivedEvent.Recipients = nil
		broadcastSessions := []*melody.Session{}
		for _, recipient := range recipients {
			sessions := communication.GetSessions(recipient)
			broadcastSessions = append(broadcastSessions, sessions...)
		}
		if len(broadcastSessions) > 0 {
			eventSerialized := receivedEvent.AsBytes()
			err := m.BroadcastMultiple(eventSerialized, broadcastSessions)
			if err != nil {
				log.Println("Failed to broadcast message:", err)
			}
			log.Printf("Message broadcasted to %d recepients", len(broadcastSessions))
		}
	}
}
