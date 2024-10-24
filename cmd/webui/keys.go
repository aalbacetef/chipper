//go:build js && wasm

package main

import (
	"fmt"
	"sync"

	"github.com/aalbacetef/chipper"
)

type WebKeyInputSource struct {
	keys     [16]bool
	mu       sync.Mutex
	listener chan int
}

func NewWebKeyInputSource() *WebKeyInputSource {
	return &WebKeyInputSource{}
}

func (ksrc *WebKeyInputSource) Set(key int, isPressed bool) {
	ksrc.mu.Lock()
	defer ksrc.mu.Unlock()

	if key < 0 || key > chipper.NumKeys {
		fmt.Printf(
			"key %d out of bounds [%d, %d]\n",
			key,
			0, chipper.NumKeys,
		)

		return
	}

	if isPressed && ksrc.listener != nil {
		fmt.Println("set key listener")
		ksrc.listener <- key
		ksrc.listener = nil
	}

	ksrc.keys[key] = isPressed
}

func (ksrc *WebKeyInputSource) Get(key int) bool {
	ksrc.mu.Lock()
	defer ksrc.mu.Unlock()

	if key < 0 || key > chipper.NumKeys {
		fmt.Printf(
			"key %d out of bounds [%d, %d]\n",
			key,
			0, chipper.NumKeys,
		)

		return false
	}

	v := ksrc.keys[key]
	return v
}

func (ksrc *WebKeyInputSource) WaitUntilKeypress() <-chan int {
	l := make(chan int, 1)

	ksrc.mu.Lock()
	ksrc.listener = l
	ksrc.mu.Unlock()

	return l
}
