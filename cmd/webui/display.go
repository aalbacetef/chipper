package main

import (
	"image"
	"image/color"

	"github.com/aalbacetef/chipper"
)

type Display struct {
	w    int
	h    int
	data []byte
}

func NewDisplay(w, h int) *Display {
	data := make([]byte, w*h)

	return &Display{
		w:    w,
		h:    h,
		data: data,
	}
}

func (d *Display) Set(x, y int, c color.Color) {
	if !d.inBounds(x, y) {
		return
	}

	p := byte(chipper.ColorBlack)
	if colorEq(c, color.White) {
		p = byte(chipper.ColorWhite)
	}

	idx := d.toIndex(x, y)
	d.data[idx] = p
}

func (d *Display) inBounds(x, y int) bool {
	r := d.Bounds()
	p := image.Point{x, y}

	return p.In(r)
}

func (d *Display) At(x, y int) color.Color {
	if !d.inBounds(x, y) {
		return color.Black
	}

	idx := d.toIndex(x, y)
	p := d.data[idx]

	if p == byte(chipper.ColorWhite) {
		return color.White
	}

	return color.Black
}

func (d *Display) Bounds() image.Rectangle {
	return image.Rect(0, 0, d.w, d.h)
}

func (d *Display) toIndex(x, y int) int {
	return x + (y * d.w)
}

func colorEq(a, b color.Color) bool {
	aR, aG, aB, aA := a.RGBA()
	bR, bG, bB, bA := b.RGBA()

	if aR != bR {
		return false
	}

	if aG != bG {
		return false
	}

	if aB != bB {
		return false
	}

	if aA != bA {
		return false
	}

	return true
}
