package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"gitlab.com/aalbacetef/chipper"
	"golang.org/x/exp/slices"
)

func main() {
	name := ""
	dump := false
	text := false

	flag.StringVar(&name, "name", name, "filepath to read")
	flag.BoolVar(&dump, "dump", dump, "dump instructions")
	flag.BoolVar(&text, "text", text, "to human readable text")

	flag.Parse()

	if !(dump || text) {
		flag.Usage()
		return
	}

	if name == "" {
		flag.Usage()
		return
	}

	data, err := os.ReadFile(name)
	if err != nil {
		log.Println("error reading file: ", err)
		return
	}

	if dump {
		fmt.Println("" +
			"------------------\n" +
			"|                |\n" +
			"|      dump      |\n" +
			"|                |\n" +
			"------------------\n")
		r := bytes.NewReader(data)

		if err := dumpInstructions(r); err != nil && !errors.Is(err, io.EOF) {
			fmt.Println("error: ", err)
			return
		}

		fmt.Println("\n\n")
	}

	if text {
		fmt.Println("" +
			"------------------\n" +
			"|                |\n" +
			"|      text      |\n" +
			"|                |\n" +
			"------------------\n")
		r := bytes.NewReader(data)

		if err := runThrough(r); err != nil {
			fmt.Println("error: ", err)
			return
		}
	}
}

func dumpInstructions(r io.Reader) error {
	p := make([]byte, 2)

	instructions := make(map[chipper.Opcode]int)
	b := &strings.Builder{}
	tw := tabwriter.NewWriter(
		b, 0, 0, 1, ' ', tabwriter.TabIndent)

	v := func() error {
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
				continue
			}

			instructions[instr.Op]++
		}
	}

	if err := v(); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("error ocurred: %w", err)
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

func runThrough(r io.Reader) error {
	p := make([]byte, 2)
	k := 0
	b := &strings.Builder{}
	tw := tabwriter.NewWriter(
		b, 0, 0, 1, ' ', tabwriter.TabIndent)

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

		fmt.Fprintf(
			tw,
			"%0#4x) %s \t=> %0#3x \t| %+v\n",
			k+0x200, instr.Op,
			addr,
			instr.Operands,
		)
		k += bytesRead
	}

}
