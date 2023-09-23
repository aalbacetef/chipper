package chipper

import (
	"bytes"
	_ "embed"
	"testing"
)

//go:embed testdata/maze-alt-david-winter-199x.ch8
var testRom []byte

func TestEmulator(t *testing.T) {
	const (
		testStackSize = 16
		testRamSize   = 4096
		testW         = 64
		testH         = 64
	)

	emu, err := NewEmulator(testStackSize, testRamSize, testW, testH)
	if err != nil {
		t.Fatalf("could not create emulator: %v", err)
	}

	r := bytes.NewReader(testRom)
	if err := emu.Load(r); err != nil {
		t.Fatalf("could not load rom: %v", err)
	}

	for {
		tickErr := emu.Tick()
		if tickErr != nil {
			t.Fatalf("error from Tick: %v", tickErr)
		}
	}
}
