package components

import (
	"image"
	"image/color"
	"strconv"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget/material"

	"rudolphmax/vbbmon/internal/api"
	t "rudolphmax/vbbmon/internal/display/theme"
)

type Line struct {
  Departure api.Departure;
  LineHeight int;
}

func (l Line) Layout(theme *material.Theme, gtx layout.Context) layout.FlexChild {
  var fgCol = l.Departure.ForegroundColor
  var bgCol = l.Departure.BackgroundColor

  return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
    layout.Stack{}.Layout(gtx,
      layout.Stacked(func(gtx layout.Context) layout.Dimensions {
  			return layout.Flex{Alignment: layout.Middle, Axis: layout.Horizontal}.Layout(gtx,
          layout.Rigid(func(gtx layout.Context) layout.Dimensions {
            width := int(1.1 * float64(l.LineHeight))

            layout.Background{}.Layout(gtx,
              func(gtx layout.Context) layout.Dimensions {
                defer clip.Rect{Max: image.Pt(width, l.LineHeight)}.Push(gtx.Ops).Pop()
                paint.Fill(gtx.Ops, color.NRGBA{R: bgCol.R, G: bgCol.G, B: bgCol.B, A: 0xFF})

           			return layout.Dimensions{Size: image.Pt(width, l.LineHeight)}
              },
              func(gtx layout.Context) layout.Dimensions {
                return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
                  layout.Flexed(1, func (gtx layout.Context) layout.Dimensions {
                    titleDim := Title{
                      Text:      l.Departure.Name,
                      Color:     color.NRGBA{R: fgCol.R, G: fgCol.G, B: fgCol.B, A: 0xFF},
                      TextSize:  t.FontBase,
                      Weight:    font.Bold,
                      Alignment: text.Middle,
                    }.Layout(theme, gtx)

                    return layout.Dimensions{Size: image.Pt(gtx.Constraints.Min.X, titleDim.Size.Y)}
                  }),
                )
              },
            )

            return layout.Dimensions{Size: image.Pt(width, l.LineHeight)}
         	}),
      		layout.Rigid(layout.Spacer{Width: 15}.Layout),
          layout.Flexed(0.35, func(gtx layout.Context) layout.Dimensions {
            titleDim := Title{
              Text:      l.Departure.Stop,
              Color:     color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF},
              TextSize:  t.FontSmall,
              Alignment: text.Start,
            }.Layout(theme, gtx)

            return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, titleDim.Size.Y)}
         	}),
          layout.Rigid(layout.Spacer{Width: 15}.Layout),
          layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
            titleDim := Title{
              Text:      l.Departure.Direction,
              Color:     color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF},
              TextSize:  t.FontBase,
              Alignment: text.Start,
            }.Layout(theme, gtx)

            return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, titleDim.Size.Y)}
         	}),
          layout.Rigid(layout.Spacer{Width: 15}.Layout),
          layout.Flexed(0.25, func(gtx layout.Context) layout.Dimensions {
            var titleText string

            if (int(l.Departure.DTime.Minutes()) <= 0) {
              titleText = "now"

            } else if (int(l.Departure.DTime.Minutes()) >= 10 + l.Departure.TimeOffset) {
              if (l.Departure.RtTime != nil) {
                titleText = l.Departure.RtTimeString
              } else {
                titleText = l.Departure.TimeString
              }

            } else {
              titleText = strconv.Itoa(int(l.Departure.DTime.Minutes()))
            }

            titleDim := Title{
              Text:      titleText,
              Color:     color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF},
              TextSize:  t.FontBase,
              Alignment: text.End,
            }.Layout(theme, gtx)

            return layout.Dimensions{Size: image.Pt(gtx.Constraints.Min.X, titleDim.Size.Y)}
         	}),
          layout.Rigid(layout.Spacer{Width: 5}.Layout),
          layout.Rigid(func(gtx layout.Context) layout.Dimensions {
            textContent := " "
            if (l.Departure.RtTime != nil) {
              textContent = "*"
            }

            titleDim := Title{
              Text:      textContent,
              Color:     color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF},
              TextSize:  t.FontSmall,
              Alignment: text.Middle,
            }.Layout(theme, gtx)

            return layout.Dimensions{Size: image.Pt(20, titleDim.Size.Y)}
         	}),
          layout.Rigid(layout.Spacer{Width: 15}.Layout),
        )
  		}),
      layout.Expanded(func(gtx layout.Context) layout.Dimensions {
  			if (l.Departure.Cancelled) {
          return layout.Inset{ Left: 10, Right: 10 }.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
           	defer op.Offset(image.Pt(0, l.LineHeight / 2)).Push(gtx.Ops).Pop()
            defer clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, 4)}.Push(gtx.Ops).Pop()
            paint.ColorOp{Color: color.NRGBA{R: 0x94, G: 0x11, B: 0x00, A: 0xE5}}.Add(gtx.Ops)
            paint.PaintOp{}.Add(gtx.Ops)

            return layout.Dimensions{Size: image.Pt(gtx.Constraints.Min.X, l.LineHeight / 2)}
         	})
        }

        return layout.Dimensions{Size: image.Pt(0, 0)}
      }),
    )

    return layout.Dimensions{Size: image.Pt(gtx.Constraints.Min.X, l.LineHeight)}
 	})
}
