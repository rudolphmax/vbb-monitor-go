package components

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget/material"

	t "rudolphmax/vbbmon/internal/display/theme"
)

// ClockBar is a component that displays the current time in a screen-wide bar.
type ClockBar struct {
  TimeString string;
}

func (cb ClockBar) Layout(theme *material.Theme, gtx layout.Context) layout.Dimensions {
  return layout.Background{}.Layout(gtx,
    func(gtx layout.Context) layout.Dimensions {
      defer clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Sp(t.FontMedium) + 20)}.Push(gtx.Ops).Pop()
     	paint.Fill(gtx.Ops, t.ForegroundColor)

 			return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, gtx.Sp(t.FontMedium) + 20)}
    },
    func(gtx layout.Context) layout.Dimensions {
      return layout.Flex{}.Layout(gtx,
        layout.Flexed(1, func (gtx layout.Context) layout.Dimensions {
          return layout.Inset{Top: 10, Bottom: 10}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
            titleDimensions := Title{
              Text:      cb.TimeString,
              Color:     &t.BackgroundColor,
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
