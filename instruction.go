package chipper

import (
	"fmt"
)

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

type InstructionNotImplementedError struct {
	op Opcode
}

func (e InstructionNotImplementedError) Error() string {
	return fmt.Sprintf("instruction not implemented: '%s'", e.op)
}

func (emu *Emulator) Execute(instr Instruction) error { //nolint: funlen,cyclop,gocyclo
	args := instr.Operands

	switch instr.Op {
	default:
		return nil

	case ExecNNN:
		return emu.execNNN(args, instr)

	case Clear:
		return emu.clearScreen()

	case ReturnFromSub:
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
		return emu.storeNNInX(args[0], args[1:])

	case AddNNToX:
		return emu.addNNToX(args[0], args[1:])

	case StoreYinX:
		return emu.storeYinX(args[0], args[1])

	case SetXToXORY:
		return emu.setXToXORY(args[0], args[1])

	case SetXToXANDY:
		return emu.setXToXANDY(args[0], args[1])

	case SetXToXXORY:
		return emu.setXToXXORY(args[0], args[1])

	case AddYToX:
		return emu.addYToX(args[0], args[1])

	case SubYFromX:
		return emu.subYFromX(args[0], args[1])

	case StoreYShiftedRightInX:
		return emu.storeYShiftedRightInX(args[0], args[1])

	case SetXToYMinusX:
		return emu.setXToYMinusX(args[0], args[1])

	case StoreYShiftedLeftInX:
		return emu.storeYShiftedLeftInX(args[0], args[1])

	case SkipIfXNotEqY:
		return emu.skipIfXNotEqY(args[0], args[1])

	case StoreMemAddrNNNInRegI:
		return emu.storeMemAddrNNNInRegI(args)

	case JumpToAddrNNNPlusV0:
		return emu.jumpToAddrNNNPlusV0(args)

	case SetXToRandomNumWithMaskNN:
		return emu.setXToRandomNumWithMaskNN(args[0], args[1:])

	case DrawSpriteInXY:
		return emu.drawSpriteInXY(args[0], args[1], args[2])

	case SkipIfKeyInXIsPressed:
		return emu.skipIfKeyInXIsPressed(args[0])

	case SkipIfKeyInXNotPressed:
		return emu.skipIfKeyInXNotPressed(args[0])

	case StoreValDTInX:
		return emu.storeValDTInX(args[0])

	case WaitForKeyAndStoreInX:
		return emu.waitForKeyAndStoreInX(args[0])

	case SetDTToX:
		return emu.setDTToX(args[0])

	case SetSTToX:
		return emu.setSTToX(args[0])

	case AddXToI:
		return emu.addXToI(args[0])

	case SetIToMemAddrOfSpriteInX:
		return emu.setIToMemAddrOfSpriteInX(args[0])

	case StoreBCDOfXInI:
		return emu.storeBCDOfXInI(args[0])

	case Store0ToXInI:
		return emu.store0ToXInI(args[0])

	case Fill0ToXWithValueInAddrI:
		return emu.fill0ToXWithValueInAddrI(args[0])
	}
}

func (emu *Emulator) execNNN(_ []int, instr Instruction) error {
	return InstructionNotImplementedError{instr.Op}
}

