package theme

import (
	"image/color"
	"rudolphmax/vbbmon/internal/utils"

	"gioui.org/unit"
)

type ThemeConfig struct {
  Font struct {
    SizeBase int;
    SizeMedium int;
    SizeSmall int;
  };
  GlobalForegroundColor string;
  GlobalBackgroundColor string;
};

var FontBase unit.Sp
var FontMedium unit.Sp
var FontSmall unit.Sp

var ForegroundColor color.NRGBA;
var BackgroundColor color.NRGBA;

func Init(config ThemeConfig) {
  FontBase = unit.Sp(config.Font.SizeBase)
  FontMedium = unit.Sp(config.Font.SizeMedium)
  FontSmall = unit.Sp(config.Font.SizeSmall)

  ForegroundColor = color.NRGBA(utils.HexToRGBA(config.GlobalForegroundColor))
  BackgroundColor = color.NRGBA(utils.HexToRGBA(config.GlobalBackgroundColor))
}
