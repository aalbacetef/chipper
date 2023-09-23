package chipper

import (
	"fmt"
	"math"
	"math/rand"
)

func ToAddr3(p []int) (uint16, error) {
	n := len(p)
	if n != 3 {
		return 0, fmt.Errorf("expected 3 values, got %d", n)
	}

	addr := uint16(p[0])<<(8) | uint16(p[1])<<4 | uint16(p[0])

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

	n := 3
	p := make([]byte, 3)

	for k := 0; k < n; k++ {
		div := int(math.Pow(10, float64(n-k)))
		p[k] = byte(v / div)
		v = v - (int(p[k]) * div)
	}

	return p, nil

}
