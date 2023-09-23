package chipper

import "fmt"

type Opcode string

const (
	Unknown                   Opcode = "Unknown"
	Nop                       Opcode = "Nop"
	ExecNNN                   Opcode = "ExecNNN"
	Clear                     Opcode = "Clear"
	Return                    Opcode = "Return"
	JumpNNN                   Opcode = "JumpNNN"
	CallSub                   Opcode = "CallSub"
	SkipIfXEqNN               Opcode = "SkipIfXEqNN"
	SkipIfXNotEqNN            Opcode = "SkipIfXNotEqNN"
	SkipIfXEqY                Opcode = "SkipIfXEqY"
	StoreNNInX                Opcode = "StoreNNInX"
	AddNNToX                  Opcode = "AddNNToX"
	StoreYinX                 Opcode = "StoreYinX"
	SetXToXORY                Opcode = "SetXToXORY"
	SetXToXANDY               Opcode = "SetXToXANDY"
	SetXToXXORY               Opcode = "SetXToXXORY"
	AddYToX                   Opcode = "AddYToX"
	SubYFromX                 Opcode = "SubYFromX"
	StoreYShiftedRightInX     Opcode = "StoreYShiftedRightInX"
	SetXToYMinusX             Opcode = "SetXToYMinusX"
	StoreYShiftedLeftInX      Opcode = "StoreYShiftedLeftInX"
	SkipIfXNotEqY             Opcode = "SkipIfXNotEqY"
	StoreMemAddrNNNInRegI     Opcode = "StoreMemAddrNNNInRegI"
	JumpToAddrNNNPlusV0       Opcode = "JumpToAddrNNNPlusV0"
	SetXToRandomNumWithMaskNN Opcode = "SetXToRandomNumWithMaskNN"
	DrawSpriteInXY            Opcode = "DrawSpriteInXY"
	SkipIfKeyInXIsPressed     Opcode = "SkipIfKeyInXIsPressed"
	SkipIfKeyInXNotPressed    Opcode = "SkipIfKeyInXNotPressed"
	StoreValDTInX             Opcode = "StoreValDTInX"
	WaitForKeyAndStoreInX     Opcode = "WaitForKeyAndStoreInX"
	SetDTToX                  Opcode = "SetDTToX"
	SetSTToX                  Opcode = "SetSTToX"
	AddXToI                   Opcode = "AddXToI"
	SetIToMemAddrOfSpriteInX  Opcode = "SetIToMemAddrOfSpriteInX"
	StoreBCDOfXInI            Opcode = "StoreBCDOfXInI"
	Store0ToXInI              Opcode = "Store0ToXInI"
	Fill0ToXWithValueInAddrI  Opcode = "Fill0ToXWithValueInAddrI"
)

// DetermineOpcode will return the appropriate Opcode given the digits passed in.
// It expects digits to have length 4.
func DetermineOpcode(digits []int) Opcode {
	first := digits[0]
	last := digits[3]

	switch first {
	case 0:
		if match(digits[1:], []int{0, 0, 0}) {
			return Nop
		}

		if match(digits[1:], []int{0, 0xE, 0}) {
			return Clear
		}

		if match(digits[1:], []int{0, 0xE, 0xE}) {
			return Return
		}

		return ExecNNN

	case 1:
		return JumpNNN

	case 2:
		return CallSub

	case 3:
		return SkipIfXEqNN

	case 4:
		return SkipIfXNotEqNN

	case 5:
		return SkipIfXEqY

	case 6:
		return StoreNNInX

	case 7:
		return AddNNToX

	case 8:
		switch last {
		case 0:
			return StoreYinX
		case 1:
			return SetXToXORY
		case 2:
			return SetXToXANDY
		case 3:
			return SetXToXXORY
		case 4:
			return AddYToX
		case 5:
			return SubYFromX
		case 6:
			return StoreYShiftedRightInX
		case 7:
			return SetXToYMinusX
		case 0xE:
			return StoreYShiftedLeftInX
		default:
			return Unknown
		}

	case 9:
		if last == 0 {
			return SkipIfXNotEqY
		}

		return Unknown

	case 0xA:
		return StoreMemAddrNNNInRegI

	case 0xB:
		return JumpToAddrNNNPlusV0

	case 0xC:
		return SetXToRandomNumWithMaskNN

	case 0xD:
		return DrawSpriteInXY

	case 0xE:
		lastTwo := digits[2:]
		if match(lastTwo, []int{9, 0xE}) {
			return SkipIfKeyInXIsPressed
		}

		if match(lastTwo, []int{0xA, 1}) {
			return SkipIfKeyInXNotPressed
		}

		return Unknown

	case 0xF:
		lastTwo := digits[2:]
		if match(lastTwo, []int{0, 7}) {
			return StoreValDTInX
		}

		if match(lastTwo, []int{0, 0xA}) {
			return WaitForKeyAndStoreInX
		}

		if match(lastTwo, []int{1, 5}) {
			return SetDTToX
		}

		if match(lastTwo, []int{1, 8}) {
			return SetSTToX
		}

		if match(lastTwo, []int{1, 0xE}) {
			return AddXToI
		}

		if match(lastTwo, []int{2, 9}) {
			return SetIToMemAddrOfSpriteInX
		}

		if match(lastTwo, []int{3, 3}) {
			return StoreBCDOfXInI
		}

		if match(lastTwo, []int{5, 5}) {
			return Store0ToXInI
		}

		if match(lastTwo, []int{6, 5}) {
			return Fill0ToXWithValueInAddrI
		}

		return Unknown

	default:
		return Unknown
	}
}

func match(in, out []int) bool {
	if len(in) != len(out) {
		return false
	}

	for k, v := range in {
		if out[k] != v {
			return false
		}
	}

	return true
}

func Decode(p []byte) (Instruction, error) {
	instr := toUint16(p)

	d0 := (instr & 0xF000) >> (3 * 4)
	d1 := (instr & 0x0F00) >> (2 * 4)
	d2 := (instr & 0x00F0) >> (1 * 4)
	d3 := (instr & 0x000F)

	digits := [4]int{
		int(d0), int(d1), int(d2), int(d3),
	}

	opcode := DetermineOpcode(digits[:])
	if opcode == Unknown {
		return Instruction{}, fmt.Errorf("unknown opcode: %#0x", instr)
	}

	return Instruction{
		Op:       opcode,
		Operands: digits[1:],
	}, nil
}
