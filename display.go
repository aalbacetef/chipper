package chipper

import (
	"fmt"
	"image"
	"strings"
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

func (d *Display) String() string {
	b := &strings.Builder{}
	rows := d.height

	for k := 0; k < rows; k++ {
		fmt.Fprintf(b, "  ")
		for j := 0; j < d.width; j++ {
			c := "."
			if int(d.At(j, k)) != 0 {
				c = "+"
			}

			fmt.Fprintf(b, "%s", c)
		}
		fmt.Fprintf(b, "  \n")
	}

	return b.String()
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

func loadSprites(emu *Emulator) error {
	data := []byte{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}

	copy(data, emu.RAM)

	return nil
}
