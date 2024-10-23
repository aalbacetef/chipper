//go:build js
// +build js

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
	emu    *chipper.Emulator
	d      chipper.Display
	KeySrc chipper.KeyInputSource
)

func printFn(s string) {
	fmt.Println("print: ", s)
}

func main() {
	const (
		stackSize = 16
		delay     = 16
		RAMSize   = 4*1024 + 1
		w         = 64
		h         = 32
	)

	_d := NewDisplay(w, h)
	d = _d

	KeySrc = NewWebKeyInputSource()

	_emu, err := chipper.NewEmulator(stackSize, RAMSize, d, KeySrc)
	if err != nil {
		log.Println("error: ", err)
		return
	}

	emu = _emu
	// emu.SetLogger(log.New(os.Stdout, "[emu] ", log.Ltime))

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

		dd.mu.Lock()
		defer dd.mu.Unlock()

		ptr := args[0]
		return js.CopyBytesToJS(ptr, dd.data)
	})

	loadROMFn := js.FuncOf(func(this js.Value, args []js.Value) any {
		const wantLen = 2
		n := len(args)

		if n != wantLen {
			fmt.Printf("want %d args, got %d\n", wantLen, n)
			return nil
		}

		buf := args[0]
		lenBytes := args[1].Int()

		romFile := make([]byte, lenBytes)
		js.CopyBytesToGo(romFile, buf)

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
				fmt.Println("ERROR: ", err)
				fmt.Printf("PC: %#0x | (%d) \n: ", emu.PC, emu.PC)
				fmt.Printf("Index: %#0x\n", emu.Index)
				fmt.Println("last instruction: ", emu.LastInstruction)
				panic(fmt.Sprintf("runUntilError: %v", err))
			}

			time.Sleep(delay)
		}
	}()

	handleKeyPress := js.FuncOf(func(this js.Value, args []js.Value) any {
		const wantLen = 3
		n := len(args)

		if n != wantLen {
			fmt.Printf("want %d args, got %d", wantLen, n)
			return nil
		}

		key := args[0].Int()
		// repeat := args[1].Bool()
		dir := args[2].Int()

		v := chipper.Direction(dir) == chipper.Down
		fmt.Printf("[main.go] key: %#0x || %t\n", key, v)

		KeySrc.Set(key, v)

		return nil
	})

	js.Global().Set("PrinterFn", printer)
	js.Global().Set("LoadROM", loadROMFn)
	js.Global().Set("GetDisplay", sendDisplayToWASM)
	js.Global().Set("SendKeyboardEvent", handleKeyPress)

	startFn := js.FuncOf(func(this js.Value, args []js.Value) any {
		startCh <- struct{}{}
		return nil
	})

	js.Global().Set("StartEmu", startFn)

	select {}
}
