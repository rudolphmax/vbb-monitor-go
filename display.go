package main

import (
	"fmt"
	"image"
	"image/color"
	"strconv"
	"time"

	"gioui.org/app"
	"gioui.org/font"
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
    layout.Stack{}.Layout(gtx,
      layout.Stacked(func(gtx layout.Context) layout.Dimensions {
  			return layout.Flex{Alignment: layout.Middle, Axis: layout.Horizontal}.Layout(gtx,
          layout.Rigid(func(gtx layout.Context) layout.Dimensions {
            layout.Background{}.Layout(gtx,
              func(gtx layout.Context) layout.Dimensions {
                defer clip.Rect{Max: image.Pt(size, size)}.Push(gtx.Ops).Pop()
                paint.Fill(gtx.Ops, color.NRGBA{R: bgCol.R, G: bgCol.G, B: bgCol.B, A: 0xFF})

           			return layout.Dimensions{Size: image.Pt(size, size)}
              },
              func(gtx layout.Context) layout.Dimensions {
                return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
                  layout.Flexed(1, func (gtx layout.Context) layout.Dimensions {
                    title := material.Body1(theme, departure.Name)
                    title.Color = color.NRGBA{R: fgCol.R, G: fgCol.G, B: fgCol.B, A: 0xFF}
                    title.Font.Weight = font.Bold
                    title.Alignment = text.Middle
                    title.Layout(gtx)

                    return layout.Dimensions{Size: image.Pt(gtx.Constraints.Min.X, gtx.Sp(20))}
                  }),
                )
              },
            )

            return layout.Dimensions{Size: image.Pt(size, size)}
         	}),
      		layout.Rigid(layout.Spacer{Width: 15}.Layout),
          layout.Flexed(0.35, func(gtx layout.Context) layout.Dimensions {
            title := material.Body1(theme, departure.Stop)
            title.Color = color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF}
            title.TextSize = unit.Sp(12)
            title.Alignment = text.Start
            title.Layout(gtx)

            return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, gtx.Sp(12))}
         	}),
          layout.Rigid(layout.Spacer{Width: 15}.Layout),
          layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
            title := material.Body1(theme, departure.Direction)
            title.Color = color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF}
            title.Alignment = text.Start
            title.Layout(gtx)

            return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, gtx.Sp(20))}
         	}),
          layout.Rigid(layout.Spacer{Width: 15}.Layout),
          layout.Rigid(func(gtx layout.Context) layout.Dimensions {
            if (departure.RtTime != nil) {
              var title material.LabelStyle
              title = material.Body1(theme, "*")

              title.TextSize = unit.Sp(13)
              title.Color = color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF}
              title.Alignment = text.Middle
              title.Layout(gtx)
            }

            return layout.Dimensions{Size: image.Pt(20, gtx.Sp(13))}
         	}),
          layout.Rigid(layout.Spacer{Width: 5}.Layout),
          layout.Rigid(func(gtx layout.Context) layout.Dimensions {
            var title material.LabelStyle

            if (int(departure.dTime.Minutes()) <= 0) {
              title = material.Body1(theme, "now")

            } else if (int(departure.dTime.Minutes()) >= 10 + departure.TimeOffset) {
              if (departure.RtTime != nil) {
                title = material.Body1(theme, departure.RtTimeString)
              } else {
                title = material.Body1(theme, departure.TimeString)
              }

            } else {
              title = material.Body1(theme, strconv.Itoa(int(departure.dTime.Minutes())))
            }

            title.Color = color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF}
            title.Alignment = text.Middle
            title.Layout(gtx)

            return layout.Dimensions{Size: image.Pt(size, gtx.Sp(20))}
         	}),
          layout.Rigid(layout.Spacer{Width: 15}.Layout),
        )
  		}),
      layout.Expanded(func(gtx layout.Context) layout.Dimensions {
  			if (departure.Cancelled) {
          return layout.Inset{ Left: 10, Right: 10 }.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
           	defer op.Offset(image.Pt(0, size / 2)).Push(gtx.Ops).Pop()
            defer clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, 4)}.Push(gtx.Ops).Pop()
            paint.ColorOp{Color: color.NRGBA{R: 0x94, G: 0x11, B: 0x00, A: 0xA5}}.Add(gtx.Ops)
            paint.PaintOp{}.Add(gtx.Ops)

            return layout.Dimensions{Size: image.Pt(gtx.Constraints.Min.X, size / 2)}
         	})
        }

        return layout.Dimensions{Size: image.Pt(0, 0)}
      }),
    )

    return layout.Dimensions{Size: image.Pt(gtx.Constraints.Min.X, size)}
 	})
}


func Display(window *app.Window, data chan []Departure) error {
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
		  timeChan <- time.Now().Local().Format("15:04")
			time.Sleep(10 * time.Second)
		}
	}()

  theme := material.NewTheme()

	var ops op.Ops

	var departures []Departure
	var timeString string

	for {
  	select {
   	case departures = <- data:
      window.Invalidate()

    case timeString = <- timeChan:
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
            return layout.Flex{Axis: layout.Vertical }.Layout(gtx,
              layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                return layout.Background{}.Layout(gtx,
                  func(gtx layout.Context) layout.Dimensions {
                    defer clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Sp(20))}.Push(gtx.Ops).Pop()
                   	paint.Fill(gtx.Ops, color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF})

               			return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, gtx.Sp(20))}
                  },
                  func(gtx layout.Context) layout.Dimensions {
                    return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
                      layout.Flexed(1, func (gtx layout.Context) layout.Dimensions {
                        var title material.LabelStyle
                        title = material.Body1(theme, timeString)
                        title.Color = color.NRGBA{0, 0, 0, 0xFF}
                        title.Alignment = text.Middle
                        title.Layout(gtx)

                        return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, gtx.Sp(20))}
                      }),
                    )
                  },
                )
             	}),
              layout.Flexed(1, func (gtx layout.Context) layout.Dimensions {
                return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
                  departureLines...,
                )
              }),
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
	window.Option(app.Title("vbbmon"))
	window.Option(app.Size(unit.Dp(800), unit.Dp(600)))

	return window
}

func destroyDisplay() {
  app.Main()
}
