package utils

import (
	"image/color"
	"strconv"
)

// HexToRGBA converts a hex color string (e.g. "#FF5733") to a `color.RGBA´ struct.
func HexToRGBA(hexColor string) color.RGBA {
  values, _ := strconv.ParseUint(string(hexColor[1:]), 16, 32)

  return color.RGBA{
    R: uint8((values >> 16) & 0xFF),
    G: uint8((values >> 8) & 0xFF),
    B: uint8((values) & 0xFFb),
    A: 255,
  }
}

// Gcd calculates the greatest common divisor of two integers using the Euclidean algorithm.
func Gcd(a int, b int) int {
  for (b != 0) {
    t := b
    b = a % b
    a = t
  }

  return a
}
