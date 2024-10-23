package chipper

import (
	"fmt"
	"strings"
)

// NewStack will return a stack of size N, with the data already allocated.
func NewStack(size int) (*Stack, error) {
	if size <= 0 {
		return nil, fmt.Errorf("stack size must be > 0, got %d", size)
	}

	return &Stack{
		data:    make([]uint16, size),
		size:    size,
		pointer: 0,
	}, nil
}

// Stack is LIFO data structure, it provides a Push/Pop API.
// @TODO: consider adding support for ShiftLeft and ShiftRight.
type Stack struct {
	size    int
	data    []uint16
	pointer int
}

func (s Stack) String() string {
	if s.pointer == 0 {
		return "Stack: <empty stack>"
	}

	b := &strings.Builder{}
	for k, v := range s.data[:s.pointer] {
		fmt.Fprintf(b, "%2d) %0#4x\n", k, v)
	}

	return b.String()
}

// Pop an element of the stack, will error if the stack is empty.
func (s *Stack) Pop() (uint16, error) {
	if s.pointer == 0 {
		return 0, fmt.Errorf("stack is empty")
	}

	s.pointer -= 1
	val := s.data[s.pointer]

	return val, nil
}

// Push an element onto the stack, will error if the stack is full.
func (s *Stack) Push(val uint16) error {
	if s.pointer == s.size {
		return fmt.Errorf("stack is full (size=%d)", s.size)
	}

	s.data[s.pointer] = val
	s.pointer += 1

	return nil
}
