package components

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

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
