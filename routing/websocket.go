package routing

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopkg.in/olahol/melody.v1"
	"log"
	"vertex/communication"
	error2 "vertex/error"
)

func setupWebSocketRoute(g *gin.Engine, m *melody.Melody) *gin.Engine {
	g.GET("/ws", func(c *gin.Context) {
		err := m.HandleRequest(c.Writer, c.Request)
		if err != nil {
			c.JSON(500, error2.ApiError{
				ErrorCode: 1,
				Data:      "Internal error occurred",
			}.AsInternalEvent())
		}
	})
	return g
}

func HandleMessage(s *melody.Session, data []byte) {
	userRequest := communication.UserRequest{}
	json.Unmarshal(data, &userRequest)
	if userRequest.ServiceAlias == "vertex" && userRequest.MethodName == "disconnect" {
		log.Println("Received disconnect signal. Closing session")
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

func HandleConnection(s *melody.Session) {
	sessionId, _ := uuid.NewUUID()
	s.Keys = map[string]interface{}{}
	s.Keys["session_id"] = sessionId.String()
	log.Printf("Connected user. Added new session %s", sessionId.String())
}
