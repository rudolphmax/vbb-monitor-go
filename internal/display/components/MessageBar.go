package components

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget/material"

	"rudolphmax/vbbmon/internal/api"
	t "rudolphmax/vbbmon/internal/display/theme"
)

// MessageBar is a component that displays a list of messages that infinitely scrolls horizontsally like a newsticker.
type MessageBar struct {
  Messages []api.Message;
  Pos int;
  ResetPos func ()
  Speed float32;
}

func (mb MessageBar) Layout(theme *material.Theme, gtx layout.Context) layout.Dimensions {
  return layout.Background{}.Layout(gtx,
    func(gtx layout.Context) layout.Dimensions {
      defer clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)}.Push(gtx.Ops).Pop()
     	paint.Fill(gtx.Ops, t.ForegroundColor)

 			return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, gtx.Sp(t.FontMedium) + 20)}
    },
    func(gtx layout.Context) layout.Dimensions {
      var parDimensions layout.Dimensions

      if (len(mb.Messages) > 0) {
        var visList = layout.List{
          Axis: layout.Horizontal,
          ScrollToEnd: true,
          Position: layout.Position{
            BeforeEnd: true,
            Offset: int(mb.Speed * float32(mb.Pos)),
          },
        }

        var listLength = 2 * len(mb.Messages) + 1

        var lastElementsWidth int = 0
        var listWidth int = 0

  			visList.Layout(
     			gtx,
     			listLength,
  		    func(gtx layout.Context, index int) layout.Dimensions {
  					var paragraph material.LabelStyle

  					if (index % 2 != 0) {
              paragraph = material.Label(theme, t.FontSmall, string(mb.Messages[index/2]))
  					} else {
              paragraph = material.Label(theme, t.FontSmall, string("  +++  "))
  					}

            paragraph.Color = t.BackgroundColor
            paragraph.Alignment = text.Start
            paragraph.TextSize = t.FontMedium
            parDimensions = paragraph.Layout(gtx)

            listWidth += parDimensions.Size.X

            if (index >= listLength - 5) {
              lastElementsWidth += parDimensions.Size.X
            }

            return parDimensions
          },
  			)

        if (int(mb.Speed) * mb.Pos >= listWidth - lastElementsWidth) {
          mb.ResetPos()
        }
      }

      return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, parDimensions.Size.Y)}
    },
  )
}
