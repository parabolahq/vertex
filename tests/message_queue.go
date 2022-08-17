package tests

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"vertex/communication"
	"vertex/config"
)

func (s *VertexTestSuite) TestWebsocketAMQPCommunication() {
	// Connect to AMQP
	mqConn, connError := amqp.Dial(config.Config.String("amqp.url"))
	if connError != nil {
		log.Fatal(connError)
	}
	defer mqConn.Close()
	// Open channel
	channel, channelError := mqConn.Channel()
	if channelError != nil {
		log.Fatal(channelError)
	}
	defer channel.Close()
	// Declare queue
	queue, _ := channel.QueueDeclare(
		"test-service",
		false,
		false,
		false,
		false,
		nil,
	)
	// Consume messages
	messages, _ := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	testServer := httptest.NewServer(s.Gin)
	defer testServer.Close()
	conn, _, err := s.Dialer.Dial(
		fmt.Sprintf("%s/ws", strings.Replace(testServer.URL, "http", "ws", 1)),
		http.Header{
			"Authorization": {"Bearer " + s.SignedToken},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	req := communication.UserRequest{
		ServiceAlias: "test-service",
		MethodName:   "test-method",
		Params:       gin.H{},
	}
	err = conn.WriteJSON(req)
	if err != nil {
		log.Fatal(err)
	}

	for message := range messages {
		body, _ := req.ToServiceRequest("0000-0000-0000-0001")
		bodyMarshaled, _ := json.Marshal(body)
		assert.JSONEq(s.T(), string(message.Body), string(bodyMarshaled))
		channel.Close()
	}

}
