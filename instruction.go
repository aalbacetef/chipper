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
		fmt.Println("unimplemented")
		return nil //	return fmt.Errorf("bad instruction: %v", instr)

	case Clear:
		return emu.clearScreen()

	case Return:
		return emu.returnFromSub()

	case JumpNNN:
		return emu.jumpNNN(args)

	case CallSub:
		return emu.callSubNNN(args)

	case SkipIfXEqNN:
		return emu.skipIfXEqNN(args[0], args[1:])

	case SkipIfXNotEqNN:
		return emu.skipIfXNotEqNN(args[0], args[1:])

	case SkipIfXEqY:
		return emu.skipIfXEqY(args[0], args[1])

	case StoreNNInX:
		return emu.storeNNInX(args[:2], args[2])

	case AddNNToX:
		return emu.addNNToX(args[:2], args[2])

	case StoreYinX:
		return emu.storeYinX(args[0], args[1])

	case SetXToXORY:
		return emu.setXToXORY(args[0], args[1])

	case SetXToXANDY:
		return emu.setXToXANDY(args[0], args[1])

	case SetXToXXORY:
		return emu.setXToXXORY(args[0], args[1])

	}
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
	addr, err := ToAddr3(args)
	if err != nil {
		return err
	}

	emu.PC = addr

	return nil
}

func (emu *Emulator) callSubNNN(args []int) error {
	addr, err := ToAddr3(args)
	if err != nil {
		return err
	}

	if err := emu.Stack.Push(emu.PC); err != nil {
		return fmt.Errorf("could not push PC on stack: %w", err)
	}

	emu.PC = addr

	return nil
}

// skipIfXEqNN will skip the next instruction if VX == NN.
func (emu *Emulator) skipIfXEqNN(x int, args []int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	vx := emu.V[x]
	value, err := ToByte(args)
	if err != nil {
		return err
	}

	if vx == value {
		emu.PC += InstructionSize
	}

	return nil
}

// skipIfXNotEqNN will skip if VX != NN.
func (emu *Emulator) skipIfXNotEqNN(x int, args []int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	vx := emu.V[x]
	value, err := ToByte(args)
	if err != nil {
		return err
	}

	if vx != value {
		emu.PC += InstructionSize
	}

	return nil
}

// skipIfXEqY will skip if VX == VY.
func (emu *Emulator) skipIfXEqY(x, y int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	if err := isInBounds(RegisterCount, y); err != nil {
		return err
	}

	vx, vy := emu.V[x], emu.V[y]

	if vx == vy {
		emu.PC += InstructionSize
	}

	return nil
}

// storeNNInX will store the value in VX.
func (emu *Emulator) storeNNInX(args []int, x int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	val, err := ToByte(args)
	if err != nil {
		return err
	}

	emu.V[x] = val

	return nil
}

// addNNToX will add NN to VX and store the (wrapped) result in VX.
func (emu *Emulator) addNNToX(args []int, x int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	val, err := ToByte(args)
	if err != nil {
		return err
	}

	emu.V[x] += val

	return nil
}

// storeYinX will store the value of VY in VX.
func (emu *Emulator) storeYinX(x, y int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	if err := isInBounds(RegisterCount, y); err != nil {
		return err
	}

	emu.V[x] = emu.V[y]

	return nil
}

// setXToXORY will set VX to VX | VY.
func (emu *Emulator) setXToXORY(x, y int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	if err := isInBounds(RegisterCount, y); err != nil {
		return err
	}

	emu.V[x] = (emu.V[x] | emu.V[y])

	return nil
}

// setXToXANDY will set VX to VX & VY.
func (emu *Emulator) setXToXANDY(x, y int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	if err := isInBounds(RegisterCount, y); err != nil {
		return err
	}

	emu.V[x] = (emu.V[x] & emu.V[y])

	return nil
}

// setXToXXORY will set VX to VX ^ VY.
func (emu *Emulator) setXToXXORY(x, y int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	if err := isInBounds(RegisterCount, y); err != nil {
		return err
	}

	emu.V[x] = (emu.V[x] ^ emu.V[y])

	return nil
}
