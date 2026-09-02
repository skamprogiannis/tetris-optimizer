package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	execute(os.Args[1:], os.Stdout)
}

func execute(args []string, output io.Writer) {
	if len(args) != 1 {
		fmt.Fprintln(output, "ERROR")
		return
	}

	data, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Fprintln(output, "ERROR")
		return
	}

	pieces, err := parseTetrominoes(data)
	if err != nil {
		fmt.Fprintln(output, "ERROR")
		return
	}

	board, size := solveTetrominoes(pieces)
	fmt.Fprint(output, formatBoard(board, size))
}
