package chipper

import "testing"

func TestUtil(t *testing.T) {
	t.Run("ToAddr3", func(t *testing.T) {
		got, err := ToAddr3([]int{0x1, 0x2, 0x3})
		if err != nil {
			t.Fatalf("got error: %v", err)
		}

		const want = 0x123
		if got != want {
			t.Fatalf("want: %#0x, got %#0x", want, got)
		}
	})

	t.Run("ToByte", func(t *testing.T) {
		cases := []struct {
			label string
			args  []int
			want  byte
		}{
			{"0x23", []int{0x2, 0x3}, 0x23},
			{"0x11", []int{0x1, 0x1}, 0x11},
		}

		for _, c := range cases {
			t.Run(c.label, func(t *testing.T) {
				got, err := ToByte(c.args)
				if err != nil {
					t.Fatalf("got error: %v", err)
				}

				if got != c.want {
					t.Fatalf("want #%0x, got #%0x", c.want, got)
				}
			})
		}
	})
}
