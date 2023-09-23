package chipper

import (
	"fmt"
	"testing"
)

func TestOpcode(t *testing.T) {
	a := []int{0x1, 0xA, 0xE, 0xF}
	b := a[2:]
	fmt.Printf("b: %#0x, %#0x\n", b[0], b[1])

}
