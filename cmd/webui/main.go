//go:build js && wasm

package main

import (
	"context"
	_ "embed"
	"fmt"
	"syscall/js"
	"time"

	"github.com/aalbacetef/chipper"
)

func main() {
	const (
		stackSize  = 16
		delay      = 16
		RAMSize    = 4*1024 + 1
		w          = 64
		h          = 32
		tickPeriod = 16 * time.Millisecond
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wrapper, err := NewWrapper(w, h, stackSize, RAMSize)
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	sendDisplayToWASM := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			fmt.Println("expected at least 1 argument, got: ", len(args))
			return 0
		}

		return wrapper.sendDisplayToWASM(args[0])
	})

	loadROMFn := js.FuncOf(func(this js.Value, args []js.Value) any {
		const wantLen = 2
		n := len(args)

		if n != wantLen {
			fmt.Printf("want %d args, got %d\n", wantLen, n)
			return 0
		}

		buf := args[0]
		lenBytes := args[1].Int()

		if err := wrapper.loadROM(buf, lenBytes); err != nil {
			fmt.Println("error: ", err)
		}

		return 0
	})

	handleKeyPress := js.FuncOf(func(this js.Value, args []js.Value) any {
		const wantLen = 3
		n := len(args)

		if n != wantLen {
			fmt.Printf("want %d args, got %d", wantLen, n)
			return nil
		}

		key := args[0].Int()
		dir := args[2].Int()

		v := chipper.Direction(dir) == chipper.Down
		fmt.Printf("[main.go] key: %#0x || %t\n", key, v)

		wrapper.keySrc.Set(key, v)

		return nil
	})

	startFn := js.FuncOf(func(this js.Value, args []js.Value) any {
		wrapper.start(ctx, tickPeriod)
		return nil
	})

	stopFn := js.FuncOf(func(this js.Value, args []js.Value) any {
		wrapper.stop()
		return nil
	})

	restartFn := js.FuncOf(func(this js.Value, args []js.Value) any {
		wrapper.restart()
		return nil
	})

	js.Global().Set("RestartEmu", restartFn)
	js.Global().Set("StartEmu", startFn)
	js.Global().Set("StopEmu", stopFn)
	js.Global().Set("LoadROM", loadROMFn)
	js.Global().Set("GetDisplay", sendDisplayToWASM)
	js.Global().Set("SendKeyboardEvent", handleKeyPress)

	select {}
}
