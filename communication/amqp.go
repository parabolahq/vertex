package communication

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"vertex/config"
	"vertex/utils"
)

var Connection *amqp.Connection
var Channel *amqp.Channel

func ConnectToQueue() {
	Connection, err := amqp.Dial(config.Config.String("amqp.url"))
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	Channel, err = Connection.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	log.Println("Connected to RabbitMQ")
}

func SendMessageToService(request ServiceRequest) error {
	_, err := Channel.QueueDeclare(
		request.UserRequest.ServiceAlias,
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
		request.UserRequest.ServiceAlias,
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
