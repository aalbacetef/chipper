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

func (emu *Emulator) Execute(instr Instruction) error {
	args := instr.Operands

	switch instr.Op {
	default:
		return nil //	return fmt.Errorf("bad instruction: %v", instr)

	case ExecNNN:
		return emu.executeSubroutine(args[0], args[1], args[2])
		/*
			case Clear:
				return emu.clear()
			case Return:
				return emu.returnFromSub()
			case JumpNNN:
				return emu.jumpNNN(args[0], args[1], args[2])
			case CallSub:
				return emu.callNNN(args[0], args[1], args[2])*/
	}
}

func (emu *Emulator) executeSubroutine(n0, n1, n2 int) error {
	addr := (uint16(n0) << 2) | (uint16(n1) << 1) | uint16(n2)

	fmt.Printf("addr: %#0x\n", addr)

	//value := emu.RAM[int(addr)]
	//fmt.Printf("value: %#0x\n", value)

	return nil
}
