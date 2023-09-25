package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"gitlab.com/aalbacetef/chipper"
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

	for {
		err := emu.Tick()
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("END")
				return
			}
			fmt.Println("ERROR: ", err)
			continue
		}

		time.Sleep(delay)
		chipper.DumpEmu(emu)
	}
}
