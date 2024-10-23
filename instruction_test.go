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

	d, err := NewDebugDisplay(testW, testH)
	if err != nil {
		t.Fatalf("could not make debug display: %v", err)
	}

	emu, err := NewEmulator(
		testStackSize,
		testRAMSize,
		d,
		&StubKeyInputSource{},
	)
	if err != nil {
		t.Fatalf("could not create emulator: %v", err)
	}

	return emu
}

func TestInstruction(t *testing.T) {
	t.Run("clearScreen", testClearScreen)
	t.Run("returnFromSub", testReturnFromSub)
	t.Run("jumpNNN", testJumpNNN)
	t.Run("callSubNNN", testCallSubNNN)
	t.Run("add NN to X", testAddNNToX)
	t.Run("set vx to random number mask with nn", testSetVXWithMask)
}

func testClearScreen(t *testing.T) {
	emu := mkEmu(t)

	colorSet := emu.Display.ColorSet()
	_ = Each(emu.Display, func(x, y int) error {
		emu.Display.Set(x, y, colorSet)

		return nil
	})

	if err := emu.clearScreen(); err != nil {
		t.Fatalf("error: %v", err)
	}

	colorClear := emu.Display.ColorClear()
	_ = Each(emu.Display, func(x, y int) error {
		at := emu.Display.At(x, y)
		if !ColorEq(at, colorClear) {
			t.Fatalf("(%d, %d) not clear", x, y)
		}

		return nil
	})
}

func testReturnFromSub(t *testing.T) {
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

func testAddNNToX(t *testing.T) {
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

func testSetVXWithMask(t *testing.T) {
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
