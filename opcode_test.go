package chipper

import (
	"testing"
)

func TestOpcode(t *testing.T) {
	b := []byte{0x1a, 0xe1}
	want := 0x1ae1
	got := toUint16(b)

	if got != uint16(want) {
		t.Fatalf("got %#0x, want %#0x", got, want)
	}

	t.Run("test overflow", func(tt *testing.T) {
		vx := byte(0xFF)
		tt.Logf("vx: %#0x", vx)
		vv := byte(1) + vx
		tt.Logf("added: %#0x", vv)
	})
}
