package chipper

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"testing"
)

//go:embed testdata/zero-demo-2007.ch8
var testZero []byte

//go:embed testdata/zero.golden.txt
var testZeroGolden []byte

//go:embed testdata/particle-demo-zero-2008.ch8
var testParticle []byte

//go:embed testdata/particle.golden.txt
var testParticleGolden []byte

//go:embed testdata/maze-alt-david-winter-199x.ch8
var testMaze []byte

//go:embed testdata/maze.golden.txt
var testMazeGolden []byte

func TestEmulator(t *testing.T) {
	cases := []struct {
		label  string
		rom    []byte
		golden []byte
	}{
		{"zero", testZero, testZeroGolden},
		{"particle", testParticle, testParticleGolden},
		{"maze", testMaze, testMazeGolden},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			testROM(t, c.rom, c.golden)
		})
	}
}

func testROM(t *testing.T, rom []byte, golden []byte) {
	t.Helper()

	const (
		testStackSize = 16
		testRAMSize   = 4096
		testW         = 64
		testH         = 32
		testMaxCount  = 200
	)

	emu, err := NewEmulator(testStackSize, testRAMSize, testW, testH)
	if err != nil {
		t.Fatalf("could not create emulator: %v", err)
	}

	r := bytes.NewReader(rom)
	if err := emu.Load(r); err != nil {
		t.Fatalf("could not load rom: %v", err)
	}

	b := &strings.Builder{}

	for count := 0; count < testMaxCount; count++ {
		tickErr := emu.Tick()
		if tickErr != nil {
			t.Logf("error from Tick: %v", tickErr)

			break
		}

		for k, v := range emu.V {
			fmt.Fprintf(
				b,
				"  (v%2d) %#0x\n",
				k, v,
			)
		}
	}

	data := []byte(b.String())
	n := len(data)
	m := len(golden)

	if m != n {
		t.Fatalf("expected size %d, got %d", m, n)
	}

	for k := 0; k < n; k++ {
		want := golden[k]
		got := data[k]

		if want != got {
			t.Fatalf("mismatch at %d, want %#0X, got %#0X", k, want, got)
		}
	}
}
