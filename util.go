package chipper

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func ToAddr3(p []int) (uint16, error) {
	n := len(p)
	if n != 3 {
		return 0, fmt.Errorf("expected 3 values, got %d", n)
	}

	addr := uint16(p[0])<<(8) | uint16(p[1])<<4 | uint16(p[2])

	return addr, nil
}

func toUint16(b []byte) uint16 {
	return (uint16(b[0]) << 8) | uint16(b[1])
}

func isInBounds(n, index int) error {
	minVal := 0
	maxVal := n - 1

	if index < minVal || index > maxVal {
		return fmt.Errorf("out of bounds: %d not in range %d - %d", index, minVal, maxVal)
	}

	return nil
}

func ToByte(p []int) (byte, error) {
	want := 2
	got := len(p)

	if got != want {
		return 0, fmt.Errorf("expected %d values, got %d", want, got)
	}

	val := byte(p[0])<<4 | byte(p[1])

	return val, nil
}

func randomNum() byte {
	return byte(rand.Intn(256))
}

func bcdOfInt(v int) ([]byte, error) {
	if v > 0xFF {
		return nil, fmt.Errorf("int (%d) exceeds max value %d", v, 0xFF)
	}

	const n = 3
	p := make([]byte, n)

	for k := 0; k < n; k++ {
		div := int(math.Pow(10, float64(n-k)))
		p[k] = byte(v / div)
		v = v - (int(p[k]) * div)
	}

	return p, nil

}

func DumpEmu(emu *Emulator) {
	p := make([]byte, 2)
	copy(p, emu.RAM[emu.PC:emu.PC+2])

	instr, _ := Decode(p)

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	fmt.Println("instruction: ", instr)
	fmt.Printf("PC: %0#4x\n", emu.PC)
	fmt.Printf("I: %0#4x\n", emu.Index)
	fmt.Println("Stack: ", emu.Stack.String())
	fmt.Println("->")
	b := &strings.Builder{}
	for k, v := range emu.V {
		fmt.Fprintf(
			b,
			"  (v%2d) %#0x\n",
			k, v,
		)
	}
	fmt.Print(b.String())
	fmt.Println(emu.Display.String())
}