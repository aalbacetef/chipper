package chipper

import "fmt"

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

	return &Emulator{
		PC:      StartAddress,
		Stack:   stack,
		RAM:     ram,
		Display: display,
	}, nil
}

func (emu *Emulator) Tick() error {
	// fetch
	instrBytes, err := emu.fetch(InstructionSize)
	if err != nil {
		return fmt.Errorf("error fetching instruction: %w", err)
	}

	// update PC
	emu.PC += uint16(InstructionSize)

	// decode
	instr, err := Decode(toUint16(instrBytes))
	if err != nil {
		return fmt.Errorf("could not decode instruction: %w", err)
	}

	// execute
	fmt.Println("instruction: ", instr)

	return fmt.Errorf("not implemented")
}

func toUint16(b []byte) uint16 {
	return (uint16(b[0]) << 2) | uint16(b[1])
}

func (emu *Emulator) fetch(numBytes int) ([]byte, error) {
	pc := int(emu.PC)
	ramSize := len(emu.RAM)

	if pc+numBytes >= ramSize {
		return nil, fmt.Errorf(
			"out of bounds (PC=%d, numBytes=%d, RAMSize=%d)",
			pc, numBytes, ramSize,
		)
	}

	read := make([]byte, numBytes)
	for k := 0; k < numBytes; k++ {
		read[k] = emu.RAM[pc+k]
	}

	return read, nil
}
