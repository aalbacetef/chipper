package chipper

import (
	"fmt"
	"image"
	"image/color"
	"strings"
)

// type Display interface {
// 	draw.Image
// }

type Display interface {
	fmt.Stringer
	Bounds() image.Rectangle
	Set(x, y int, c Color) error
	At(x, y int) Color
}

func Each(d Display, fn func(int, int) error) error {
	b := d.Bounds()

	for x := b.Min.X; x < b.Max.X; x++ {
		for y := b.Min.Y; y < b.Max.Y; y++ {
			if err := fn(x, y); err != nil {
				return err
			}
		}
	}

	return nil
}

const (
	ColorClear Color = iota
	ColorSet
)

type Color byte

func (c Color) String() string {
	switch c {
	case ColorClear:
		return "ColorBlack"
	case ColorSet:
		return "ColorWhite"
	default:
		return fmt.Sprintf("%0#x", c)
	}
}

func (c Color) RGBA() (uint32, uint32, uint32, uint32) {
	switch c {
	case ColorSet:
		return color.White.RGBA()
	default:
		return color.Black.RGBA()
	}
}

type DebugDisplay struct {
	width  int
	height int
	data   []Color // storage is in row-major order.
}

func (d *DebugDisplay) String() string {
	b := &strings.Builder{}
	rows := d.height
	cols := d.width

	top := &strings.Builder{}
	top.WriteString("     ")

	for k := 0; k < cols*2; k++ {
		fmt.Fprintf(top, "-")
	}

	b.WriteString(top.String())
	fmt.Fprintf(b, "\n")

	for y := 0; y < rows; y++ {
		fmt.Fprintf(b, " %2d |", y)

		for x := 0; x < cols; x++ {
			c := " ."
			if d.At(x, y) == ColorSet {
				c = " o"
			}

			fmt.Fprintf(b, "%s", c)
		}

		fmt.Fprintf(b, "|\n")
	}

	b.WriteString(top.String())

	return b.String()
}

func (d DebugDisplay) toIndex(x, y int) int {
	return x + (y * d.width)
}

func (d *DebugDisplay) Bounds() image.Rectangle {
	return image.Rect(0, 0, d.width, d.height)
}

func (d *DebugDisplay) At(x, y int) Color {
	return d.data[d.toIndex(x, y)]
}

func NewDebugDisplay(w, h int) (*DebugDisplay, error) {
	if w < 0 || h < 0 {
		return nil, fmt.Errorf("width and height must be >= 0 (w=%d, h=%d)", w, h)
	}

	return &DebugDisplay{
		width:  w,
		height: h,
		data:   make([]Color, w*h),
	}, nil
}

func (d *DebugDisplay) Set(x, y int, c Color) error {
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

func loadSprites(emu *Emulator) error { //nolint: unparam
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

	copy(emu.RAM, data)

	return nil
}
