package display

import (
	"math"
	"time"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"rudolphmax/vbbmon/internal/api"
	"rudolphmax/vbbmon/internal/display/components"
	t "rudolphmax/vbbmon/internal/display/theme"
)

type DisplayConfig struct {
  Theme t.ThemeConfig;
  NumLines int;
  ScrollSpeed float32;
}

var displayConfig DisplayConfig

func Run(window *app.Window, data chan api.Data) error {
 	events := make(chan event.Event)
 	acks := make(chan struct{})
  timeChan := make(chan string)

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

  go func() {
		for {
		  timeChan <- time.Now().Local().Format("15:04:05")
			time.Sleep(500 * time.Millisecond)
		}
	}()

  theme := material.NewTheme()

	var ops op.Ops

	var error error;
	var departures []api.Departure;
	var messages api.Messages;
	var timeString string

	var ClockBarHeight int
  var MessageBarHeight int

	messagesOffset := 0

	for {
  	select {
    case d := <- data:
      error = d.Error
      departures = d.Departures
      messages = d.Messages

      // Appending first and second element to end of list for "continous scrolling"
      if (len(d.Messages) > 0) {
        messages = append(messages, messages[0])
        messages = append(messages, messages[min(1, len(messages))])
      }

      window.Invalidate()

    case timeString = <- timeChan:
      window.Invalidate()

		case e := <-events:
  	  switch e := e.(type) {
  		case app.DestroyEvent:
        acks <- struct{}{}
   			return e.Err

  		case app.FrameEvent:
   			gtx := app.NewContext(&ops, e)

        var departureLines []layout.FlexChild

        contentHeight := gtx.Constraints.Max.Y - (ClockBarHeight + MessageBarHeight)
        lineHeight := math.Ceil(float64(contentHeight) / float64(displayConfig.NumLines))

        for i := 0; i < min(len(departures), int(displayConfig.NumLines)); i++ {
          departureLines = append(
            departureLines,
            components.Line{
              Departure: departures[i],
              LineHeight: int(lineHeight),
            }.Layout(theme, gtx),
          )
        }

        layout.Stack{}.Layout(gtx,
          layout.Expanded(func (gtx layout.Context) layout.Dimensions {
            defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
           	paint.ColorOp{Color: t.BackgroundColor}.Add(gtx.Ops)
           	paint.PaintOp{}.Add(gtx.Ops)

            return layout.Dimensions{Size: gtx.Constraints.Max}
          }),
          layout.Stacked(func (gtx layout.Context) layout.Dimensions {
            return layout.Flex{ Axis: layout.Vertical }.Layout(gtx,
              layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                dimensions := components.ClockBar{
                  TimeString: timeString,
                }.Layout(theme, gtx)

                ClockBarHeight = dimensions.Size.Y

                return dimensions
             	}),
              layout.Flexed(1, func (gtx layout.Context) layout.Dimensions {
                if (error == nil) {
                  return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
                    departureLines...,
                  )
                } else {
                  return components.ErrorBox{
                    Error: error.Error(),
                  }.Layout(theme, gtx)
                }
              }),
              layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                var dimensions layout.Dimensions

                if (len(messages) > 0) {
                  dimensions = components.MessageBar{
                    Messages: messages,
                    Pos: messagesOffset,
                    ResetPos: func () { messagesOffset = 0 },
                    Speed: displayConfig.ScrollSpeed,
                  }.Layout(theme, gtx)

                  MessageBarHeight = dimensions.Size.Y
                }

                return dimensions
             	}),
            )
          }),
        )

        inv := op.InvalidateCmd{At: gtx.Now.Add(time.Second / 25)}
  			gtx.Execute(inv)

        // Pass the drawing operations to the GPU.
   			e.Frame(gtx.Ops)
        messagesOffset++
  		}

      acks <- struct{}{}
  	}
	}
}

func Init(config DisplayConfig) *app.Window {
	displayConfig = config

  window := new(app.Window)
	window.Option(app.Title("vbbmon"))
	window.Option(app.Size(unit.Dp(800), unit.Dp(600)))

	t.Init(config.Theme)

	return window
}

func Destroy() {
  app.Main()
}
