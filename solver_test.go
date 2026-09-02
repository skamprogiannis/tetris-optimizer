package main

import (
	"strings"
	"testing"
)

func TestSolveTetrominoesFindsSmallestSquare(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantSize   int
		wantSpaces int
	}{
		{
			name:       "one square",
			input:      "....\n.##.\n.##.\n....\n",
			wantSize:   2,
			wantSpaces: 0,
		},
		{
			name: "official good example 01",
			input: "...#\n...#\n...#\n...#\n\n" +
				"....\n....\n....\n####\n\n" +
				".###\n...#\n....\n....\n\n" +
				"....\n..##\n.##.\n....\n",
			wantSize:   5,
			wantSpaces: 9,
		},
		{
			name: "official good example 02",
			input: "...#\n...#\n...#\n...#\n\n" +
				"....\n....\n....\n####\n\n" +
				".###\n...#\n....\n....\n\n" +
				"....\n..##\n.##.\n....\n\n" +
				"....\n.##.\n.##.\n....\n\n" +
				"....\n....\n##..\n.##.\n\n" +
				"##..\n.#..\n.#..\n....\n\n" +
				"....\n###.\n.#..\n....\n",
			wantSize:   6,
			wantSpaces: 4,
		},
		{
			name: "official good example 03",
			input: "....\n.##.\n.##.\n....\n\n" +
				"...#\n...#\n...#\n...#\n\n" +
				"....\n..##\n.##.\n....\n\n" +
				"....\n.##.\n.##.\n....\n\n" +
				"....\n..#.\n.##.\n.#..\n\n" +
				".###\n...#\n....\n....\n\n" +
				"##..\n.#..\n.#..\n....\n\n" +
				"....\n..##\n.##.\n....\n\n" +
				"##..\n.#..\n.#..\n....\n\n" +
				".#..\n.##.\n..#.\n....\n\n" +
				"....\n###.\n.#..\n....\n",
			wantSize:   7,
			wantSpaces: 5,
		},
		{
			name: "official hard example",
			input: "....\n.##.\n.##.\n....\n\n" +
				".#..\n.##.\n.#..\n....\n\n" +
				"....\n..##\n.##.\n....\n\n" +
				"....\n.##.\n.##.\n....\n\n" +
				"....\n..#.\n.##.\n.#..\n\n" +
				".###\n...#\n....\n....\n\n" +
				"##..\n.#..\n.#..\n....\n\n" +
				"....\n.##.\n.##.\n....\n\n" +
				"....\n..##\n.##.\n....\n\n" +
				"##..\n.#..\n.#..\n....\n\n" +
				".#..\n.##.\n..#.\n....\n\n" +
				"....\n###.\n.#..\n....\n",
			wantSize:   7,
			wantSpaces: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pieces, err := parseTetrominoes([]byte(test.input))
			if err != nil {
				t.Fatalf("parseTetrominoes() error = %v", err)
			}

			board, size := solveTetrominoes(pieces)
			if size != test.wantSize {
				t.Fatalf("solveTetrominoes() size = %d, want %d", size, test.wantSize)
			}
			if got := strings.Count(string(board), "."); got != test.wantSpaces {
				t.Fatalf("solveTetrominoes() empty cells = %d, want %d\n%s", got, test.wantSpaces, formatBoard(board, size))
			}

			assertValidSolution(t, board, size, pieces)
		})
	}
}

func assertValidSolution(t *testing.T, board []byte, size int, pieces []tetromino) {
	t.Helper()

	if len(board) != size*size {
		t.Fatalf("board has %d cells, want %d", len(board), size*size)
	}

	for _, piece := range pieces {
		if got := strings.Count(string(board), string(piece.letter)); got != 4 {
			t.Fatalf("piece %q occupies %d cells, want 4", piece.letter, got)
		}
	}
}

func TestFormatBoard(t *testing.T) {
	board := []byte("ABB." + "ABBA" + "A..." + "....")
	const want = "ABB.\nABBA\nA...\n....\n"

	if got := formatBoard(board, 4); got != want {
		t.Fatalf("formatBoard() = %q, want %q", got, want)
	}
}
