package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"
)

// ApiParams represents the parameters accepted by vbbmon's api functions
type ApiParams struct {
  Base string;
  AccessId string;
  Stops []ApiStop;
  RemoveStopSuffix string;
}

// ApiStop represents a single stop as represented in the api
type ApiStop struct {
  ID string;
  Lines string;
  MaxDepartures int;
  TimeOffset int;
  Direction string;
}

// ApiColor represents a color as represented in the api
type ApiColor struct {
  R uint8;
  G uint8;
  B uint8;
}

// ApiMessage represents a HIM message as represented in the api
type ApiMessage struct {
  AffectedProduct []struct{
    Name string;
  };
  Act bool;
  Head string;
  Text string;
}

// ApiDeparture represents a departure as represented in the api
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

// Message represents the internal data model of a single (disruption) message
type Message string;

// Departure represents the internal data model of a single departure
type Departure struct {
  Name string;
  Stop string;
  Direction string;
  Cancelled bool;
  TimeString string;
  Time *time.Time;
  RtTimeString string;
  RtTime *time.Time;
  DTime time.Duration;
  TimeOffset int;
  ForegroundColor ApiColor;
  BackgroundColor ApiColor;
};

// Data represents the chunk in which data is passed through the application.
// Departures and Messages are to be nil, iff an error is present. Else, Error is to be nil.
type Data struct {
  Departures []Departure;
  Messages []Message;
  Error error;
}

// preprocessDepartures transform raw departures as returned by the api into the internal data model and sorts them
// by departure time.
// This includes calculating `dTime` (difference between departure time and current time), parsing time data,
// removing `removeStopSuffix` from stop names, and, of course, sorting.
func preprocessDepartures(data []ApiDeparture, timeOffsets []time.Duration, removeStopSuffix string) []Departure {
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
      Stop: strings.ReplaceAll(dep.Stop, removeStopSuffix, ""),
      Direction: dep.Direction,
      Cancelled: dep.Cancelled,
      TimeString: parsedTimeString,
      Time: parsedTime,
      RtTimeString: parsedRtTimeString,
      RtTime: parsedRtTime,
      DTime: dtime,
      TimeOffset: int(timeOffset.Minutes()),
      ForegroundColor: dep.ProductAtStop.Icon.ForegroundColor,
      BackgroundColor: dep.ProductAtStop.Icon.BackgroundColor,
    })
  }

  slices.SortStableFunc(departures, func(a, b Departure) int {
    var aCmpTime *time.Time
    var bCmpTime *time.Time

    if (a.RtTime != nil) {
      aCmpTime = a.RtTime
    } else {
      aCmpTime = a.Time
    }

    if (b.RtTime != nil) {
      bCmpTime = b.RtTime
    } else {
      bCmpTime = b.Time
    }

		return (*aCmpTime).Compare(*bCmpTime)
	})

  return departures
}

// FetchDepartures calls the api once for every configured stop, filters out wrong directions, and preprocesses the result.
// Returns an array of Departures and nil, iff no error occurred. Else, returns nil and the error.
func FetchDepartures(params ApiParams) ([]Departure, error) {
  var error error
  var departures []ApiDeparture
  var timeOffsets []time.Duration

  for i := range params.Stops {
    stop := params.Stops[i]

    timeOffset := time.Minute * time.Duration(stop.TimeOffset)
    offsetTimestamp := time.Now().Local().Add(timeOffset).Format("15:04")

    escapedStopId := url.PathEscape(stop.ID)
    resp, err := http.Get(params.Base + "/departureBoard/?accessId=" + params.AccessId + "&id=" + escapedStopId + "&time=" + offsetTimestamp + "&lines=" + stop.Lines + "&maxJourneys=" + strconv.Itoa(stop.MaxDepartures) + "&format=json")

    if err != nil {
      error = err
      break
    }

    body, err := io.ReadAll(resp.Body)

    if err != nil {
      error = err
      break
    }

    var res struct { Departure []ApiDeparture }
    err = json.Unmarshal([]byte(body), &res);

    if err != nil {
      error = err
      break
   	}

    // Filtering directions
    res.Departure = slices.DeleteFunc(
      res.Departure,
      func(d ApiDeparture) bool {
        return stop.Direction != "" && d.Direction != stop.Direction && d.DirectionFlag != stop.Direction
      },
    )

    for range res.Departure {
      timeOffsets = append(timeOffsets, timeOffset)
    }

    departures = append(departures, (res.Departure)...)

    resp.Body.Close()
  }

  if (error != nil) {
    return nil, error
  }

  return preprocessDepartures(departures, timeOffsets, params.RemoveStopSuffix), nil
}

// FetchMessages calls the api for HIM messages for the configured lines, filters out inactives, and preprocesses the result.
// A list of configured lines is constructed from the configured stops.
// Returns an array of Messages and nil, iff no error occurred. Else, returns nil and the error.
func FetchMessages(params ApiParams) ([]Message, error) {
  var himSearchLines string
  for i, e := range params.Stops {
    if (i == 0) {
      himSearchLines = strings.ToUpper(e.Lines)
      continue
    }

    himSearchLines = himSearchLines + "," + strings.ToUpper(e.Lines)
  }

  resp, err := http.Get(params.Base + "/himsearch/?accessId=" + params.AccessId + "&lines=" + himSearchLines + "&himcategory=1&format=json")

  if err != nil {
    return nil, err
  }

  body, err := io.ReadAll(resp.Body)

  if err != nil {
    return nil, err
  }

  var res struct { Message []ApiMessage }
  err = json.Unmarshal([]byte(body), &res);

  if err != nil {
    return nil, err
  }

  // Filtering non-active messages
  res.Message = slices.DeleteFunc(
    res.Message,
    func(d ApiMessage) bool {
      return d.Act == false
    },
  )

  var messages []Message
  for _, msg := range res.Message {
    messages = append(messages, Message(msg.Text))
  }

  return messages, nil
}
