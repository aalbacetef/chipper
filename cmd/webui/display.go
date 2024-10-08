package main

import (
	"encoding/json"
	"image"
	"image/color"
	"sync"

	"github.com/aalbacetef/chipper"
)

type Display struct {
	w       int
	h       int
	data    []byte
	palette color.Palette
	mu      sync.Mutex
}

func NewDisplay(w, h int) *Display {
	data := make([]byte, w*h)

	return &Display{
		w:    w,
		h:    h,
		data: data,
		palette: color.Palette{
			color.Black,
			color.White,
		},
	}
}

func (d *Display) stringify() string {
	if len(d.data) == 0 {
		return "[]"
	}

	type O struct {
		Data []int `json:"data"`
	}

	o := O{
		Data: make([]int, len(d.data)),
	}

	for k, dd := range d.data {
		o.Data[k] = int(dd)
	}

	data, _ := json.MarshalIndent(o, "", "  ")
	return string(data)
}

func (d *Display) String() string {
	return "<display>"
}

func (d *Display) ColorClear() color.Color {
	return d.palette[0]
}

func (d *Display) ColorSet() color.Color {
	return d.palette[1]
}

func (d *Display) Set(x, y int, c color.Color) {
	if !d.inBounds(x, y) {
		return
	}

	p := byte(chipper.ColorClear)
	if chipper.ColorEq(c, color.White) {
		p = byte(chipper.ColorSet)
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

	if p == byte(chipper.ColorSet) {
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

func (d *Display) ColorModel() color.Model {
	return d.palette
}
