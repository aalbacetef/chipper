package main

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/aalbacetef/chipper"
)

func toHumanReadable(r io.Reader) error {
	const instructionWidth = 2
	p := make([]byte, instructionWidth)

	k := 0
	b := &strings.Builder{}
	tw := tabwriter.NewWriter(
		b,
		0, 0, 1, ' ',
		tabwriter.TabIndent,
	)

	defer func() {
		tw.Flush()
		fmt.Println(b.String())
	}()

	for {
		bytesRead, err := r.Read(p)
		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return fmt.Errorf("error reading: %w", err)
		}

		if bytesRead == 0 {
			return fmt.Errorf("read 0 bytes")
		}

		instr, err := chipper.Decode(p)
		if err != nil {
			fmt.Fprintf(
				tw,
				"%0#4x) could not decode \t=> %0#4x\n",
				k,
				uint16(p[0])<<8|uint16(p[1]),
			)

			continue
		}

		addr, err := chipper.ToAddr3(instr.Operands)
		if err != nil {
			return fmt.Errorf("could not parse addr3: %w", err)
		}

		const offset = 0x200

		fmt.Fprintf(
			tw,
			"%0#4x) %s \t=> %0#3x \t| %+v\n",
			k+offset, instr.Op,
			addr,
			instr.Operands,
		)

		k += bytesRead
	}
}
