package main

import (
	"image"
	"image/color"
	"math"
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

var fontBase = unit.Sp(17)
var fontMedium = unit.Sp(15)
var fontSmall = unit.Sp(12)

type Title struct {
  Color color.NRGBA
  TextSize unit.Sp
  Weight font.Weight
  Alignment text.Alignment
  Text string
}

func (t Title) Layout(theme *material.Theme, gtx layout.Context) layout.Dimensions {
  title := material.Body1(theme, t.Text)
  title.Color = t.Color
  title.TextSize = t.TextSize
  title.Font.Weight = t.Weight
  title.Alignment = t.Alignment
  title.LineHeightScale = 1
  return title.Layout(gtx)
}

func MessageBar(theme *material.Theme, gtx layout.Context, messages Messages, pos int, resetPos func ()) layout.Dimensions {
  return layout.Background{}.Layout(gtx,
    func(gtx layout.Context) layout.Dimensions {
      defer clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)}.Push(gtx.Ops).Pop()
     	paint.Fill(gtx.Ops, color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF})

 			return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, gtx.Sp(fontMedium) + 20)}
    },
    func(gtx layout.Context) layout.Dimensions {
      var parDimensions layout.Dimensions

      if (len(messages) > 0) {
        var visList = layout.List{
          Axis: layout.Horizontal,
          ScrollToEnd: true,
          Position: layout.Position{
            BeforeEnd: true,
            Offset: 5 * pos,
          },
        }

        inv := op.InvalidateCmd{At: gtx.Now.Add(time.Second / 25)}
  			gtx.Execute(inv)

        var listLength = 2 * len(messages) + 1

        var lastElementsWidth int = 0
        var listWidth int = 0

  			visList.Layout(
     			gtx,
     			listLength,
  		    func(gtx layout.Context, index int) layout.Dimensions {
  					var paragraph material.LabelStyle

  					if (index % 2 != 0) {
              paragraph = material.Label(theme, fontSmall, string(messages[index/2]))
  					} else {
              paragraph = material.Label(theme, fontSmall, string("  +++  "))
  					}

            paragraph.Alignment = text.Start
            paragraph.TextSize = fontMedium
            parDimensions = paragraph.Layout(gtx)

            listWidth += parDimensions.Size.X

            if (index >= listLength - 5) {
              lastElementsWidth += parDimensions.Size.X
            }

            return parDimensions
          },
  			)

        if (5 * pos >= listWidth - lastElementsWidth) {
          resetPos()
        }
      }

      return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, parDimensions.Size.Y)}
    },
  )
}

