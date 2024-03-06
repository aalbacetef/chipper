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

	if err := run(data, dump, text); err != nil {
		log.Println("error: ", err)

		return
	}
}

func run(data []byte, dump bool, text bool) error {
	if dump {
		fmt.Print(printHeader("dump"))

		r := bytes.NewReader(data)

		if err := dumpInstructions(r); err != nil && !errors.Is(err, io.EOF) {
			return err
		}

		fmt.Print("\n\n")
	}

	if text {
		fmt.Print(printHeader("text"))

		r := bytes.NewReader(data)

		if err := toHumanReadable(r); err != nil {
			return err
		}
	}

	return nil
}

// NOTE: assumes text is 4 characters wide.
func printHeader(text string) string {
	const (
		hrRow  = "------------------"
		pad    = "|                |"
		center = "|      %s      |"
	)

	return strings.Join(
		[]string{
			hrRow,
			pad,
			fmt.Sprintf(center, text),
			pad,
			hrRow,
		},
		"\n",
	) + "\n"
}