func (emu *Emulator) clearScreen() error {
	b := emu.Display.Bounds()
	dx, dy := b.Dx(), b.Dy()

	clearColor := emu.Display.ColorClear()

	for y := 0; y < dy; y++ {
		for x := 0; x < dx; x++ {
			emu.Display.Set(x, y, clearColor)
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
func (emu *Emulator) storeNNInX(x int, args []int) error {
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
func (emu *Emulator) addNNToX(x int, args []int) error {
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

// addYToX will add VY to VX, setting VF if it overflows and keeping the lower 8 bits.
func (emu *Emulator) addYToX(x, y int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	if err := isInBounds(RegisterCount, y); err != nil {
		return err
	}

	val := int(emu.V[x]) + int(emu.V[y])

	emu.V[0xF] = 0
	if val > max8BitVal {
		emu.V[0xF] = 1
	}

	emu.V[x] = byte(val & max8BitVal)

	return nil
}

// subYFromX will subtract VY from VX, storing the value in VX. It
// will clear VF if it borrows.
func (emu *Emulator) subYFromX(x, y int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	if err := isInBounds(RegisterCount, y); err != nil {
		return err
	}

	val := int(emu.V[x]) - int(emu.V[y])

	clearOnBorrow := 1
	if val < 0 {
		clearOnBorrow = 0
		val = val * -1
	}

	emu.V[x] = byte(val)
	emu.V[0xF] = byte(clearOnBorrow)

	return nil
}

// storeYShiftedRightInX will shift VY right and store it in X.
// It will place the dropped bit in VF.
func (emu *Emulator) storeYShiftedRightInX(x, y int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	if err := isInBounds(RegisterCount, y); err != nil {
		return err
	}

	vy := emu.V[y]

	const bitMask = 0x1
	droppedBit := vy & bitMask

	shifted := vy >> 1
	emu.V[0xF] = droppedBit
	emu.V[x] = shifted

	return nil
}

// SetXToYMinusX will set VX to VY - VX, clearing VF if a borrow occurs.
func (emu *Emulator) setXToYMinusX(x, y int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	if err := isInBounds(RegisterCount, y); err != nil {
		return err
	}

	vx, vy := emu.V[x], emu.V[y]
	val := int(vy) - int(vx)

	clearOnBorrow := 1
	if val < 0 {
		clearOnBorrow = 0
		val = -1 * val
	}

	emu.V[0xF] = byte(clearOnBorrow)
	emu.V[x] = byte(val)

	return nil
}

// storeYShiftedLeftInX will shift VY left and store it in VX.
// It will place the dropped bit in VF.
func (emu *Emulator) storeYShiftedLeftInX(x, y int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	if err := isInBounds(RegisterCount, y); err != nil {
		return err
	}

	const shiftBits = 7

	vy := emu.V[y]
	droppedBit := vy >> shiftBits

	emu.V[0xF] = droppedBit
	emu.V[x] = vy << 1

	return nil
}

// skipIfXNotEqY will skip to the next instruction if VX != VY.
func (emu *Emulator) skipIfXNotEqY(x, y int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	if err := isInBounds(RegisterCount, y); err != nil {
		return err
	}

	if emu.V[x] != emu.V[y] {
		emu.PC += InstructionSize
	}

	return nil
}

// storeMemAddrNNNInRegI will store the address NNN in the Index register.
func (emu *Emulator) storeMemAddrNNNInRegI(args []int) error {
	addr, err := ToAddr3(args)
	if err != nil {
		return err
	}

	emu.Index = addr

	return nil
}

// jumpToAddrNNNPlusV0 will JUMP to the address NNN + V0.
func (emu *Emulator) jumpToAddrNNNPlusV0(args []int) error {
	addr, err := ToAddr3(args)
	if err != nil {
		return err
	}

	newAddr := addr + uint16(emu.V[0])
	ramSize := len(emu.RAM)

	if err := isInBounds(ramSize, int(newAddr)); err != nil {
		return err
	}

	emu.PC = newAddr

	return nil
}

// setXToRandomNumWithMaskNN will set VX to (randInt(0, 255)  & NN).
func (emu *Emulator) setXToRandomNumWithMaskNN(x int, args []int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	val, err := ToByte(args)
	if err != nil {
		return err
	}

	rn := randomNum()
	emu.V[x] = rn & val

	return nil
}

// drawSpriteInXY NEEDS TO BE DONE PROPERLY.
func (emu *Emulator) drawSpriteInXY(x, y, n int) error { //nolint: gocognit
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	if err := isInBounds(RegisterCount, y); err != nil {
		return err
	}

	b := emu.Display.Bounds()
	displayWidth, displayHeight := b.Dx(), b.Dy()

	posx, posy := int(emu.V[x]), int(emu.V[y])

	emu.V[0xF] = 0
	clearColor := emu.Display.ColorClear()
	setColor := emu.Display.ColorSet()

	const shiftMask = 7

	for yline := 0; yline < n; yline++ {
		addr := int(emu.Index) + yline
		pixels := emu.RAM[addr]

		for xline := 0; xline < 8; xline++ {
			value := (pixels >> (shiftMask - xline)) & 1

			xpos := (posx + xline) % displayWidth
			ypos := (posy + yline) % displayWidth
			at := emu.Display.At(xpos, ypos)
			bit := byte(1)

			if ColorEq(at, clearColor) {
				bit = 0
			}

			newPix := value ^ bit

			if value == bit {
				emu.V[0xF] = 1
			}

			c := clearColor
			if newPix == 1 {
				c = setColor
			}

			emu.Display.Set(xpos, ypos, c)
		}

		if (posy + yline) >= displayHeight {
			break
		}
	}

	return nil
}

func (emu *Emulator) skipIfKeyInXIsPressed(x int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	fmt.Println("skip if key pressed: ", emu.V[x])

	v := emu.Keys.Get(int(emu.V[x]))
	if v {
		emu.PC += InstructionSize
	}

	return nil
}

func (emu *Emulator) skipIfKeyInXNotPressed(x int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	fmt.Println("skip if key not pressed: ", emu.V[x])

	v := emu.Keys.Get(int(emu.V[x]))
	if !v {
		emu.PC += InstructionSize
	}

	return nil
}

func (emu *Emulator) storeValDTInX(x int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	emu.V[x] = emu.DelayTimer

	return nil
}

func (emu *Emulator) waitForKeyAndStoreInX(x int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	fmt.Println("[instruction.go] waiting for key")

	key := <-emu.Keys.WaitUntilKeypress()

	fmt.Println("[instruction.go] got key: ", key)
	emu.V[x] = byte(key)

	return nil
}

func (emu *Emulator) setDTToX(x int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	emu.DelayTimer = emu.V[x]

	return nil
}

func (emu *Emulator) setSTToX(x int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	emu.SoundTimer = emu.V[x]

	return nil
}

func (emu *Emulator) addXToI(x int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	emu.Index += uint16(emu.V[x])

	return nil
}

func (emu *Emulator) setIToMemAddrOfSpriteInX(x int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	const width = 5

	emu.Index = uint16(emu.V[x]) * width

	return nil
}

func (emu *Emulator) storeBCDOfXInI(x int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	val := emu.V[x]

	bcd, err := bcdOfInt(int(val))
	if err != nil {
		return err
	}

	addr := int(emu.Index)

	for k, p := range bcd {
		emu.RAM[addr+k] = p
	}

	return nil
}

func (emu *Emulator) store0ToXInI(x int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	addr := int(emu.Index)

	for k, p := range emu.V[:x] {
		emu.RAM[addr+k] = p
	}

	return nil
}

func (emu *Emulator) fill0ToXWithValueInAddrI(x int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	addr := int(emu.Index)
	for k := 0; k < x; k++ {
		emu.V[k] = emu.RAM[addr+k]
	}

	return nil
}
