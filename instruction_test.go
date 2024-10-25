package chipper

import "testing"

func mkEmu(t *testing.T) *Emulator {
	t.Helper()

	const (
		testStackSize = 16
		testRAMSize   = 4096
		testW         = 64
		testH         = 32
	)

	display, err := NewDebugDisplay(testW, testH)
	if err != nil {
		t.Fatalf("could not make debug display: %v", err)
	}

	emu, err := NewEmulator(
		testStackSize,
		testRAMSize,
		display,
		&StubKeyInputSource{},
	)
	if err != nil {
		t.Fatalf("could not create emulator: %v", err)
	}

	return emu
}

func TestInstruction(t *testing.T) {
	tests := []struct {
		label string
		fn    func(t *testing.T)
	}{
		{"clearScreen", testClearScreen},
		{"returnFromSub", testReturnFromSub},
		{"jumpNNN", testJumpNNN},
		{"callSubNNN", testCallSubNNN},
		{"skipIfXEqNN", testSkipIfXEqNN},
		{"skipIfXNotEqNN", testSkipIfXNotEqNN},
		{"skipIfXEqY", testSkipIfXEqY},
		{"storeNNInX", testStoreNNInX},
		{"addNNtoX", testAddNNToX},
		{"storeYInX", testStoreYInX},
		{"setXToXORY", testSetXToXORY},
		{"setXToXANDY", testSetXToXANDY},
		{"setXToXXORY", testSetXToXXORY},
		{"addYToX", testAddYToX},
		{"subYFromX", testSubYFromX},
		{"storeYShiftedRightInX", testStoreYShiftedRightInX},
		{"setXToYMinusX", testSetXToYMinusX},
		{"storeYShiftedLeftInX", testStoreYShiftedLeftInX},
		{"skipIfXNotEqY", testSkipIfXNotEqY},
		{"set vx to random number mask with nn", testSetVXWithMask},
	}

	for _, c := range tests {
		t.Run(c.label, c.fn)
	}
}

func testClearScreen(t *testing.T) {
	t.Helper()

	emu := mkEmu(t)

	colorSet := emu.Display.ColorSet()
	colorClear := emu.Display.ColorClear()

	if err := Each(emu.Display, func(x, y int) error {
		emu.Display.Set(x, y, colorSet)

		return nil
	}); err != nil {
		t.Fatalf("Each: %v", err)
	}

	if err := emu.clearScreen(); err != nil {
		t.Fatalf("error: %v", err)
	}

	if err := Each(emu.Display, func(x, y int) error {
		at := emu.Display.At(x, y)
		if !ColorEq(at, colorClear) {
			t.Fatalf("(%d, %d) not clear", x, y)
		}

		return nil
	}); err != nil {
		t.Fatalf("Each: %v", err)
	}
}

func testReturnFromSub(t *testing.T) {
	t.Helper()

	emu := mkEmu(t)

	const testAddr = 0x111
	if err := emu.Stack.Push(testAddr); err != nil {
		t.Fatalf("could not push: %v", err)
	}

	if err := emu.returnFromSub(); err != nil {
		t.Fatalf("error: %v", err)
	}

	if emu.PC != testAddr {
		t.Fatalf("got %#0x, want %#0x", emu.PC, testAddr)
	}
}

func testJumpNNN(t *testing.T) {
	t.Helper()

	emu := mkEmu(t)

	const testAddr = 0x111

	if err := emu.jumpNNN([]int{1, 1, 1}); err != nil {
		t.Fatalf("error: %v", err)
	}

	if emu.PC != testAddr {
		t.Fatalf("got %#0x, want %#0x", emu.PC, testAddr)
	}
}

func testCallSubNNN(t *testing.T) {
	t.Helper()

	emu := mkEmu(t)

	const (
		testAddr = 0x111
		origPC   = 0x210
	)

	emu.PC = origPC

	if err := emu.callSubNNN([]int{1, 1, 1}); err != nil {
		t.Fatalf("error: %v", err)
	}

	if emu.PC != testAddr {
		t.Fatalf("got %#0x, want %#0x", emu.PC, testAddr)
	}

	val, err := emu.Stack.Pop()
	if err != nil {
		t.Fatalf("could not pop stack: %v", err)
	}

	if val != origPC {
		t.Fatalf("got %#0x, want %#0x", val, origPC)
	}
}

