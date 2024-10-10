package chipper

import "testing"

func TestDisplay(t *testing.T) { //nolint: gocognit
	t.Run("should set", func(tt *testing.T) {
		d, err := NewDebugDisplay(5, 5)
		if err != nil {
			tt.Fatalf("could not create display: %v", err)
		}

		colorSet := d.ColorSet()
		d.Set(2, 1, colorSet)

		c := d.At(2, 1)
		if !ColorEq(c, colorSet) {
			tt.Fatalf("expected %#0x, got %#0x", ColorSet, c)
		}
	})

	t.Run("simple display", func(tt *testing.T) {
		d, err := NewDebugDisplay(5, 5)
		if err != nil {
			tt.Fatalf("could not create: %v", err)
		}

		points := [][2]int{
			{1, 2},
			{1, 3},
			{2, 2},
			{2, 3},
			{3, 2},
			{3, 3},
		}

		colorSet := d.ColorSet()
		for _, p := range points {
			x := p[0]
			y := p[1]

			d.Set(x, y, colorSet)

			if !ColorEq(d.At(x, y), colorSet) {
				tt.Fatalf("point (%d, %d) was not set", x, y)
			}
		}
	})
}
