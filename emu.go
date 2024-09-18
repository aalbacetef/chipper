package chipper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
)

const (
	ProgramCounterSize = 2     // Size in bytes.
	RegisterCount      = 16    // V0-VF.
	StartAddress       = 0x200 // starting address of PC.
	NumKeys            = 16
	InstructionSize    = 2 // each instruction is 2 bytes wide.
)

type Emulator struct {
	PC         uint16
	V          [RegisterCount]byte
	DelayTimer byte
	SoundTimer byte
	Index      uint16
	Keys       [NumKeys]bool
	Stack      *Stack
	RAM        []byte
	Display    *Display
	cancelFns  []context.CancelFunc
	listeners  Listeners
}

func (emu *Emulator) Close() {
	for _, cancel := range emu.cancelFns {
		cancel()
	}
}

type EventType string

type Listener struct {
	ID        string
	EventType EventType
	ch        chan Event
}

type Listeners struct {
	listeners []Listener
	mu        sync.Mutex
}

func (ls *Listeners) All() []Listener {
	ls.mu.Lock()
	listeners := make([]Listener, len(ls.listeners))
	copy(listeners, ls.listeners)
	ls.mu.Unlock()

	return listeners
}

func (ls *Listeners) Del(id string) {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	listeners := make([]Listener, 0, len(ls.listeners)-1)
	for _, l := range ls.listeners {
		if l.ID == id {
			continue
		}

		listeners = append(listeners, l)
	}

	ls.listeners = listeners
}

func (ls *Listeners) Add(l Listener) {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	ls.listeners = append(ls.listeners, l)
}

type Event struct{}

func (emu *Emulator) RegisterEventSource(src <-chan Event) {
	ctx, cancel := context.WithCancel(context.Background())
	emu.cancelFns = append(emu.cancelFns, cancel)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case ev := <-src:
				for _, l := range emu.listeners.All() {
					l.ch <- ev
				}
			}
		}
	}()
}

func NewRAM(size int) ([]byte, error) {
	if size <= 0 {
		return nil, fmt.Errorf("size must be > 0, got %d", size)
	}

	return make([]byte, size), nil
}

func NewEmulator(stackSize, ramSize, w, h int) (*Emulator, error) {
	stack, err := NewStack(stackSize)
	if err != nil {
		return nil, fmt.Errorf("could not create stack: %w", err)
	}

	ram, err := NewRAM(ramSize)
	if err != nil {
		return nil, fmt.Errorf("could not create ram: %w", err)
	}

	display, err := NewDisplay(w, h)
	if err != nil {
		return nil, fmt.Errorf("could not create display: %w", err)
	}

	emu := &Emulator{
		PC:      StartAddress,
		Stack:   stack,
		RAM:     ram,
		Display: display,
	}

	if err := loadSprites(emu); err != nil {
		return nil, fmt.Errorf("could not load sprites into emulator: %w", err)
	}

	return emu, nil
}

// Tick .
func (emu *Emulator) Tick() error {
	instrBytes, err := emu.Fetch(InstructionSize)
	if errors.Is(err, io.EOF) {
		return fmt.Errorf("reached last instruction: %w", io.EOF)
	}

	if err != nil {
		return fmt.Errorf("error fetching instruction: %w", err)
	}

	// update PC
	emu.PC += uint16(InstructionSize)

	// decode
	instr, err := Decode(instrBytes)
	if err != nil {
		return fmt.Errorf("could not decode instruction: %w", err)
	}

	// execute
	execErr := emu.Execute(instr)
	if execErr != nil {
		return execErr
	}

	return nil
}

// fetch .
func (emu *Emulator) Fetch(numBytes int) ([]byte, error) {
	pc := int(emu.PC)
	ramSize := len(emu.RAM)

	if pc+numBytes >= ramSize {
		return nil, fmt.Errorf(
			"out of bounds (PC=%d, numBytes=%d, RAMSize=%d): %w",
			pc, numBytes, ramSize, io.EOF,
		)
	}

	read := make([]byte, numBytes)
	for k := 0; k < numBytes; k++ {
		read[k] = emu.RAM[pc+k]
	}

	return read, nil
}

// Load .
func (emu *Emulator) Load(r io.Reader) error {
	if err := loadSprites(emu); err != nil {
		return fmt.Errorf("could not load sprites: %w", err)
	}

	ramSize := len(emu.RAM)
	maxSize := ramSize - StartAddress
	p := make([]byte, maxSize)

	bytesRead, err := r.Read(p)
	if err != nil {
		return fmt.Errorf("error reading ROM: %w", err)
	}

	p = p[:bytesRead]

	for k, b := range p {
		emu.RAM[int(StartAddress)+k] = b
	}

	return nil
}
