package main

import (
	"flag"
	"log"
	"os"
	"time"
)

func main() {
  configPath := flag.String("config", "./config.json", "The path to the config file, e.g. ~/vbbmon/config.json")
  flag.Parse()

  config, error := readConfig(*configPath)

  if (error != nil) {
    log.Fatal("Error reading config", error)
  }

  departureData := make(chan []Departure)
  messageData := make(chan Messages)

  go func() {
    for {
      fetchDepartures(
        config.Api,
        departureData,
      )

      time.Sleep(time.Duration(config.DepartureFetchInterval) * time.Second)
    }
  }()

  go func() {
    for {
      fetchMessages(
        config.Api,
        messageData,
      )

      time.Sleep(time.Duration(config.MessageFetchInterval) * time.Second)
    }
  }()

  go func() {
		window := initDisplay()
		err := Display(window, departureData, messageData)

		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}()

	destroyDisplay()
	close(messageData)
	close(departureData)
}
