package chipper

import "testing"

func TestStack(t *testing.T) {
	const (
		testStackSize = 2
		testValue     = uint16(5)
	)

	t.Run("push", func(tt *testing.T) {
		tt.Run("it should be able to push an element", func(ttt *testing.T) {
			s, _ := NewStack(testStackSize)
			wantVal := testValue
			wantPtr := 1

			if err := s.Push(wantVal); err != nil {
				ttt.Fatalf("error pushing: %v", err)
			}

			if s.pointer != wantPtr {
				ttt.Fatalf("(pointer) got %d, want %d", s.pointer, wantPtr)
			}

			if s.data[wantPtr-1] != wantVal {
				ttt.Fatalf("(value) got %d, want %d", s.data[wantPtr], wantVal)
			}
		})

		tt.Run("it should not be able push when full", func(ttt *testing.T) {
			s, _ := NewStack(2)

			if err := s.Push(1); err != nil {
				ttt.Fatalf("(1) error pushing: %v", err)
			}

			if err := s.Push(2); err != nil {
				ttt.Fatalf("(2) error pushing: %v", err)
			}

			if err := s.Push(3); err == nil {
				ttt.Fatalf("(3) should have errored, got nil")
			}
		})
	})

	t.Run("pop", func(tt *testing.T) {
		tt.Run("it should be able to pop an element", func(ttt *testing.T) {
			s, _ := NewStack(testStackSize)
			want := testValue

			if err := s.Push(want); err != nil {
				ttt.Fatalf("could not push element: %v", err)
			}

			got, err := s.Pop()
			if err != nil {
				ttt.Fatalf("could not pop element: %v", err)
			}

			if got != want {
				ttt.Fatalf("(value) got %d, want %d", got, want)
			}
		})

		tt.Run("it should not be able to pop an empty stack", func(ttt *testing.T) {
			s, _ := NewStack(testStackSize)
			if _, err := s.Pop(); err == nil {
				ttt.Fatalf("expected error, got nil")
			}
		})
	})
}
