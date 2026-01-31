package main

import (
	"fmt"
	"flag"
	"log"
	"os"
	"time"
)

func main() {
  apiBase := flag.String("base", "", "The base url of the API, e.g.: hhtps://example.com/api/info/v2")
  apiAccessId := flag.String("accessId", "", "The HAFAS Access-ID to use with the API.")
  apiStopId := flag.String("stop", "", "The ID of the stop to monitor.")

  flag.Parse()

  if (apiBase == nil || *apiBase == "") {
    log.Fatal("No API base provided, use -base to provide one.")
  }
  if (apiAccessId == nil || *apiAccessId == "") {
    log.Fatal("No AccessID provided, use -accessId to provide one.")
  }
  if (apiStopId == nil || *apiStopId == "") {
    log.Fatal("No StopID provided, use -stop to provide one.")
  }

  data := make(chan []Departure)

  go func() {
    for {
      fmt.Println("Fetching data...")

      fetchData(
        ApiParams{
          Base: *apiBase,
          AccessId: *apiAccessId,
          StopID: *apiStopId,
        },
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
