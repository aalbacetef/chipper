package chipper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"
)

const (
	ProgramCounterSize = 2     // Size in bytes.
	RegisterCount      = 16    // V0-VF.
	StartAddress       = 0x200 // starting address of PC.
	NumKeys            = 16
	InstructionSize    = 2 // each instruction is 2 bytes wide.
)

type Emulator struct {
	PC              uint16
	V               [RegisterCount]byte
	DelayTimer      byte
	SoundTimer      byte
	Index           uint16
	Keys            KeyInputSource
	Stack           *Stack
	RAM             []byte
	Display         Display
	cancelFns       []context.CancelFunc
	LastInstruction Instruction
	logger          *log.Logger
	lastTick        time.Time
}

func (emu *Emulator) SetLogger(l *log.Logger) {
	emu.logger = l
}

func (emu *Emulator) Close() {
	for _, cancel := range emu.cancelFns {
		cancel()
	}
}

func NewRAM(size int) ([]byte, error) {
	if size <= 0 {
		return nil, fmt.Errorf("size must be > 0, got %d", size)
	}

	return make([]byte, size), nil
}

func NewEmulator(stackSize, ramSize int, d Display, keys KeyInputSource) (*Emulator, error) {
	stack, err := NewStack(stackSize)
	if err != nil {
		return nil, fmt.Errorf("could not create stack: %w", err)
	}

	ram, err := NewRAM(ramSize)
	if err != nil {
		return nil, fmt.Errorf("could not create ram: %w", err)
	}

	emu := &Emulator{
		PC:      StartAddress,
		Stack:   stack,
		RAM:     ram,
		Keys:    keys,
		Display: d,
	}

	if err := loadSprites(emu); err != nil {
		return nil, fmt.Errorf("could not load sprites into emulator: %w", err)
	}

	return emu, nil
}

func (emu *Emulator) subtractTimers() {
	period := time.Millisecond * 16
	elapsed := time.Since(emu.lastTick)
	if elapsed < period {
		return
	}

	times := elapsed / period
	sub := int(times)

	if emu.DelayTimer > 0 {
		dt := int(emu.DelayTimer) - sub
		if dt < 0 {
			dt = 0
		}
		emu.DelayTimer = byte(dt)
	}

	if emu.SoundTimer > 0 {
		st := int(emu.SoundTimer) - sub
		if st < 0 {
			st = 0
		}
		emu.SoundTimer = byte(st)
	}
}

// Tick .
func (emu *Emulator) Tick() error {
	if emu.lastTick.IsZero() {
		emu.lastTick = time.Now()
	}

	emu.subtractTimers()

	logger := emu.logger
	if logger == nil {
		logger = log.New(io.Discard, "[emu] ", log.Ltime)
	}

	logger.Println("fetch")
	instrBytes, err := emu.Fetch(InstructionSize)
	if errors.Is(err, io.EOF) {
		return fmt.Errorf("reached last instruction: %w", io.EOF)
	}

	if err != nil {
		log.Printf("error fetching: %v", err)
		return fmt.Errorf("error fetching instruction: %w", err)
	}

	// update PC
	emu.PC += uint16(InstructionSize)

	logger.Println("decoding")

	// decode
	instr, err := Decode(instrBytes)
	if err != nil {
		log.Printf("error decoding: %v", err)
		fmt.Println("could not decode : ", err)
		return fmt.Errorf("could not decode instruction: %w", err)
	}

	emu.LastInstruction = instr

	logger.Println("executing instruction: ", instr.String())

	// execute
	execErr := emu.Execute(instr)
	if execErr != nil {
		logger.Printf("error executing: %v", err)
		fmt.Println("could not execute: ", err)
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
