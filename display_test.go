package chipper

import "testing"

func TestDisplay(t *testing.T) {
	t.Run("should set", func(tt *testing.T) {

		d, err := NewDisplay(5, 5)
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
		d, err := NewDisplay(5, 5)
		if err != nil {
			tt.Fatalf("could not create: %v", err)
		}

		d.Set(1, 2, 1)
		d.Set(1, 3, 1)
		d.Set(2, 2, 1)
		d.Set(2, 3, 1)
		d.Set(3, 2, 1)
		d.Set(3, 3, 1)
		t.Logf("")
		t.Logf("\n%s", d.String())
	})
}
