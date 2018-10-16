package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"unicode"

	xunicode "golang.org/x/text/encoding/unicode"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
)

func main() {
	// The samples below are all the string "Hello World" represented in different unicode encodings.
	samples := [][]byte{
		[]byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64},                                                                               // UTF-8
		[]byte{0xef, 0xbb, 0xbf, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64},                                                             // UTF-8 w/ BOM
		[]byte{0xff, 0xfe, 0x48, 0x00, 0x65, 0x00, 0x6c, 0x00, 0x6c, 0x00, 0x6f, 0x00, 0x20, 0x00, 0x57, 0x00, 0x6f, 0x00, 0x72, 0x00, 0x6c, 0x00, 0x64, 0x00}, // UTF-16-LE
		[]byte{0xfe, 0xff, 0x00, 0x48, 0x00, 0x65, 0x00, 0x6c, 0x00, 0x6c, 0x00, 0x6f, 0x00, 0x20, 0x00, 0x57, 0x00, 0x6f, 0x00, 0x72, 0x00, 0x6c, 0x00, 0x64}, // UTF-16-BE
		[]byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0xe4, 0xb8, 0x96, 0xe7, 0x95, 0x8c},                                                                         // UTF-8
	}

	// The decoder will override to the correct unicode decoder if it detects a byte order mark, but will default to the
	// fallback decoder specified (unicode.UTF8).
	decoder := xunicode.BOMOverride(xunicode.UTF8.NewDecoder())
	asciiEncoder := runes.Map(asciiEncoderMapping)
	graphicEncoder := runes.Map(graphicEncoderMapping)

	// This is a bit of a strange writer. It replaces all unprintable
	// output with unicode \ufffd characters.
	graphicStdout := transform.NewWriter(os.Stdout, graphicEncoder)

	for _, sample := range samples {
		f := bytes.NewReader(sample)         // Wrap sample in a reader.
		r := transform.NewReader(f, decoder) // Use a transforming reader with the decoder specified above

		// Scan the result
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			fmt.Fprintln(graphicStdout, "Input:       ", sample)
			fmt.Fprintln(graphicStdout, "Output Bytes:", sc.Bytes())
			fmt.Fprintln(graphicStdout, "Output Text: ", sc.Text())

			// Let's transform this for ASCII output
			var bout bytes.Buffer
			wout := transform.NewWriter(&bout, asciiEncoder)
			wout.Write(sc.Bytes())
			fmt.Fprintln(graphicStdout, "ASCII Bytes: ", bout.Bytes())
			fmt.Fprintln(graphicStdout, "ASCII Text:  ", bout.String())
			fmt.Fprintln(graphicStdout)
		}
		if err := sc.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading sample: %s\n", err)
			os.Exit(1)
		}
	}
}

func asciiEncoderMapping(inrune rune) rune {
	if inrune <= 127 {
		return inrune
	} else {
		return rune(0x1a)
	}
}

func graphicEncoderMapping(inrune rune) rune {
	if unicode.IsGraphic(inrune) || unicode.IsSpace(inrune) {
		return inrune
	} else {
		return rune(0xfffd)
	}
}
