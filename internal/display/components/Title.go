package components

import (
	"image/color"
	t "rudolphmax/vbbmon/internal/display/theme"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// Title is a component that displays a title with customizable color, text size, weight and alignment.
type Title struct {
  Color *color.NRGBA
  TextSize unit.Sp
  Weight font.Weight
  Alignment text.Alignment
  Text string
}

func (title Title) Layout(theme *material.Theme, gtx layout.Context) layout.Dimensions {
  element := material.Body1(theme, title.Text)

  if title.Color != nil {
    element.Color = *title.Color
  } else {
    element.Color = t.ForegroundColor
  }

  element.TextSize = title.TextSize
  element.Font.Weight = title.Weight
  element.Alignment = title.Alignment
  element.LineHeightScale = 1

  return element.Layout(gtx)
}
