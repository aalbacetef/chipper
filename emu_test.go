package chipper

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

//go:embed testdata/zero-demo-2007.ch8
var testZero []byte

//go:embed testdata/particle-demo-zero-2008.ch8
var testParticle []byte

//go:embed testdata/maze-alt-david-winter-199x.ch8
var testMaze []byte

func TestEmulator(t *testing.T) {
	const (
		testStackSize = 16
		testRamSize   = 4096
		testW         = 64
		testH         = 32
		testDelay     = 500 * time.Millisecond
	)

	testRom := testMaze
	if readFile := os.Getenv("EMU_ROM_PATH"); readFile != "" {
		data, err := os.ReadFile(readFile)
		if err != nil {
			t.Fatalf("could not read file(%s): %v", readFile, err)
		}

		testRom = data
	}

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
			t.Logf("error from Tick: %v", tickErr)
		}

		if testDelay > 0 {
			time.Sleep(testDelay)
		}

		fmt.Printf("PC: %0#4x\n", emu.PC)
		fmt.Println("Stack: ", emu.Stack.String())
		fmt.Println("->")
		b := &strings.Builder{}
		for k, v := range emu.V {
			fmt.Fprintf(
				b,
				"  (v%2d) %#0x\n",
				k, v,
			)
		}
		fmt.Println(b.String())
		fmt.Println(emu.Display.String())
		fmt.Println("---------------")

	}
}
