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

  config, configError := config.Read(*configPath)

  if (configError != nil) {
    log.Fatal("Error reading config", configError)
  }

  errorData := make(chan string)
  departureData := make(chan []api.Departure)
  messageData := make(chan api.Messages)

  go func() {
    for {
      errorData <- ""

      api.FetchDepartures(
        config.Api,
        departureData,
        errorData,
      )

      time.Sleep(time.Duration(config.DepartureFetchInterval) * time.Second)
    }
  }()

  go func() {
    for {
      errorData <- ""

      api.FetchMessages(
        config.Api,
        messageData,
        errorData,
      )

      time.Sleep(time.Duration(config.MessageFetchInterval) * time.Second)
    }
  }()

  go func() {
		window := display.Init(config.Display)
		err := display.Run(window, departureData, messageData, errorData)

		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}()

	display.Destroy()
	close(messageData)
	close(departureData)
}
