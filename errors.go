package chipper

import "fmt"

type ArgCountError struct {
	got  int
	want int
}

func (e ArgCountError) Error() string {
	return fmt.Sprintf("want %d args, got %d", e.want, e.got)
}
