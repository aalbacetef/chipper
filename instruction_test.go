package chipper

import "testing"

func mkEmu(t *testing.T) *Emulator {
	t.Helper()
	const (
		testStackSize = 16
		testRamSize   = 4096
		testW         = 64
		testH         = 32
	)

	emu, err := NewEmulator(testStackSize, testRamSize, testW, testH)
	if err != nil {
		t.Fatalf("could not create emulator: %v", err)
	}

	return emu
}

func TestInstruction(t *testing.T) {
	t.Run("clearScreen", func(tt *testing.T) {
		emu := mkEmu(tt)
		for k := range emu.Display.data {
			emu.Display.data[k] = 1
		}

		if err := emu.clearScreen(); err != nil {
			tt.Fatalf("error: %v", err)
		}

		for k, v := range emu.Display.data {
			if v != 0 {
				tt.Fatalf("index %d is not cleared", k)
			}
		}
	})

	t.Run("returnFromSub", func(tt *testing.T) {
		emu := mkEmu(tt)
		const testAddr = 0x111
		if err := emu.Stack.Push(testAddr); err != nil {
			tt.Fatalf("could not push: %v", err)
		}

		if err := emu.returnFromSub(); err != nil {
			tt.Fatalf("error: %v", err)
		}

		if emu.PC != testAddr {
			tt.Fatalf("got %#0x, want %#0x", emu.PC, testAddr)
		}
	})

	t.Run("jumpNNN", func(tt *testing.T) {
		emu := mkEmu(tt)
		const testAddr = 0x111
		if err := emu.jumpNNN([]int{1, 1, 1}); err != nil {
			tt.Fatalf("error: %v", err)
		}

		if emu.PC != testAddr {
			tt.Fatalf("got %#0x, want %#0x", emu.PC, testAddr)
		}
	})

	t.Run("callSubNNN", func(tt *testing.T) {
		emu := mkEmu(tt)
		const testAddr = 0x111
		const origPC = 0x210

		emu.PC = origPC

		if err := emu.callSubNNN([]int{1, 1, 1}); err != nil {
			tt.Fatalf("error: %v", err)
		}

		if emu.PC != testAddr {
			tt.Fatalf("got %#0x, want %#0x", emu.PC, testAddr)
		}

		val, err := emu.Stack.Pop()
		if err != nil {
			tt.Fatalf("could not pop stack: %v", err)
		}

		if val != origPC {
			tt.Fatalf("got %#0x, want %#0x", val, origPC)
		}
	})

	t.Run("add NN to X", func(tt *testing.T) {
		emu := mkEmu(tt)
		const (
			testX = 5
			testN = 0x1
		)
		if err := emu.addNNToX(testX, []int{testN, testN}); err != nil {
			tt.Fatalf("error: %v", err)
		}

		got := emu.V[testX]
		want, err := ToByte([]int{testN, testN})
		if err != nil {
			tt.Fatalf("could not convert to byte: %v", err)
		}

		if got != want {
			tt.Fatalf("got %d, want %d", got, want)
		}
	})

	t.Run("set vx to random number mask with nn", func(tt *testing.T) {
		emu := mkEmu(tt)
		// run 5 times, check if at least once the value was set
		const iters = 10
		const testV = 5
		const testN = 0xF
		for k := 0; k < iters; k++ {
			if err := emu.setXToRandomNumWithMaskNN(testV, []int{testN, testN}); err != nil {
				tt.Fatalf("error: %v", err)
			}

			if emu.V[testV] != 0 {
				return
			}
		}

		tt.Fatalf("register V%d was never set", testV)
	})
}
