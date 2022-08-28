package communication

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"gopkg.in/olahol/melody.v1"
	"log"
	"vertex/config"
	"vertex/utils"
)

var Connection *amqp.Connection
var Channel *amqp.Channel

type OnMessageCallback func(*melody.Melody, amqp.Delivery)

func ConnectToQueue() {
	Connection, err := amqp.Dial(config.Config.String("amqp.url"))
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	Channel, err = Connection.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	log.Println("Connected to RabbitMQ")
}

func SendMessageToService(request ServiceRequest) error {
	_, err := Channel.QueueDeclare(
		request.UserRequest.Service,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(request)
	err = Channel.PublishWithContext(
		context.Background(),
		"",
		request.UserRequest.Service,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return err
	} else {
		return nil
	}
}

func CloseEverything() {
	err := Channel.Close()
	utils.FailOnError(err, "Failed to close channel")
	err = Connection.Close()
	utils.FailOnError(err, "Failed to close connection")
}

func ListenForEvents(m *melody.Melody, callback OnMessageCallback) {
	var forever chan struct{}
	log.Println("Listening for events")
	queue, err := Channel.QueueDeclare(
		config.Config.String("amqp.queue"),
		false,
		false,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "Failed to declare a queue")
	messages, err := Channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "Failed to register a consumer")
	for d := range messages {
		log.Println(string(d.Body))
		callback(m, d)
	}
	<-forever
}
