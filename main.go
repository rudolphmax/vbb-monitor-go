package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
  data := make(chan []Departure)

  go func() {
    for {
      fmt.Println("Fetching data...")
      fetchData(data)
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
