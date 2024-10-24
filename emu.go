package chipper

import (
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
	LastInstruction Instruction
	logger          *log.Logger
	lastUpdate      time.Time
}

func (emu *Emulator) SetLogger(l *log.Logger) {
	emu.logger = l
}

func (emu *Emulator) Close() { /* noop for now */ }

func NewRAM(size int) ([]byte, error) {
	if size <= 0 {
		return nil, fmt.Errorf("size must be > 0, got %d", size)
	}

	return make([]byte, size), nil
}

func NewEmulator(stackSize, ramSize int, display Display, keys KeyInputSource) (*Emulator, error) {
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
		Display: display,
	}

	if err := loadSprites(emu); err != nil {
		return nil, fmt.Errorf("could not load sprites into emulator: %w", err)
	}

	return emu, nil
}

// Load will read the ROM from the passed in io.Reader and load the sprites.
func (emu *Emulator) Load(r io.Reader) error {
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

// subtractTimers is our routine for reducing the timers. It is meant to run at 60Hz (~16ms).
func (emu *Emulator) subtractTimers() {
	const timerPeriod = 16 * time.Millisecond

	if emu.lastUpdate.IsZero() || time.Since(emu.lastUpdate) < timerPeriod {
		return
	}

	elapsed := time.Since(emu.lastUpdate)
	if elapsed < timerPeriod {
		return
	}

	emu.lastUpdate = time.Now()

	times := elapsed / timerPeriod
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

// Tick is the core Fetch-Decode-Execute loop of the emulator.
func (emu *Emulator) Tick() error {
	if emu.lastUpdate.IsZero() {
		emu.lastUpdate = time.Now()
	}

	emu.subtractTimers()

	logger := emu.logger
	if logger == nil {
		logger = log.New(io.Discard, "[emu] ", log.Ltime)
	}

	logger.Printf("fetch (PC=%#0x)\n", emu.PC)

	instrBytes, err := emu.Fetch(InstructionSize)
	if errors.Is(err, io.EOF) {
		return fmt.Errorf("reached last instruction: %w", io.EOF)
	}

	if err != nil {
		return fmt.Errorf("error fetching instruction: %w", err)
	}

	// update PC
	emu.PC += uint16(InstructionSize)

	logger.Println("decoding")

	// decode
	instr, err := Decode(instrBytes)
	if err != nil {
		return fmt.Errorf("could not decode instruction: %w", err)
	}

	// store the last instruction, useful for debugging
	emu.LastInstruction = instr

	logger.Println("executing instruction: ", instr.String())

	// execute
	execErr := emu.Execute(instr)
	if execErr != nil {
		return execErr
	}

	return nil
}

// Fetch will read the instruction pointed at by the PC. It will do a bounds check.
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
