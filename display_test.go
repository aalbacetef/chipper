package chipper

import "testing"

func TestDisplay(t *testing.T) {
	d, err := NewDisplay(5, 5)
	if err != nil {
		t.Fatalf("could not create display: %v", err)
	}

	if err := d.Set(2, 1, ColorWhite); err != nil {
		t.Fatalf("could not set (%d, %d) to %#0x", 2, 1, ColorWhite)
	}

	c := d.At(2, 1)
	if c != ColorWhite {
		t.Fatalf("expected %#0x, got %#0x", ColorWhite, c)
	}
}
