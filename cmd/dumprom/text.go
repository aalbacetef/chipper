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
	b := &strings.Builder{}
	tw := tabwriter.NewWriter(
		b,
		0, 0, 1, ' ',
		tabwriter.TabIndent,
	)

	err := readUntilError(r, tw)
	tw.Flush()
	fmt.Println(b.String())

	return err
}

func readUntilError(r io.Reader, w io.Writer) error {
	const width = 2
	p := make([]byte, width)
	cnt := 0

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
				w,
				"%0#4x) could not decode (%v) \t=> %0#4x\n",
				cnt,
				err,
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
			w,
			"%0#4x) %s \t=> %0#3x \t| %+v\n",
			cnt+offset, instr.Op,
			addr,
			instr.Operands,
		)

		cnt += bytesRead
	}
}
