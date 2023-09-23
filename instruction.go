package chipper

type Instruction struct {
	Op       Opcode
	Operands []int
}

func Execute(emu *Emulator, instr Instruction) error {
	/*
		args := instr.Operands

		switch instr.Op {
		default:
			return fmt.Errorf("bad instruction: %v", instr)

		case ExecNNN:
			return emu.executeSubroutine(args[0], args[1], args[2])
		case Clear:
			return emu.clear()
		case Return:
			return emu.returnFromSub()
		case JumpNNN:
			return emu.jumpNNN(args[0], args[1], args[2])
		case CallSub:
			return emu.callNNN(args[0], args[1], args[2])
		}
	*/
	return nil
}
