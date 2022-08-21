package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
	"vertex/communication"
	"vertex/config"
)

func (s *VertexTestSuite) TestWebsocketAMQPCommunication() {
	channel, chanErr := s.AMQPConnection.Channel()
	if chanErr != nil {
		s.Fail("Can't connect message channel")
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
		_, err := channel.QueueDelete(queue.Name, false, false, false)
		if err != nil {
			s.T().Fail()
		}
		channel.Close()
	}
}

func (s *VertexTestSuite) TestWebsocketAMQPMessageReceive() {
	channel, chanErr := s.AMQPConnection.Channel()
	if chanErr != nil {
		s.Fail("Can't connect message channel")
	}
	defer channel.Close()
	testServer := httptest.NewServer(s.Gin)
	defer testServer.Close()
	conn, _, err := s.Dialer.Dial(
		fmt.Sprintf("%s/ws", strings.Replace(testServer.URL, "http", "ws", 1)),
		http.Header{
			"Authorization": {"Bearer " + s.SignedToken},
		},
	)
	time.Sleep(time.Microsecond * 100)
	if err != nil {
		log.Fatal(err)
	}
	MessageBody, _ := json.Marshal(communication.Event{
		ServiceAlias: "testService",
		EventType:    "testEventType",
		RecipientIds: []string{"0000-0000-0000-0001"},
		Data:         &map[string]interface{}{},
	})
	var event *communication.Event
	channel.PublishWithContext(
		context.Background(),
		"",
		config.Config.String("amqp.queue"),
		false,
		false,
		amqp.Publishing{
			Body: MessageBody,
		},
	)
	s.T().Log("Message published")
	_, message, _ := conn.ReadMessage()
	json.Unmarshal(message, &event)
	if err != nil {
		s.Assert().Fail(err.Error())
	}
	assert.Equal(s.T(), &communication.Event{
		ServiceAlias: "testService",
		EventType:    "testEventType",
		RecipientIds: nil,
		Data:         &map[string]interface{}{},
	}, event)
	conn.Close()
}
