package main

import (
	"log"
	"tiny-submarine/config"
	"tiny-submarine/routing"
)

func main() {
	config.LoadConfigs()
	g, _ := routing.SetupRoutes()
	err := g.Run(config.Config.String("bindaddr"))
	if err != nil {
		log.Fatal(err)
	}
}
