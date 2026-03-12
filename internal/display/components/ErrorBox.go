package components

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/widget/material"

	t "rudolphmax/vbbmon/internal/display/theme"
)

// ErrorBox is a component that displays an error message.
type ErrorBox struct {
  Error string
}

func (eb ErrorBox) Layout(theme *material.Theme, gtx layout.Context) layout.Dimensions {
  return layout.UniformInset(10).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
    dimensions := Title{
      Text: "Error: " + eb.Error,
      TextSize:  t.FontMedium,
      Color: &color.NRGBA{R: 0xF8, G: 0x43, B: 0xD, A: 0xFF},
    }.Layout(theme, gtx)

    return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, dimensions.Size.Y)}
  })
}
