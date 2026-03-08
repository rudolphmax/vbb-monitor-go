package components

import (
	"image"
	"image/color"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget/material"

	"rudolphmax/vbbmon/internal/api"
	t "rudolphmax/vbbmon/internal/display/theme"
)

func MessageBar(theme *material.Theme, gtx layout.Context, messages api.Messages, pos int, resetPos func ()) layout.Dimensions {
  return layout.Background{}.Layout(gtx,
    func(gtx layout.Context) layout.Dimensions {
      defer clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y)}.Push(gtx.Ops).Pop()
     	paint.Fill(gtx.Ops, color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF})

 			return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, gtx.Sp(t.FontMedium) + 20)}
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
              paragraph = material.Label(theme, t.FontSmall, string(messages[index/2]))
  					} else {
              paragraph = material.Label(theme, t.FontSmall, string("  +++  "))
  					}

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

        if (5 * pos >= listWidth - lastElementsWidth) {
          resetPos()
        }
      }

      return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, parDimensions.Size.Y)}
    },
  )
}
