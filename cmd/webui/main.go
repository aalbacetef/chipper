package main

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"syscall/js"
	"time"

	"github.com/aalbacetef/chipper"
)

var (
	emu *chipper.Emulator
	d   chipper.Display
)

//go:embed zero-demo-2007.ch8
var romFile []byte

func printFn(s string) {
	fmt.Println("print: ", s)
}

func main() {
	const (
		stackSize = 16
		delay     = 16
		RAMSize   = 4 * 1024
		w         = 64
		h         = 32
	)
	_d := NewDisplay(w, h)
	d = _d

	_emu, err := chipper.NewEmulator(stackSize, RAMSize, d)
	if err != nil {
		log.Println("error: ", err)
		return
	}

	emu = _emu
	d = _d

	printer := js.FuncOf(func(this js.Value, args []js.Value) any {
		a := make([]string, len(args))
		for k, arg := range args {
			a[k] = arg.String()
		}

		printFn(fmt.Sprintf("[%s]: %s", this.String(), strings.Join(a, "||")))

		return nil
	})

	sendDisplayToWASM := js.FuncOf(func(this js.Value, args []js.Value) any {
		dd, ok := d.(*Display)
		if !ok {
			fmt.Println("could not cast")
			panic("could not cast")
		}

		if len(args) != 1 {
			fmt.Println("expected at least 1 argument, got: ", len(args))
			return 0
		}

		ptr := args[0]
		return js.CopyBytesToJS(ptr, dd.data)
	})

	loadROMFn := js.FuncOf(func(this js.Value, args []js.Value) any {
		r := bytes.NewReader(romFile)
		if err := emu.Load(r); err != nil {
			fmt.Println("error loading rom: ", err)

			return nil
		}

		return nil
	})

	startCh := make(chan struct{})
	go func() {
		<-startCh
		fmt.Println("starting")
		for {
			err := emu.Tick()
			if errors.Is(err, io.EOF) {
				return
			}

			if err != nil {
				panic(fmt.Sprintf("runUntilError: %v", err))
			}

			time.Sleep(delay)
		}
	}()

	js.Global().Set("PrinterFn", printer)
	js.Global().Set("LoadROM", loadROMFn)
	js.Global().Set("GetDisplay", sendDisplayToWASM)

	startFn := js.FuncOf(func(this js.Value, args []js.Value) any {
		startCh <- struct{}{}
		return nil
	})

	js.Global().Set("StartEmu", startFn)

	select {}
}
