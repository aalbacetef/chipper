//go:build js && wasm

package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"syscall/js"
	"time"

	"github.com/aalbacetef/chipper"
)

func NewWrapper(w, h, stackSize, RAMSize int) (*WASMWrapper, error) {
	wrapper := &WASMWrapper{}
	wrapper.settings = Settings{
		stackSize: stackSize,
		ramSize:   RAMSize,
		w:         w,
		h:         h,
	}

	if err := wrapper.init(); err != nil {
		return nil, fmt.Errorf("could not initialize: %w", err)
	}

	return wrapper, nil
}

type Settings struct {
	stackSize int
	ramSize   int
	w         int
	h         int
}

type WASMWrapper struct {
	settings   Settings
	emu        *chipper.Emulator
	d          *Display
	keySrc     chipper.KeyInputSource
	cancelFunc context.CancelFunc
	mu         sync.Mutex
}

func (wrapper *WASMWrapper) loadROM(buf js.Value, lenBytes int) error {
	romFile := make([]byte, lenBytes)
	js.CopyBytesToGo(romFile, buf)

	r := bytes.NewReader(romFile)
	if err := wrapper.emu.Load(r); err != nil {
		return fmt.Errorf("could not load rom: %w", err)
	}

	return nil
}

func (wrapper *WASMWrapper) sendDisplayToWASM(ptr js.Value) int {
	wrapper.d.mu.Lock()
	bytesCopied := js.CopyBytesToJS(ptr, wrapper.d.data)
	wrapper.d.mu.Unlock()

	return bytesCopied
}

func (wrapper *WASMWrapper) start(mainCtx context.Context, period time.Duration) {
	wrapper.mu.Lock()
	if wrapper.cancelFunc != nil {
		wrapper.cancelFunc()
	}

	ctx, cancel := context.WithCancel(mainCtx)
	wrapper.cancelFunc = cancel
	wrapper.mu.Unlock()

	ticker := time.NewTicker(period)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := wrapper.emu.Tick()
				if errors.Is(err, io.EOF) {
					fmt.Println("end of file, exiting")
					return
				}

				if err != nil {
					pc := wrapper.emu.PC
					indx := wrapper.emu.Index
					last := wrapper.emu.LastInstruction

					fmt.Println("error: ", err)
					fmt.Printf("PC: %#0x | (%d) \n: ", pc, pc)
					fmt.Printf("Index: %#0x\n", indx)
					fmt.Println("last instruction: ", last)

					return
				}

			}
		}
	}()
}

func (wrapper *WASMWrapper) stop() {
	wrapper.mu.Lock()
	defer wrapper.mu.Unlock()

	if wrapper.cancelFunc != nil {
		wrapper.cancelFunc()
	}
}

func (wrapper *WASMWrapper) reset() {
	wrapper.stop()

	wrapper.mu.Lock()
	wrapper.init()
	wrapper.mu.Unlock()

	return
}

func (wrapper *WASMWrapper) init() error {
	w := wrapper.settings.w
	h := wrapper.settings.h
	stackSize := wrapper.settings.stackSize
	ramSize := wrapper.settings.ramSize

	d := NewDisplay(w, h)
	keySrc := NewWebKeyInputSource()

	emu, err := chipper.NewEmulator(stackSize, ramSize, d, keySrc)
	if err != nil {
		return fmt.Errorf("could not start emulator: %w", err)
	}

	wrapper.emu = emu
	wrapper.d = d
	wrapper.keySrc = keySrc

	return nil
}
