package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"gitlab.com/aalbacetef/chipper"
)

func main() {
	fname := os.Args[1]
	data, err := os.ReadFile(fname)
	if err != nil {
		log.Println("error reading file: ", err)
		return
	}

	p := make([]byte, 2)
	r := bytes.NewReader(data)
	k := 0
	b := &strings.Builder{}
	tw := tabwriter.NewWriter(
		b, 0, 0, 1, ' ', tabwriter.TabIndent)

	for {
		bytesRead, err := r.Read(p)
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			log.Println("error reading: ", err)
			break
		}

		if bytesRead == 0 {
			log.Println("read 0 bytes")
			break
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

		fmt.Fprintf(
			tw,
			"%0#4x) %s \t=> %0#3x \t| %+v\n",
			k+0x200, instr.Op, chipper.ToAddr3(instr.Operands), instr.Operands,
		)
		k += bytesRead
	}

	tw.Flush()
	fmt.Println(b.String())
}
