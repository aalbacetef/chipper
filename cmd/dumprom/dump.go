package main

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"gitlab.com/aalbacetef/chipper"
	"golang.org/x/exp/slices"
)

func dumpInstructions(r io.Reader) error {
	const instructionWidth = 2
	p := make([]byte, instructionWidth)

	b := &strings.Builder{}
	tw := tabwriter.NewWriter(
		b, 0, 0, 1, ' ', tabwriter.TabIndent)

	instructions, err := countInstructions(p, r)
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("error occurred: %w", err)
	}

	keys := make([]string, 0, len(instructions))
	for key := range instructions {
		keys = append(keys, string(key))
	}

	slices.Sort(keys)

	for k, key := range keys {
		fmt.Fprintf(tw, "%2d) %s \t(count: %3d)\n", k, key, instructions[chipper.Opcode(key)])
	}

	tw.Flush()
	fmt.Println(b.String())

	return nil
}

func countInstructions(p []byte, r io.Reader) (map[chipper.Opcode]int, error) {
	instructions := make(map[chipper.Opcode]int)

	for {
		bytesRead, err := r.Read(p)
		if errors.Is(err, io.EOF) {
			return instructions, nil
		}

		if err != nil {
			return nil, fmt.Errorf("error reading: %w", err)
		}

		if bytesRead == 0 {
			return nil, fmt.Errorf("read 0 bytes")
		}

		instr, err := chipper.Decode(p)
		if err != nil {
			continue
		}

		instructions[instr.Op]++
	}
}
