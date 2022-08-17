package main

import (
	"vertex/communication"
	"vertex/config"
	"vertex/routing"
	"vertex/utils"
)

func main() {
	config.LoadConfigs()
	config.LoadKeys()
	communication.ConnectToQueue()
	// Closing RabbitMQ session on exit
	defer communication.CloseEverything()
	g, _ := routing.SetupRoutes()
	err := g.Run(config.Config.String("bindaddr"))
	utils.FailOnError(err, "Failed to start gin server")
}
