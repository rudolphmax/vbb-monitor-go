package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type ApiColor struct {
  R uint8;
  G uint8;
  B uint8;
}

type ApiDeparture struct {
  Name string;
  Direction string;
  Cancelled bool;
  Date string;
  Time string;
  RtTime string;
  RtDate string;
  ProductAtStop struct {
    Icon struct {
      ForegroundColor ApiColor;
      BackgroundColor ApiColor;
    }
  }
};

type ApiResponse struct {
  Departure []ApiDeparture;
}

type Departure struct {
  Name string;
  Direction string;
  Cancelled bool;
  Date string;
  Time string;
  RtTime string;
  RtDate string;
  dTime time.Duration;
  ForegroundColor ApiColor;
  BackgroundColor ApiColor;
};

func preprocess(data ApiResponse) []Departure {
  var departures []Departure

  for i := range data.Departure {
    dep := data.Departure[i]

    var dtime time.Duration = 0
    var parsedTime time.Time
    var parsedRtTime time.Time

    parsedTime, error := time.ParseInLocation(time.DateTime, dep.Date + " " + dep.Time, time.Local)

    if (dep.RtTime != "") {
      parsedRtTime, error := time.ParseInLocation(time.DateTime, dep.RtDate + " " + dep.RtTime, time.Local)

      if (error == nil) {
        dtime = parsedRtTime.Sub(time.Now().Local())
      }

    } else if (error == nil) {
      dtime = parsedTime.Sub(time.Now().Local())
    }

    departures = append(departures, Departure{
      Name: dep.Name,
      Direction: dep.Direction,
      Cancelled: dep.Cancelled,
      Time: parsedTime.Format("15:04"),
      RtTime: parsedRtTime.Format("15:04"),
      dTime: dtime,
      ForegroundColor: dep.ProductAtStop.Icon.ForegroundColor,
      BackgroundColor: dep.ProductAtStop.Icon.BackgroundColor,
    })
  }

  return departures
}

func fetchData(data chan []Departure) {
  escapedStopId := url.PathEscape("")

  resp, err := http.Get("" + escapedStopId + "")

  if err != nil {
    panic(err)
  }

  body, err := io.ReadAll(resp.Body)

  var res ApiResponse

  err = json.Unmarshal([]byte(body), &res);
 	if err != nil {
		fmt.Println("error:", err)
 	}

  data <- preprocess(res)

  resp.Body.Close()
}
