package main

import (
	// "bufio"
	"fmt"
	"io"
	"os"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <filename>\n", os.Args[0])
		os.Exit(1)
	}

	filename := os.Args[1]

	var f io.ReadCloser
	var err error
	if filename == "-" {
		f = os.Stdin
	} else {
		f, err = os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file \"%s\": %s\n", filename, err)
			os.Exit(2)
		}
	}
	defer f.Close()

	r := ConvertUnicode(f)

	/* *** DIAGNOSTIC OUTPUT
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		fmt.Println(sc.Bytes())
		fmt.Println(sc.Text())
	}
	if err = sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
		os.Exit(3)
	}
	*/

	// Write output as utf-8 without BOM to stdout
	if _, err := io.Copy(os.Stdout, r); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(3)
	}
	fmt.Println()
}

func ConvertUnicode(r io.Reader) io.Reader {
	decoder := unicode.BOMOverride(unicode.UTF8.NewDecoder())

	return transform.NewReader(r, decoder)
}
