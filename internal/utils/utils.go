package utils

import (
	"image/color"
	"strconv"
)

func HexToRGBA(hexColor string) color.RGBA {
  values, _ := strconv.ParseUint(string(hexColor[1:]), 16, 32)

  return color.RGBA{
    R: uint8((values >> 16) & 0xFF),
    G: uint8((values >> 8) & 0xFF),
    B: uint8((values) & 0xFFb),
    A: 255,
  }
}
