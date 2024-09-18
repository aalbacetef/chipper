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

	emu, err := chipper.NewEmulator(
		stackSize,
		ramSize,
		w, h,
	)
	if err != nil {
		fmt.Println("error creating emulator: ", err)
	}

	if err := emu.Load(r); err != nil {
		fmt.Println("could not load ROM: ", err)

		return
	}

	if err := runUntilError(emu, delay); err != nil {
		fmt.Println("error: ", err)
	}
}

func runUntilError(emu *chipper.Emulator, delay time.Duration) error {
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