func Line(theme *material.Theme, gtx layout.Context, departure Departure, lineHeight int) layout.FlexChild {
  var fgCol = departure.ForegroundColor
  var bgCol = departure.BackgroundColor

  return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
    layout.Stack{}.Layout(gtx,
      layout.Stacked(func(gtx layout.Context) layout.Dimensions {
  			return layout.Flex{Alignment: layout.Middle, Axis: layout.Horizontal}.Layout(gtx,
          layout.Rigid(func(gtx layout.Context) layout.Dimensions {
            width := int(1.1 * float64(lineHeight))

            layout.Background{}.Layout(gtx,
              func(gtx layout.Context) layout.Dimensions {
                defer clip.Rect{Max: image.Pt(width, lineHeight)}.Push(gtx.Ops).Pop()
                paint.Fill(gtx.Ops, color.NRGBA{R: bgCol.R, G: bgCol.G, B: bgCol.B, A: 0xFF})

           			return layout.Dimensions{Size: image.Pt(width, lineHeight)}
              },
              func(gtx layout.Context) layout.Dimensions {
                return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
                  layout.Flexed(1, func (gtx layout.Context) layout.Dimensions {
                    titleDim := Title{
                      Text:      departure.Name,
                      Color:     color.NRGBA{R: fgCol.R, G: fgCol.G, B: fgCol.B, A: 0xFF},
                      TextSize:  fontBase,
                      Weight:    font.Bold,
                      Alignment: text.Middle,
                    }.Layout(theme, gtx)

                    return layout.Dimensions{Size: image.Pt(gtx.Constraints.Min.X, titleDim.Size.Y)}
                  }),
                )
              },
            )

            return layout.Dimensions{Size: image.Pt(width, lineHeight)}
         	}),
      		layout.Rigid(layout.Spacer{Width: 15}.Layout),
          layout.Flexed(0.35, func(gtx layout.Context) layout.Dimensions {
            titleDim := Title{
              Text:      departure.Stop,
              Color:     color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF},
              TextSize:  fontSmall,
              Alignment: text.Start,
            }.Layout(theme, gtx)

            return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, titleDim.Size.Y)}
         	}),
          layout.Rigid(layout.Spacer{Width: 15}.Layout),
          layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
            titleDim := Title{
              Text:      departure.Direction,
              Color:     color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF},
              TextSize:  fontBase,
              Alignment: text.Start,
            }.Layout(theme, gtx)

            return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, titleDim.Size.Y)}
         	}),
          layout.Rigid(layout.Spacer{Width: 15}.Layout),
          layout.Flexed(0.25, func(gtx layout.Context) layout.Dimensions {
            var titleText string

            if (int(departure.dTime.Minutes()) <= 0) {
              titleText = "now"

            } else if (int(departure.dTime.Minutes()) >= 10 + departure.TimeOffset) {
              if (departure.RtTime != nil) {
                titleText = departure.RtTimeString
              } else {
                titleText = departure.TimeString
              }

            } else {
              titleText = strconv.Itoa(int(departure.dTime.Minutes()))
            }

            titleDim := Title{
              Text:      titleText,
              Color:     color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF},
              TextSize:  fontBase,
              Alignment: text.End,
            }.Layout(theme, gtx)

            return layout.Dimensions{Size: image.Pt(gtx.Constraints.Min.X, titleDim.Size.Y)}
         	}),
          layout.Rigid(layout.Spacer{Width: 5}.Layout),
          layout.Rigid(func(gtx layout.Context) layout.Dimensions {
            textContent := " "
            if (departure.RtTime != nil) {
              textContent = "*"
            }

            titleDim := Title{
              Text:      textContent,
              Color:     color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF},
              TextSize: fontSmall,
              Alignment: text.Middle,
            }.Layout(theme, gtx)

            return layout.Dimensions{Size: image.Pt(20, titleDim.Size.Y)}
         	}),
          layout.Rigid(layout.Spacer{Width: 15}.Layout),
        )
  		}),
      layout.Expanded(func(gtx layout.Context) layout.Dimensions {
  			if (departure.Cancelled) {
          return layout.Inset{ Left: 10, Right: 10 }.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
           	defer op.Offset(image.Pt(0, lineHeight / 2)).Push(gtx.Ops).Pop()
            defer clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, 4)}.Push(gtx.Ops).Pop()
            paint.ColorOp{Color: color.NRGBA{R: 0x94, G: 0x11, B: 0x00, A: 0xE5}}.Add(gtx.Ops)
            paint.PaintOp{}.Add(gtx.Ops)

            return layout.Dimensions{Size: image.Pt(gtx.Constraints.Min.X, lineHeight / 2)}
         	})
        }

        return layout.Dimensions{Size: image.Pt(0, 0)}
      }),
    )

    return layout.Dimensions{Size: image.Pt(gtx.Constraints.Min.X, lineHeight)}
 	})
}


func Display(window *app.Window, departureData chan []Departure, messageData chan Messages) error {
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

	var departures []Departure;
	var messages Messages;
	var timeString string

	messagesOffset := 0

	for {
  	select {
   	case departures = <- departureData:
      window.Invalidate()

   	case messages = <- messageData:
      // Appending first and second element to end of list for "continous scrolling"
      if (len(messages) > 0) {
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

        contentHeight := gtx.Constraints.Max.Y - 2 * (gtx.Sp(fontMedium) + 20)
        desiredLineHeight := int(3.5 * float32(gtx.Sp(fontBase)))
        numLines := math.Floor(float64(contentHeight / desiredLineHeight))
        lineHeight := math.Ceil(float64(contentHeight) / numLines)

        for i := 0; i < min(len(departures), int(numLines)); i++ {
          departureLines = append(
            departureLines,
            Line(theme, gtx, departures[i], int(lineHeight)),
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
            return layout.Flex{ Axis: layout.Vertical }.Layout(gtx,
              layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                return layout.Background{}.Layout(gtx,
                  func(gtx layout.Context) layout.Dimensions {
                    defer clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Sp(fontMedium) + 20)}.Push(gtx.Ops).Pop()
                   	paint.Fill(gtx.Ops, color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF})

               			return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, gtx.Sp(fontMedium) + 20)}
                  },
                  func(gtx layout.Context) layout.Dimensions {
                    return layout.Flex{}.Layout(gtx,
                      layout.Flexed(1, func (gtx layout.Context) layout.Dimensions {
                        return layout.Inset{Top: 10, Bottom: 10}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                          titleDimensions := Title{
                            Text:      timeString,
                            Color:     color.NRGBA{0, 0, 0, 0xFF},
                            Alignment: text.Middle,
                            TextSize:  fontMedium,
                          }.Layout(theme, gtx)

                          return layout.Dimensions{Size: titleDimensions.Size}
                        })
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
              layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                return MessageBar(theme, gtx, messages, messagesOffset, func () { messagesOffset = 0 })
             	}),
            )
          }),
        )

        // Pass the drawing operations to the GPU.
   			e.Frame(gtx.Ops)
        messagesOffset++
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
