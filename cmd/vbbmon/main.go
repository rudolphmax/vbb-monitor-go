package main

import (
	"flag"
	"log"
	"os"
	"time"

	"rudolphmax/vbbmon/internal/api"
	"rudolphmax/vbbmon/internal/config"
	"rudolphmax/vbbmon/internal/display"
	"rudolphmax/vbbmon/internal/utils"
)

func main() {
  configPath := flag.String("config", "./configs/config.json", "The path to the config file, e.g. ~/vbbmon/config.json")
  flag.Parse()

  config, configError := config.Read(*configPath)

  if (configError != nil) {
    log.Fatal("Error reading config", configError)
  }

  data := make(chan api.Data)

  go func() {
    var interval int
    sleepDuration := utils.Gcd(config.DepartureFetchInterval, config.MessageFetchInterval)

    var departures []api.Departure
    var messages api.Messages
    var error error

    for {
      if (interval % config.DepartureFetchInterval == 0) {
        departures, error = api.FetchDepartures(config.Api)
      }

      if (interval % config.MessageFetchInterval == 0) {
        messages, error = api.FetchMessages(config.Api)
      }

      if (interval % min(config.DepartureFetchInterval, config.MessageFetchInterval) == 0) {
        if (error != nil) {
          data <- api.Data{Error: error}

        } else {
          data <- api.Data{
            Departures: departures,
            Messages: messages,
            Error: nil,
          }
        }
      }

      if (interval == max(config.DepartureFetchInterval, config.MessageFetchInterval)) {
        interval = 0
      }

      interval += sleepDuration
      time.Sleep(time.Duration(sleepDuration) * time.Second)
    }
  }()

  go func() {
		window := display.Init(config.Display)
		err := display.Run(window, data)

		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}()

	display.Destroy()
	close(data)
}
