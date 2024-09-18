package chipper

import "testing"

func TestStack(t *testing.T) { //nolint:funlen,gocognit
	const (
		testStackSize = 2
		testValue     = uint16(5)
	)

	t.Run("push", func(tt *testing.T) {
		tt.Run("it should be able to push an element", func(ttt *testing.T) {
			stack, err := NewStack(testStackSize)
			if err != nil {
				ttt.Fatalf("could not make new stack: %v", err)
			}

			wantVal := testValue
			wantPtr := 1

			if err := stack.Push(wantVal); err != nil {
				ttt.Fatalf("error pushing: %v", err)
			}

			if stack.pointer != wantPtr {
				ttt.Fatalf("(pointer) got %d, want %d", stack.pointer, wantPtr)
			}

			if stack.data[wantPtr-1] != wantVal {
				ttt.Fatalf("(value) got %d, want %d", stack.data[wantPtr], wantVal)
			}
		})

		tt.Run("it should not be able push when full", func(ttt *testing.T) {
			stack, err := NewStack(2)
			if err != nil {
				ttt.Fatalf("could not make new stack: %v", err)
			}

			if err := stack.Push(1); err != nil {
				ttt.Fatalf("(1) error pushing: %v", err)
			}

			if err := stack.Push(2); err != nil {
				ttt.Fatalf("(2) error pushing: %v", err)
			}

			if err := stack.Push(3); err == nil {
				ttt.Fatalf("(3) should have errored, got nil")
			}
		})
	})

	t.Run("pop", func(tt *testing.T) {
		tt.Run("it should be able to pop an element", func(ttt *testing.T) {
			stack, err := NewStack(testStackSize)
			if err != nil {
				ttt.Fatalf("could not make new stack: %v", err)
			}

			want := testValue

			if pushErr := stack.Push(want); pushErr != nil {
				ttt.Fatalf("could not push element: %v", pushErr)
			}

			got, err := stack.Pop()
			if err != nil {
				ttt.Fatalf("could not pop element: %v", err)
			}

			if got != want {
				ttt.Fatalf("(value) got %d, want %d", got, want)
			}
		})

		tt.Run("it should not be able to pop an empty stack", func(ttt *testing.T) {
			s, err := NewStack(testStackSize)
			if err != nil {
				ttt.Fatalf("could not make new stack: %v", err)
			}

			if _, err := s.Pop(); err == nil {
				ttt.Fatalf("expected error, got nil")
			}
		})
	})
}
