package chipper

import "testing"

func TestDisplay(t *testing.T) {
	t.Run("should set", func(tt *testing.T) {
		display, err := NewDebugDisplay(5, 5)
		if err != nil {
			tt.Fatalf("could not create display: %v", err)
		}

		colorSet := display.ColorSet()
		display.Set(2, 1, colorSet)

		c := display.At(2, 1)
		if !ColorEq(c, colorSet) {
			tt.Fatalf("expected %#0x, got %#0x", ColorSet, c)
		}
	})

	t.Run("simple display", func(tt *testing.T) {
		display, err := NewDebugDisplay(5, 5)
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

		colorSet := display.ColorSet()

		for _, p := range points {
			x := p[0]
			y := p[1]

			display.Set(x, y, colorSet)

			if !ColorEq(display.At(x, y), colorSet) {
				tt.Fatalf("point (%d, %d) was not set", x, y)
			}
		}
	})
}
