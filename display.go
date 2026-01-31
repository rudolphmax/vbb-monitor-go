package main

import (
	"fmt"
	"image"
	"image/color"
	"strconv"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func Line(theme *material.Theme, gtx layout.Context, textContent string, departure Departure) layout.FlexChild {
  var size = 100

  var fgCol = departure.ForegroundColor
  var bgCol = departure.BackgroundColor

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
        var title material.LabelStyle

        if (int(departure.dTime.Minutes()) <= 0) {
          title = material.Body1(theme, "now")
        } else if (int(departure.dTime.Minutes()) >= 10) {
          title = material.Body1(theme, departure.Time)
        } else {
          title = material.Body1(theme, strconv.Itoa(int(departure.dTime.Minutes())))
        }

        title.Color = color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF}
        title.Alignment = text.Middle
        title.Layout(gtx)

        return layout.Dimensions{Size: image.Pt(size, size)}
     	}),
    )

    return layout.Dimensions{Size: image.Pt(gtx.Constraints.Min.X, size)}
 	})
}


func Display(window *app.Window, data chan []Departure) error {
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
   	case departures = <-data:
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

func initDisplay() *app.Window {
  window := new(app.Window)
	window.Option(app.Title("Egg timer"))
	window.Option(app.Size(unit.Dp(800), unit.Dp(600)))

	return window
}

func destroyDisplay() {
  app.Main()
}
