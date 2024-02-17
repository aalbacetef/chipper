package chipper

import "testing"

func TestUtil(t *testing.T) {
	v, err := ToAddr3([]int{0x1, 0x2, 0x3})
	t.Logf("v: %0#x", v)
	t.Logf("err: %v", err)
	if v != 0x123 {
		t.Fatalf("failed")
	}
}
