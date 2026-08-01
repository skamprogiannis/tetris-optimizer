package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <input.txt>\n", os.Args[0])
		os.Exit(1)
	}

	path := os.Args[1]
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read %q: %v\n", path, err)
		os.Exit(1)
	}

	if _, err := parseTetrominoes(data); err != nil {
		fmt.Fprintf(os.Stderr, "error: parse %q: %v\n", path, err)
		os.Exit(1)
	}
}
