package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	// "gioui.org/widget/material"
)

func main() {
  data := make(chan ApiResponse)

  go func() {
    for {
      fmt.Println("Fetching data...")

      escapedStopId := url.PathEscape("")
      // escapedStopId := url.PathEscape("")

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

      data <- res

      resp.Body.Close()
      time.Sleep(10 * time.Second)
    }
  }()

  go func() {
		window := new(app.Window)
		window.Option(app.Title("Egg timer"))
		window.Option(app.Size(unit.Dp(800), unit.Dp(600)))

		err := run(window, data)

		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}()

	app.Main()

	close(data)
}

type ApiColor struct {
  R uint8;
  G uint8;
  B uint8;
}

type Departure struct {
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
  Departure []Departure;
}

func Line(theme *material.Theme, gtx layout.Context, textContent string, departure Departure) layout.FlexChild {
  var size = 100

  var fgCol = departure.ProductAtStop.Icon.ForegroundColor
  var bgCol = departure.ProductAtStop.Icon.BackgroundColor

  return layout.Rigid(func(gtx layout.Context) layout.Dimensions {

    layout.Flex{Alignment: layout.End, Axis: layout.Horizontal}.Layout(gtx,
      layout.Rigid(func(gtx layout.Context) layout.Dimensions {
        defer clip.Rect{Max: image.Pt(size, size)}.Push(gtx.Ops).Pop()
       	paint.ColorOp{Color: color.NRGBA{R: bgCol.R, G: bgCol.G, B: bgCol.B, A: 0xFF}}.Add(gtx.Ops)
       	paint.PaintOp{}.Add(gtx.Ops)

        title := material.Body1(theme, departure.Name)
        title.Color = color.NRGBA{R: fgCol.R, G: fgCol.G, B: fgCol.B, A: 0xFF}
        title.Alignment = text.Middle
        title.Layout(gtx)

        return layout.Dimensions{Size: image.Pt(size, size)}
     	}),
      layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
        title := material.Body1(theme, departure.Direction)
        title.Color = color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF}
        title.Alignment = text.Middle
        title.Layout(gtx)

        return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, size)}
     	}),
      layout.Rigid(func(gtx layout.Context) layout.Dimensions {
        parsedTime, error := time.ParseInLocation(time.DateTime, departure.Date + " " + departure.Time, time.Local)

        if (error == nil) {
          dtime := parsedTime.Sub(time.Now().Local())

          var title material.LabelStyle

          if (int(dtime.Minutes()) <= 0) {
            title = material.Body1(theme, "now")
          } else if (int(dtime.Minutes()) >= 10) {
            title = material.Body1(theme, parsedTime.Format("15:04"))
          } else {
            title = material.Body1(theme, strconv.Itoa(int(dtime.Minutes())))
          }

          title.Color = color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF}
          title.Alignment = text.Middle
          title.Layout(gtx)
        }

        return layout.Dimensions{Size: image.Pt(size, size)}
     	}),
    )

    return layout.Dimensions{Size: image.Pt(gtx.Constraints.Min.X, size)}
 	})
}


func run(window *app.Window, data chan ApiResponse) error {
 	events := make(chan event.Event)
 	acks := make(chan struct{})

  go func() {
		for {
			ev := window.Event()
			events <- ev
			<-acks
			if _, ok := ev.(app.DestroyEvent); ok {
				return
			}
		}
	}()

  theme := material.NewTheme()

	var ops op.Ops

	var departures []Departure

	for {
  	select {
   	case apiRes := <-data:
      departures = apiRes.Departure
      window.Invalidate()

		case e := <-events:
  	  switch e := e.(type) {
  		case app.DestroyEvent:
        acks <- struct{}{}
   			return e.Err

  		case app.FrameEvent:
        fmt.Println("Frame event")
   			gtx := app.NewContext(&ops, e)

        var departureLines []layout.FlexChild

        for i := range departures {
          departureLines = append(
            departureLines,
            Line(theme, gtx, "Oben", departures[i]),
          )
        }

        layout.Stack{}.Layout(gtx,
          layout.Expanded(func (gtx layout.Context) layout.Dimensions {
            defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
           	paint.ColorOp{Color: color.NRGBA{0, 0, 0, 0xFF}}.Add(gtx.Ops)
           	paint.PaintOp{}.Add(gtx.Ops)

            return layout.Dimensions{Size: gtx.Constraints.Max}
          }),
          layout.Stacked(func (gtx layout.Context) layout.Dimensions {
            return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
              departureLines...,
            )
          }),
        )

   			// Pass the drawing operations to the GPU.
   			e.Frame(gtx.Ops)
  		}

      acks <- struct{}{}
  	}
	}
}