func testSkipIfXEqNN(t *testing.T) {
	t.Helper()

	emu := mkEmu(t)

	const (
		x  = 1
		nn = 0x12
	)

	args := []int{0x1, 0x2}

	t.Run("check it skips when Vx equals nn", func(t *testing.T) {
		emu.V[x] = nn
		wantPC := emu.PC + InstructionSize

		if err := emu.skipIfXEqNN(x, args); err != nil {
			t.Fatalf("error: %v", err)
		}

		if emu.PC != wantPC {
			t.Fatalf("got %#0x, want %#0x", emu.PC, wantPC)
		}
	})

	t.Run("check it doesn't skip when Vx is not equal to nn", func(t *testing.T) {
		emu.V[x] = nn + 1
		wantPC := emu.PC

		if err := emu.skipIfXEqNN(x, args); err != nil {
			t.Fatalf("error: %v", err)
		}

		if emu.PC != wantPC {
			t.Fatalf("got #%0x, want #%0x", emu.PC, wantPC)
		}
	})
}

func testSkipIfXNotEqNN(t *testing.T) {
	t.Helper()

	emu := mkEmu(t)

	const (
		x  = 1
		nn = 0x12
	)

	args := []int{0x1, 0x2}

	t.Run("it skips if Vx not eq nn", func(t *testing.T) {
		emu.V[x] = nn + 1
		wantPC := emu.PC + InstructionSize

		if err := emu.skipIfXNotEqNN(x, args); err != nil {
			t.Fatalf("error: %v", err)
		}

		if emu.PC != wantPC {
			t.Fatalf("got %#0x, want %#0x", emu.PC, wantPC)
		}
	})

	t.Run("it doesn't skip if Vx eq nn", func(t *testing.T) {
		emu.V[x] = nn
		wantPC := emu.PC

		if err := emu.skipIfXNotEqNN(x, args); err != nil {
			t.Fatalf("error: %v", err)
		}

		if emu.PC != wantPC {
			t.Fatalf("got #%0x, want #%0x", emu.PC, wantPC)
		}
	})
}

func testSkipIfXEqY(t *testing.T) {
	t.Helper()

	emu := mkEmu(t)

	const (
		x       = 1
		y       = 2
		testVal = 5
	)

	t.Run("it skips if Vx eq Vy", func(t *testing.T) {
		emu.V[x] = testVal
		emu.V[y] = emu.V[x]
		wantPC := emu.PC + InstructionSize

		if err := emu.skipIfXEqY(x, y); err != nil {
			t.Fatalf("error: %v", err)
		}

		if emu.PC != wantPC {
			t.Fatalf("got %#0x, want %#0x", emu.PC, wantPC)
		}
	})

	t.Run("it doesn't skip if Vx not eq Vy", func(t *testing.T) {
		emu.V[x] = emu.V[y] + 1
		wantPC := emu.PC

		if err := emu.skipIfXEqY(x, y); err != nil {
			t.Fatalf("error: %v", err)
		}

		if emu.PC != wantPC {
			t.Fatalf("got %#0x, want %#0x", emu.PC, wantPC)
		}
	})
}

func testStoreNNInX(t *testing.T) {
	t.Helper()

	emu := mkEmu(t)

	const (
		x  = 1
		nn = 0x11
	)

	emu.V[x] = 0
	args := []int{0x1, 0x1}

	if err := emu.storeNNInX(x, args); err != nil {
		t.Fatalf("error: %v", err)
	}

	if emu.V[x] != nn {
		t.Fatalf("got %#0x, want %#0x", emu.V[x], nn)
	}
}

