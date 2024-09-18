package chipper

import "testing"

func TestUtil(t *testing.T) {
	got, err := ToAddr3([]int{0x1, 0x2, 0x3})
	if err != nil {
		t.Fatalf("got error: %v", err)
	}

	const want = 0x123
	if got != want {
		t.Fatalf("want: %#0x, got %#0x", want, got)
	}
}
