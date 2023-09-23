package chipper

import "fmt"

type Instruction struct {
	Op       Opcode
	Operands []int
}

func (instr Instruction) String() string {
	return fmt.Sprintf(
		"{Op: %s, Operands: [%#0x, %#0x, %#0x]}",
		instr.Op,
		instr.Operands[0],
		instr.Operands[1],
		instr.Operands[2],
	)
}

// @NOTE: consider returning an error, or not having to do this outright.
func addrFrom3(args []int) uint16 {
	return 0
}

func (emu *Emulator) Execute(instr Instruction) error {
	args := instr.Operands

	switch instr.Op {
	default:
		return nil //	return fmt.Errorf("bad instruction: %v", instr)

	case Clear:
		return emu.clearScreen()
	case Return:
		return emu.returnFromSub()
	case JumpNNN:
		return emu.jumpNNN(args)
		// case CallSub:
		// 	return emu.callNNN(args[0], args[1], args[2])
	}
}

func (emu *Emulator) executeSubroutine(n0, n1, n2 int) error {
	addr := (uint16(n0) << 8) | (uint16(n1) << 4) | uint16(n2)

	fmt.Printf("addr: %#0x\n", addr)

	//value := emu.RAM[int(addr)]
	//fmt.Printf("value: %#0x\n", value)

	return nil
}

func (emu *Emulator) clearScreen() error {
	b := emu.Display.Bounds()
	dx, dy := b.Dx(), b.Dy()

	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {
			if err := emu.Display.Set(x, y, ColorWhite); err != nil {
				return fmt.Errorf("could not clear pixel at (%d, %d): %w", x, y, err)
			}
		}
	}

	return nil
}

func (emu *Emulator) returnFromSub() error {
	retAddr, err := emu.Stack.Pop()
	if err != nil {
		return fmt.Errorf("returnFromSub: %w", err)
	}

	emu.PC = retAddr

	return nil
}

func (emu *Emulator) jumpNNN(args []int) error {
	return nil
}