func testAddNNToX(t *testing.T) {
	t.Helper()

	emu := mkEmu(t)

	const (
		testX = 5
		testN = 0x1
	)

	if err := emu.addNNToX(testX, []int{testN, testN}); err != nil {
		t.Fatalf("error: %v", err)
	}

	got := emu.V[testX]
	want, err := ToByte([]int{testN, testN})

	if err != nil {
		t.Fatalf("could not convert to byte: %v", err)
	}

	if got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

func testStoreYInX(t *testing.T) {
	t.Helper()

	emu := mkEmu(t)

	const (
		x  = 1
		y  = 2
		nn = 5
	)

	emu.V[x] = nn - 1
	emu.V[y] = nn

	if err := emu.storeYinX(x, y); err != nil {
		t.Fatalf("error: %v", err)
	}

	if emu.V[x] != emu.V[y] {
		t.Fatalf("got %#0x, want %#0x", emu.V[x], emu.V[y])
	}
}

func testSetXToXORY(t *testing.T) {
	t.Helper()
	emu := mkEmu(t)

	const (
		x    = 1
		y    = 2
		valX = byte(0b00010001)
		valY = byte(0b00100010)
		want = byte(0b00110011)
	)

	emu.V[x] = valX
	emu.V[y] = valY

	if err := emu.setXToXORY(x, y); err != nil {
		t.Fatalf("error: %v", err)
	}

	if emu.V[x] != want {
		t.Fatalf("got %#0x, want %#0x", emu.V[x], want)
	}
}

func testSetXToXANDY(t *testing.T) {
	t.Helper()
	emu := mkEmu(t)

	const (
		x    = 1
		y    = 2
		valX = byte(0b00010001)
		valY = byte(0b00110011)
		want = byte(0b00010001)
	)

	emu.V[x] = valX
	emu.V[y] = valY

	if err := emu.setXToXANDY(x, y); err != nil {
		t.Fatalf("error: %v", err)
	}

	got := emu.V[x]
	if got != want {
		t.Fatalf("got %#0x, want %#0x", got, want)
	}
}

func testSetXToXXORY(t *testing.T) {
	t.Helper()
	emu := mkEmu(t)

	const (
		x    = 1
		y    = 2
		valX = byte(0b00010001)
		valY = byte(0b00110011)
		want = byte(0b00100010)
	)

	emu.V[x] = valX
	emu.V[y] = valY

	if err := emu.setXToXXORY(x, y); err != nil {
		t.Fatalf("error: %v", err)
	}

	got := emu.V[x]
	if got != want {
		t.Fatalf("got %#0x, want %#0x", got, want)
	}
}

func testAddYToX(t *testing.T) {
	t.Helper()
	emu := mkEmu(t)

	const (
		x    = 1
		y    = 2
		valX = 1
		valY = 10
	)

	emu.V[x] = valX
	emu.V[y] = valY

	t.Run("it adds VY to VX without overflow", func(t *testing.T) {
		want := valX + valY

		if err := emu.addYToX(x, y); err != nil {
			t.Fatalf("error: %v", err)
		}

		if emu.V[x] != byte(want) {
			t.Fatalf("got %#0x, want %#0x", emu.V[x], want)
		}

		if emu.V[0xF] != 0 {
			t.Fatalf("overflow flag was set")
		}
	})

	t.Run("it adds VY to VX with overflow", func(t *testing.T) {
		emu.V[x] = 0xFF
		emu.V[y] = valY

		want := (valY + emu.V[x]) & 0xFF

		if err := emu.addYToX(x, y); err != nil {
			t.Fatalf("error: %v", err)
		}

		if emu.V[x] != byte(want) {
			t.Fatalf("got %#0x, want %#0x", emu.V[x], want)
		}

		if emu.V[0xF] != 1 {
			t.Fatalf("overflow flag not set")
		}
	})
}

func testSubYFromX(t *testing.T) {
	t.Helper()
	emu := mkEmu(t)

	const (
		x    = 1
		y    = 2
		valX = 10
		valY = 1
	)

	emu.V[x] = valX
	emu.V[y] = valY

	t.Run("it subs VY from VX without borrow", func(t *testing.T) {
		want := valX - valY

		if err := emu.subYFromX(x, y); err != nil {
			t.Fatalf("error: %v", err)
		}

		if emu.V[x] != byte(want) {
			t.Fatalf("got %#0x, want %#0x", emu.V[x], want)
		}

		if emu.V[0xF] != 1 {
			t.Fatalf("borrow flag was cleared")
		}
	})

	t.Run("it subs VY from VX with borrow", func(t *testing.T) {
		emu.V[x] = 0x0F
		emu.V[y] = 0xFF

		want := 0xF0

		if err := emu.subYFromX(x, y); err != nil {
			t.Fatalf("error: %v", err)
		}

		if emu.V[x] != byte(want) {
			t.Fatalf("got %#0x, want %#0x", emu.V[x], want)
		}

		if emu.V[0xF] != 0 {
			t.Fatalf("borrow flag was not cleared")
		}
	})
}

func testStoreYShiftedRightInX(t *testing.T) {
	t.Helper()
	emu := mkEmu(t)

	const (
		x = 1
		y = 2
	)

	t.Run("it will not affect VF if bit is not dropped", func(t *testing.T) {
		emu.V[y] = 0b00010000
		want := byte(0b00001000)

		if err := emu.storeYShiftedRightInX(x, y); err != nil {
			t.Fatalf("error: %v", err)
		}

		if emu.V[x] != want {
			t.Fatalf("got %#0x, want %#0x", emu.V[x], want)
		}

		if emu.V[0xF] != 0 {
			t.Fatalf("VF was set")
		}
	})

	t.Run("it will affect VF if bit is dropped", func(t *testing.T) {
		emu.V[y] = 0b00000111
		want := byte(0b00000011)

		if err := emu.storeYShiftedRightInX(x, y); err != nil {
			t.Fatalf("error: %v", err)
		}

		if emu.V[x] != want {
			t.Fatalf("got %#0x, want %#0x", emu.V[x], want)
		}

		if emu.V[0xF] != 1 {
			t.Fatalf("VF was not set")
		}
	})
}

func testSetXToYMinusX(t *testing.T) {
	t.Helper()
	emu := mkEmu(t)

	const (
		x    = 1
		y    = 2
		valX = 1
		valY = 10
	)

	emu.V[x] = valX
	emu.V[y] = valY

	t.Run("it subs VX from VY without borrow", func(t *testing.T) {
		want := valY - valX

		if err := emu.setXToYMinusX(x, y); err != nil {
			t.Fatalf("error: %v", err)
		}

		if emu.V[x] != byte(want) {
			t.Fatalf("got %#0x, want %#0x", emu.V[x], want)
		}

		if emu.V[0xF] != 1 {
			t.Fatalf("borrow flag was cleared")
		}
	})

	t.Run("it subs VX from VY with borrow", func(t *testing.T) {
		emu.V[x] = 0xFF
		emu.V[y] = 0x0F

		want := 0xF0

		if err := emu.setXToYMinusX(x, y); err != nil {
			t.Fatalf("error: %v", err)
		}

		if emu.V[x] != byte(want) {
			t.Fatalf("got %#0x, want %#0x", emu.V[x], want)
		}

		if emu.V[0xF] != 0 {
			t.Fatalf("borrow flag was not cleared")
		}
	})
}

func testStoreYShiftedLeftInX(t *testing.T) {
	t.Helper()
	emu := mkEmu(t)

	const (
		x = 1
		y = 2
	)

	t.Run("it will not affect VF if bit is not dropped", func(t *testing.T) {
		emu.V[y] = 0b00010000
		want := byte(0b00100000)

		if err := emu.storeYShiftedLeftInX(x, y); err != nil {
			t.Fatalf("error: %v", err)
		}

		if emu.V[x] != want {
			t.Fatalf("got %#0x, want %#0x", emu.V[x], want)
		}

		if emu.V[0xF] != 0 {
			t.Fatalf("VF was set")
		}
	})

	t.Run("it will affect VF if bit is dropped", func(t *testing.T) {
		emu.V[y] = 0b10000111
		want := byte(0b00001110)

		if err := emu.storeYShiftedLeftInX(x, y); err != nil {
			t.Fatalf("error: %v", err)
		}

		if emu.V[x] != want {
			t.Fatalf("got %#0x, want %#0x", emu.V[x], want)
		}

		if emu.V[0xF] != 1 {
			t.Fatalf("VF was not set")
		}
	})
}

func testSkipIfXNotEqY(t *testing.T) {
	t.Helper()

	emu := mkEmu(t)

	const (
		x       = 1
		y       = 2
		testVal = 5
	)

	t.Run("it doesn't skip if Vx eq Vy", func(t *testing.T) {
		emu.V[x] = testVal
		emu.V[y] = emu.V[x]
		wantPC := emu.PC

		if err := emu.skipIfXNotEqY(x, y); err != nil {
			t.Fatalf("error: %v", err)
		}

		if emu.PC != wantPC {
			t.Fatalf("got %#0x, want %#0x", emu.PC, wantPC)
		}
	})

	t.Run("it skips if Vx not eq Vy", func(t *testing.T) {
		emu.V[x] = emu.V[y] + 1
		wantPC := emu.PC + InstructionSize

		if err := emu.skipIfXNotEqY(x, y); err != nil {
			t.Fatalf("error: %v", err)
		}

		if emu.PC != wantPC {
			t.Fatalf("got %#0x, want %#0x", emu.PC, wantPC)
		}
	})
}

func testSetVXWithMask(t *testing.T) {
	t.Helper()

	emu := mkEmu(t)
	// run 5 times, check if at least once the value was set
	const (
		iters = 10
		testV = 5
		testN = 0xF
	)

	for k := 0; k < iters; k++ {
		if err := emu.setXToRandomNumWithMaskNN(testV, []int{testN, testN}); err != nil {
			t.Fatalf("error: %v", err)
		}

		if emu.V[testV] != 0 {
			return
		}
	}

	t.Fatalf("register V%d was never set", testV)
}
