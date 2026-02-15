package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"time"
)

type Stop struct {
  ID string;
  Lines string;
  MaxDepartures int;
  TimeOffset int;
  Direction string;
}

type ApiParams struct {
  Base string;
  AccessId string;
  Stops []Stop;
}

type ApiColor struct {
  R uint8;
  G uint8;
  B uint8;
}

type ApiDeparture struct {
  Name string;
  Stop string;
  Direction string;
  DirectionFlag string;
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
  Stop string;
  Direction string;
  Cancelled bool;
  TimeString string;
  Time *time.Time;
  RtTimeString string;
  RtTime *time.Time;
  dTime time.Duration;
  TimeOffset int;
  ForegroundColor ApiColor;
  BackgroundColor ApiColor;
};

func preprocess(data []ApiDeparture, timeOffsets []time.Duration) []Departure {
  var departures []Departure

  for i := range data {
    dep := data[i]
    timeOffset := timeOffsets[i]

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
      Stop: dep.Stop,
      Direction: dep.Direction,
      Cancelled: dep.Cancelled,
      TimeString: parsedTimeString,
      Time: parsedTime,
      RtTimeString: parsedRtTimeString,
      RtTime: parsedRtTime,
      dTime: dtime,
      TimeOffset: int(timeOffset.Minutes()),
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
  var departures []ApiDeparture
  var timeOffsets []time.Duration

  for i := range params.Stops {
    stop := params.Stops[i]

    timeOffset := time.Minute * time.Duration(stop.TimeOffset)
    offsetTimestamp := time.Now().Local().Add(timeOffset).Format("15:04")

    escapedStopId := url.PathEscape(stop.ID)
    resp, err := http.Get(params.Base + "?accessId=" + params.AccessId + "&id=" + escapedStopId + "&time=" + offsetTimestamp + "&lines=" + stop.Lines + "&maxJourneys=" + strconv.Itoa(stop.MaxDepartures) + "&format=json")

    if err != nil {
      panic(err)
    }

    body, err := io.ReadAll(resp.Body)

    var res ApiResponse

    err = json.Unmarshal([]byte(body), &res);
   	if err != nil {
  		fmt.Println("error:", err)
   	}

    // Filtering directions
    res.Departure = slices.DeleteFunc(
      res.Departure,
      func(d ApiDeparture) bool {
        return stop.Direction != "" && d.Direction != stop.Direction && d.DirectionFlag != stop.Direction
      },
    )

    for _ = range res.Departure {
      timeOffsets = append(timeOffsets, timeOffset)
    }

    departures = append(departures, (res.Departure)...)

    resp.Body.Close()
  }

  data <- preprocess(departures, timeOffsets)
}
