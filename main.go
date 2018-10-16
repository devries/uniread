package main

import (
	"bufio"
	"fmt"
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
	decoder := unicode.BOMOverride(unicode.UTF8.NewDecoder())

	f, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %s\n", err)
		os.Exit(2)
	}
	defer f.Close()

	r := transform.NewReader(f, decoder)

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		fmt.Println(sc.Bytes())
		fmt.Println(sc.Text())
	}
	if err = sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
		os.Exit(3)
	}
}
