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
}
