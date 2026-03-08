package components

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget/material"

	t "rudolphmax/vbbmon/internal/display/theme"
)

type ClockBar struct {
  TimeString string;
}

func (cb ClockBar) Layout(theme *material.Theme, gtx layout.Context) layout.Dimensions {
  return layout.Background{}.Layout(gtx,
    func(gtx layout.Context) layout.Dimensions {
      defer clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Sp(t.FontMedium) + 20)}.Push(gtx.Ops).Pop()
     	paint.Fill(gtx.Ops, color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF})

 			return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, gtx.Sp(t.FontMedium) + 20)}
    },
    func(gtx layout.Context) layout.Dimensions {
      return layout.Flex{}.Layout(gtx,
        layout.Flexed(1, func (gtx layout.Context) layout.Dimensions {
          return layout.Inset{Top: 10, Bottom: 10}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
            titleDimensions := Title{
              Text:      cb.TimeString,
              Color:     color.NRGBA{0, 0, 0, 0xFF},
              Alignment: text.Middle,
              TextSize:  t.FontMedium,
            }.Layout(theme, gtx)

            return layout.Dimensions{Size: titleDimensions.Size}
          })
        }),
      )
    },
  )
}
