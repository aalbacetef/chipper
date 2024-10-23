package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aalbacetef/chipper"
)

const (
	ramSize = 4096
	w       = 64
	h       = 32
)

func main() {
	fname := ""
	delayms := 500
	stackSize := 16

	flag.StringVar(&fname, "name", fname, "name of rom (path)")
	flag.IntVar(&delayms, "delay", delayms, "delay in ms")
	flag.IntVar(&stackSize, "stack", stackSize, "stack size")

	flag.Parse()

	if fname == "" {
		flag.Usage()

		return
	}

	delay := time.Duration(delayms) * time.Millisecond

	data, err := os.ReadFile(fname)
	if err != nil {
		fmt.Println("could not open file: ", err)

		return
	}

	r := bytes.NewReader(data)

	emu, err := mkEmu(stackSize, ramSize, w, h)
	if err != nil {
		fmt.Println(err)

		return
	}

	if err := runUntilError(r, emu, delay); err != nil {
		fmt.Println("error: ", err)
	}
}

func mkEmu(stackSize, ramSize, w, h int) (*chipper.Emulator, error) {
	display, err := chipper.NewDebugDisplay(w, h)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	emu, err := chipper.NewEmulator(
		stackSize,
		ramSize,
		display,
		&chipper.StubKeyInputSource{},
	)
	if err != nil {
		return nil, fmt.Errorf("error creating emulator: %w", err)
	}

	return emu, nil
}

func runUntilError(r io.Reader, emu *chipper.Emulator, delay time.Duration) error {
	if err := emu.Load(r); err != nil {
		return fmt.Errorf("could not load ROM: %w", err)
	}

	for {
		err := emu.Tick()
		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return fmt.Errorf("runUntilError: %w", err)
		}

		time.Sleep(delay)
		chipper.DumpEmu(emu)
	}
}
