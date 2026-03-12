package theme

import (
	"image/color"
	"rudolphmax/vbbmon/internal/utils"

	"gioui.org/unit"
)

type ThemeConfig struct {
  Font struct {
    SizeLarge int;
    SizeBase int;
    SizeMedium int;
    SizeSmall int;
  };
  GlobalForegroundColor string;
  GlobalBackgroundColor string;
};

var FontLarge unit.Sp
var FontBase unit.Sp
var FontMedium unit.Sp
var FontSmall unit.Sp

var ForegroundColor color.NRGBA;
var BackgroundColor color.NRGBA;

// Init initializes the theme with the given configuration.
// Theme aspects (like font sizes and colors) are then available as global variables.
// Colors are converted from hex strings (as in the config) to `color.NRGBA` values for direct use with Gioui.
func Init(config ThemeConfig) {
  FontLarge = unit.Sp(config.Font.SizeLarge)
  FontBase = unit.Sp(config.Font.SizeBase)
  FontMedium = unit.Sp(config.Font.SizeMedium)
  FontSmall = unit.Sp(config.Font.SizeSmall)

  ForegroundColor = color.NRGBA(utils.HexToRGBA(config.GlobalForegroundColor))
  BackgroundColor = color.NRGBA(utils.HexToRGBA(config.GlobalBackgroundColor))
}
