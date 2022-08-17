package communication

import (
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

func CloseEverything() {
	err := Channel.Close()
	utils.FailOnError(err, "Failed to close channel")
	err = Connection.Close()
	utils.FailOnError(err, "Failed to close connection")
}
