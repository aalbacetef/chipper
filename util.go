package chipper

import "fmt"

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
