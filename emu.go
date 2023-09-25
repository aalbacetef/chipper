package chipper

import (
	"errors"
	"fmt"
	"io"
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
	// pc := emu.PC

	// fetch
	instrBytes, err := emu.Fetch(InstructionSize)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return fmt.Errorf("reached last instruction: %w", io.EOF)
		}
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
	fmt.Println("instruction: ", instr)
	//fmt.Printf("(%#0x) %+v\n", pc, instr)
	execErr := emu.Execute(instr)
	if execErr != nil {
		//fmt.Println("execution error: ", execErr)
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
