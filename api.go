package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"time"
)

type ApiParams struct {
  Base string;
  AccessId string;
  StopID string;
}

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
  TimeString string;
  Time *time.Time;
  RtTimeString string;
  RtTime *time.Time;
  dTime time.Duration;
  ForegroundColor ApiColor;
  BackgroundColor ApiColor;
};

func preprocess(data ApiResponse) []Departure {
  var departures []Departure

  for i := range data.Departure {
    dep := data.Departure[i]

    var dtime time.Duration = 0
    var parsedTime *time.Time = nil
    var parsedTimeString string = ""
    var parsedRtTime *time.Time = nil
    var parsedRtTimeString string = ""

    var error error

    pt, error := time.ParseInLocation(time.DateTime, dep.Date + " " + dep.Time, time.Local)

    if (error == nil) {
      parsedTime = &pt
      parsedTimeString = (*parsedTime).Format("15:04")
    }

    if (dep.RtTime != "") {
      prt, error := time.ParseInLocation(time.DateTime, dep.RtDate + " " + dep.RtTime, time.Local)

      if (error == nil) {
        parsedRtTime = &prt
        parsedRtTimeString = (*parsedRtTime).Format("15:04")
        dtime = (*parsedRtTime).Sub(time.Now().Local())
      }

    } else if (error == nil) {
      dtime = (*parsedTime).Sub(time.Now().Local())
    }

    departures = append(departures, Departure{
      Name: dep.Name,
      Direction: dep.Direction,
      Cancelled: dep.Cancelled,
      TimeString: parsedTimeString,
      Time: parsedTime,
      RtTimeString: parsedRtTimeString,
      RtTime: parsedRtTime,
      dTime: dtime,
      ForegroundColor: dep.ProductAtStop.Icon.ForegroundColor,
      BackgroundColor: dep.ProductAtStop.Icon.BackgroundColor,
    })
  }

  slices.SortStableFunc(departures, func(a, b Departure) int {
		return (*a.Time).Compare(*b.Time)
	})

  return departures
}

func fetchData(params ApiParams, data chan []Departure) {
  escapedStopId := url.PathEscape(params.StopID)

  resp, err := http.Get(params.Base + "?accessId=" + params.AccessId + "&id=" + escapedStopId + "&format=json")

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
