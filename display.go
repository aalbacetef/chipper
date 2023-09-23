package chipper

import (
	"fmt"
	"image"
)

type Color byte

const (
	ColorBlack Color = 0
	ColorWhite Color = 1
)

type Display struct {
	width  int
	height int
	data   []Color // storage is in row-major order.
}

func (d Display) toIndex(x, y int) int {
	return x + (y * d.width)
}

func (d *Display) Bounds() image.Rectangle {
	return image.Rect(0, 0, d.width, d.height)
}

func (d *Display) At(x, y int) Color {
	return d.data[d.toIndex(x, y)]
}

func NewDisplay(w, h int) (*Display, error) {
	if w < 0 || h < 0 {
		return nil, fmt.Errorf("width and height must be >= 0 (w=%d, h=%d)", w, h)
	}

	return &Display{
		width:  w,
		height: h,
		data:   make([]Color, w*h),
	}, nil
}

func (d *Display) Set(x, y int, c Color) error {
	point := image.Pt(x, y)
	bounds := d.Bounds()

	if !point.In(bounds) {
		return fmt.Errorf(
			"out of bounds, point(%d, %d) not in rect(%d, %d)",
			x, y, bounds.Dx(), bounds.Dy(),
		)
	}

	d.data[d.toIndex(x, y)] = c

	return nil
}
