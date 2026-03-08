package main

import (
	"flag"
	"log"
	"os"
	"time"

	"rudolphmax/vbbmon/internal/api"
	"rudolphmax/vbbmon/internal/config"
	"rudolphmax/vbbmon/internal/display"
)

func main() {
  configPath := flag.String("config", "./configs/config.json", "The path to the config file, e.g. ~/vbbmon/config.json")
  flag.Parse()

  config, error := config.Read(*configPath)

  if (error != nil) {
    log.Fatal("Error reading config", error)
  }

  departureData := make(chan []api.Departure)
  messageData := make(chan api.Messages)

  go func() {
    for {
      api.FetchDepartures(
        config.Api,
        departureData,
      )

      time.Sleep(time.Duration(config.DepartureFetchInterval) * time.Second)
    }
  }()

  go func() {
    for {
      api.FetchMessages(
        config.Api,
        messageData,
      )

      time.Sleep(time.Duration(config.MessageFetchInterval) * time.Second)
    }
  }()

  go func() {
		window := display.Init(config.Display)
		err := display.Run(window, departureData, messageData)

		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}()

	display.Destroy()
	close(messageData)
	close(departureData)
}
