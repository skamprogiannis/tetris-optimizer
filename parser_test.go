package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseTetrominoesNormalizesValidInput(t *testing.T) {
	input := []byte("....\n..#.\n..#.\n..##\n")
	want := []tetromino{
		{
			blocks: [4]point{{x: 0, y: 0}, {x: 0, y: 1}, {x: 0, y: 2}, {x: 1, y: 2}},
			width:  2,
			height: 3,
			letter: 'A',
		},
	}

	got, err := parseTetrominoes(input)
	if err != nil {
		t.Fatalf("parseTetrominoes() error = %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseTetrominoes() = %#v, want %#v", got, want)
	}
}

func TestParseTetrominoesRequiresFinalNewline(t *testing.T) {
	input := []byte("####\n....\n....\n....")

	_, err := parseTetrominoes(input)
	if err == nil {
		t.Fatal("parseTetrominoes() error = nil, want final-newline error")
	}
}

func TestParseTetrominoesRejectsDisconnectedBlocks(t *testing.T) {
	input := []byte("#..#\n....\n....\n#..#\n")

	_, err := parseTetrominoes(input)
	if err == nil {
		t.Fatal("parseTetrominoes() error = nil, want disconnected-blocks error")
	}

	const want = "piece 1: blocks are not edge-connected"
	if err.Error() != want {
		t.Fatalf("parseTetrominoes() error = %q, want %q", err, want)
	}
}

func TestParseTetrominoesRejectsMoreThanTwentySixPieces(t *testing.T) {
	const section = "####\n....\n....\n...."
	input := []byte(strings.Repeat(section+"\n\n", 26) + section + "\n")

	_, err := parseTetrominoes(input)
	if err == nil {
		t.Fatal("parseTetrominoes() error = nil, want piece-limit error")
	}

	const want = "input contains 27 pieces; maximum is 26"
	if err.Error() != want {
		t.Fatalf("parseTetrominoes() error = %q, want %q", err, want)
	}
}

func TestParseTetrominoesReportsExactBlockCount(t *testing.T) {
	input := []byte("####\n#...\n....\n....\n")

	_, err := parseTetrominoes(input)
	if err == nil {
		t.Fatal("parseTetrominoes() error = nil, want block-count error")
	}

	const want = "piece 1: expected 4 blocks, got 5"
	if err.Error() != want {
		t.Fatalf("parseTetrominoes() error = %q, want %q", err, want)
	}
}

func TestParseTetrominoesLabelsPiecesInInputOrder(t *testing.T) {
	input := []byte("####\n....\n....\n....\n\n.##.\n.##.\n....\n....\n")

	got, err := parseTetrominoes(input)
	if err != nil {
		t.Fatalf("parseTetrominoes() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(parseTetrominoes()) = %d, want 2", len(got))
	}
	if got[0].letter != 'A' || got[1].letter != 'B' {
		t.Fatalf("piece letters = %q, %q; want A, B", got[0].letter, got[1].letter)
	}
	if got[1].width != 2 || got[1].height != 2 {
		t.Fatalf("second piece dimensions = %dx%d, want 2x2", got[1].width, got[1].height)
	}
}

func TestParseTetrominoesAcceptsTwentySixPieces(t *testing.T) {
	const section = "####\n....\n....\n...."
	input := []byte(strings.Repeat(section+"\n\n", 25) + section + "\n")

	got, err := parseTetrominoes(input)
	if err != nil {
		t.Fatalf("parseTetrominoes() error = %v", err)
	}
	if len(got) != 26 {
		t.Fatalf("len(parseTetrominoes()) = %d, want 26", len(got))
	}
	if got[25].letter != 'Z' {
		t.Fatalf("last piece letter = %q, want Z", got[25].letter)
	}
}

func TestParseTetrominoesAcceptsCanonicalShapes(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		width  int
		height int
	}{
		{name: "I", input: "####\n....\n....\n....\n", width: 4, height: 1},
		{name: "O", input: ".##.\n.##.\n....\n....\n", width: 2, height: 2},
		{name: "T", input: "###.\n.#..\n....\n....\n", width: 3, height: 2},
		{name: "L", input: "#...\n#...\n##..\n....\n", width: 2, height: 3},
		{name: "J", input: ".#..\n.#..\n##..\n....\n", width: 2, height: 3},
		{name: "S", input: ".##.\n##..\n....\n....\n", width: 3, height: 2},
		{name: "Z", input: "##..\n.##.\n....\n....\n", width: 3, height: 2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseTetrominoes([]byte(test.input))
			if err != nil {
				t.Fatalf("parseTetrominoes() error = %v", err)
			}
			if len(got) != 1 {
				t.Fatalf("len(parseTetrominoes()) = %d, want 1", len(got))
			}
			if got[0].width != test.width || got[0].height != test.height {
				t.Fatalf("piece dimensions = %dx%d, want %dx%d", got[0].width, got[0].height, test.width, test.height)
			}
		})
	}
}

func TestParseTetrominoesRejectsMalformedInput(t *testing.T) {
	const section = "####\n....\n....\n...."
	tests := []struct {
		name      string
		input     string
		wantError string
	}{
		{name: "empty", wantError: "input is empty"},
		{name: "empty line", input: "\n", wantError: "input is empty"},
		{name: "CRLF", input: "####\r\n....\r\n....\r\n....\r\n", wantError: "piece 1: row 1: expected 4 cells, got 5"},
		{name: "missing separator", input: section + "\n" + section + "\n", wantError: "piece 1: expected 4 rows, got 8"},
		{name: "extra separator", input: section + "\n\n\n" + section + "\n", wantError: "piece 2: expected 4 rows, got 5"},
		{name: "trailing blank line", input: section + "\n\n", wantError: "piece 1: expected 4 rows, got 5"},
		{name: "three rows", input: "####\n....\n....\n", wantError: "piece 1: expected 4 rows, got 3"},
		{name: "short row", input: "###\n....\n....\n#...\n", wantError: "piece 1: row 1: expected 4 cells, got 3"},
		{name: "invalid character", input: "###x\n....\n....\n#...\n", wantError: "piece 1: row 1, column 4: invalid character 'x'"},
		{name: "too few blocks", input: "###.\n....\n....\n....\n", wantError: "piece 1: expected 4 blocks, got 3"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := parseTetrominoes([]byte(test.input))
			if err == nil {
				t.Fatal("parseTetrominoes() error = nil, want malformed-input error")
			}
			if err.Error() != test.wantError {
				t.Fatalf("parseTetrominoes() error = %q, want %q", err, test.wantError)
			}
		})
	}
}
