package chipper

import (
	"fmt"
	"time"
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

func (emu *Emulator) Execute(instr Instruction) error {
	args := instr.Operands

	switch instr.Op {
	default:
		fmt.Println("unimplemented")
		return nil //	return fmt.Errorf("bad instruction: %v", instr)

	case ExecNNN:
		return nil
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

func (emu *Emulator) clearScreen() error {
	b := emu.Display.Bounds()
	dx, dy := b.Dx(), b.Dy()

	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {
			if err := emu.Display.Set(x, y, ColorBlack); err != nil {
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

// addYToX will add VY to VX, setting VF if it overflows.
func (emu *Emulator) addYToX(x, y int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	if err := isInBounds(RegisterCount, y); err != nil {
		return err
	}

	val := int(emu.V[x]) + int(emu.V[y])
	if val > 0xFF {
		emu.V[0xF] = 1
	}

	// note: the typecase is enough, but I prefer to be explicit
	emu.V[x] = byte(val % 0x100)

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

	val := int(emu.V[y]) - int(emu.V[x])

	clearOnBorrow := 1
	if val < 0 {
		clearOnBorrow = 0
	}

	emu.V[x] = byte(val % 0x100)
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
	droppedBit := vy & 0x1
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
	}

	emu.V[0xF] = byte(clearOnBorrow)
	emu.V[x] = byte(val % 0x100)

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

	vy := emu.V[y]
	droppedBit := vy >> 7

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

	if err := isInBounds(int(newAddr), ramSize); err != nil {
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
func (emu *Emulator) drawSpriteInXY(x, y, n int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	if err := isInBounds(RegisterCount, y); err != nil {
		return err
	}

	b := emu.Display.Bounds()
	dx, dy := b.Dx(), b.Dy()

	posx, posy := int(emu.V[x]), int(emu.V[y])

	emu.V[0xF] = 0

	for yline := 0; yline < n; yline++ {
		addr := int(emu.Index) + yline
		pixels := emu.RAM[addr]

		for xline := 0; xline < 8; xline++ {
			v := (pixels >> (7 - xline)) & 1

			xpos := (posx + xline) % dx
			ypos := (posy + yline) % dx
			at := emu.Display.At(xpos, ypos)

			newPix := v ^ byte(at)
			if Color(v) == at {
				emu.V[0xF] = 1
			}

			if err := emu.Display.Set(xpos, ypos, Color(newPix)); err != nil {
				return fmt.Errorf("could not set (%d, %d, %s): %w", xpos, ypos, Color(newPix), err)
			}
			// mask := byte(0b10000000 >> xline)
			// val := (pixels & mask)
			// if val == byte(ColorBlack) {
			// 	continue
			// }
			//
			// c := Color(val)
			//
			// xpos := (posx + xline) % dx
			// ypos := (posy + yline) % dy
			//
			// pix := emu.Display.At(xpos, ypos)
			//
			// if pix == ColorWhite && pix == c {
			// 	emu.V[0xF] = 1
			// 	c = ColorBlack
			// }
			//
			// if err := emu.Display.Set(xpos, ypos, c); err != nil {
			// 	return fmt.Errorf("could not set (%d, %d): %w", xpos, ypos, err)
			// }
		}

		if (posy + yline) >= dy {
			break
		}
	}

	return nil
}

func (emu *Emulator) skipIfKeyInXIsPressed(x int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	key := emu.Keys[emu.V[x]]
	if key {
		emu.PC += InstructionSize
	}

	return nil
}

func (emu *Emulator) skipIfKeyInXNotPressed(x int) error {
	if err := isInBounds(RegisterCount, x); err != nil {
		return err
	}

	key := emu.Keys[emu.V[x]]
	if !key {
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

	l := Listener{
		EventType: "KeyEvent",
		ID:        "waitForKeyAndStoreInX",
		ch:        make(chan Event, 1),
	}

	emu.listeners.Add(l)
	ev := <-l.ch
	emu.listeners.Del(l.ID)

	for {
		for k, key := range emu.Keys {
			if key {
				emu.V[x] = byte(k)
				return nil
			}
		}
	}
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

	c := int(emu.V[x]) * 5
	emu.Index = uint16(c)

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
