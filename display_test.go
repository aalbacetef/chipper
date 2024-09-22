package chipper

import "testing"

func TestDisplay(t *testing.T) { //nolint: gocognit
	t.Run("should set", func(tt *testing.T) {
		d, err := NewDebugDisplay(5, 5)
		if err != nil {
			tt.Fatalf("could not create display: %v", err)
		}

		if err := d.Set(2, 1, ColorWhite); err != nil {
			tt.Fatalf("could not set (%d, %d) to %#0x", 2, 1, ColorWhite)
		}

		c := d.At(2, 1)
		if c != ColorWhite {
			tt.Fatalf("expected %#0x, got %#0x", ColorWhite, c)
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

		for _, p := range points {
			x := p[0]
			y := p[1]

			if err := d.Set(x, y, 1); err != nil {
				tt.Fatalf("could not set (%d, %d): %v", x, y, err)
			}

			if v := d.At(x, y); v != 1 {
				tt.Fatalf("point (%d, %d) was not set", x, y)
			}
		}
	})
}
