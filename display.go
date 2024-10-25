package chipper

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strings"
)

type Display interface {
	fmt.Stringer
	draw.Image
	ColorClear() color.Color
	ColorSet() color.Color
}

func Each(d Display, fn func(int, int) error) error {
	b := d.Bounds()

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if err := fn(x, y); err != nil {
				return err
			}
		}
	}

	return nil
}

const (
	ColorClear = iota
	ColorSet
)

type DebugDisplay struct {
	width   int
	height  int
	data    []bool // storage is in row-major order.
	palette color.Palette
}

func (d *DebugDisplay) ColorClear() color.Color {
	return color.Black
}

func (d *DebugDisplay) ColorSet() color.Color {
	return color.White
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
			if ColorEq(d.At(x, y), d.ColorSet()) {
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

func (d *DebugDisplay) At(x, y int) color.Color {
	if d.data[d.toIndex(x, y)] {
		return d.ColorSet()
	}

	return d.ColorClear()
}

func (d *DebugDisplay) ColorModel() color.Model {
	return d.palette
}

func NewDebugDisplay(w, h int) (*DebugDisplay, error) {
	if w < 0 || h < 0 {
		return nil, fmt.Errorf("width and height must be >= 0 (w=%d, h=%d)", w, h)
	}

	return &DebugDisplay{
		width:  w,
		height: h,
		data:   make([]bool, w*h),
		palette: []color.Color{
			color.Black,
			color.White,
		},
	}, nil
}

func (d *DebugDisplay) Set(x, y int, c color.Color) {
	point := image.Pt(x, y)
	bounds := d.Bounds()

	if !point.In(bounds) {
		panic(fmt.Sprintf(
			"out of bounds, point(%d, %d) not in rect(%d, %d)",
			x, y, bounds.Dx(), bounds.Dy(),
		))
	}

	is := ColorEq(c, d.ColorSet())
	d.data[d.toIndex(x, y)] = is
}

func ColorEq(c1, c2 color.Color) bool {
	//nolint: varnamelen
	r1, g1, b1, a1 := c1.RGBA()

	//nolint: varnamelen
	r2, g2, b2, a2 := c2.RGBA()

	if r1 != r2 {
		return false
	}

	if g1 != g2 {
		return false
	}

	if b1 != b2 {
		return false
	}

	if a1 != a2 {
		return false
	}

	return true
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
