package main

import (
	"fmt"
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

  data := make(chan []Departure)

  go func() {
    for {
      fmt.Println("Fetching data...")

      fetchData(
        config.Api,
        data,
      )

      time.Sleep(10 * time.Second)
    }
  }()

  go func() {
		window := initDisplay()
		err := Display(window, data)

		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}()

	destroyDisplay()
	close(data)
}
