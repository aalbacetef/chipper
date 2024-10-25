package chipper

type Direction int

const (
	Up Direction = iota
	Down
)

type KeyInputSource interface {
	Get(key int) bool
	Set(key int, v bool)
	WaitUntilKeypress() <-chan int
}

type StubKeyInputSource struct{}

func (stub *StubKeyInputSource) Get(_ int) bool {
	return false
}

func (stub *StubKeyInputSource) Set(_ int, _ bool) {
}

func (stub *StubKeyInputSource) WaitUntilKeypress() <-chan int {
	l := make(chan int, 1)
	l <- 0

	return l
}
