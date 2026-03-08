package theme

import (
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

var ForegroundColor string;
var BackgroundColor string;

func Init(config ThemeConfig) {
  FontBase = unit.Sp(config.Font.SizeBase)
  FontMedium = unit.Sp(config.Font.SizeMedium)
  FontSmall = unit.Sp(config.Font.SizeSmall)

  ForegroundColor = config.GlobalForegroundColor
  BackgroundColor = config.GlobalBackgroundColor
}
