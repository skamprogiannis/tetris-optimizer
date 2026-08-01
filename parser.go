package main

import (
	"fmt"
	"strings"
)

func parseTetrominoes(data []byte) ([]tetromino, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("input is empty")
	}
	if data[len(data)-1] != '\n' {
		return nil, fmt.Errorf("input must end with a newline")
	}

	content := string(data[:len(data)-1])
	if content == "" {
		return nil, fmt.Errorf("input is empty")
	}

	sections := strings.Split(content, "\n\n")
	if len(sections) > 26 {
		return nil, fmt.Errorf("input contains %d pieces; maximum is 26", len(sections))
	}

	tetrominoes := make([]tetromino, 0, len(sections))

	for index, section := range sections {
		piece, err := parseTetromino(section, byte('A'+index))
		if err != nil {
			return nil, fmt.Errorf("piece %d: %w", index+1, err)
		}

		tetrominoes = append(tetrominoes, piece)
	}

	return tetrominoes, nil
}

func parseTetromino(section string, letter byte) (tetromino, error) {
	lines := strings.Split(section, "\n")
	if len(lines) != 4 {
		return tetromino{}, fmt.Errorf("expected 4 rows, got %d", len(lines))
	}

	var blocks [4]point
	blockCount := 0

	for y, line := range lines {
		if len(line) != 4 {
			return tetromino{}, fmt.Errorf("row %d: expected 4 cells, got %d", y+1, len(line))
		}

		for x := range line {
			switch line[x] {
			case '.':
			case '#':
				if blockCount < len(blocks) {
					blocks[blockCount] = point{x: x, y: y}
				}

				blockCount++
			default:
				return tetromino{}, fmt.Errorf("row %d, column %d: invalid character %q", y+1, x+1, line[x])
			}
		}
	}

	if blockCount != len(blocks) {
		return tetromino{}, fmt.Errorf("expected 4 blocks, got %d", blockCount)
	}
	if !blocksAreConnected(blocks) {
		return tetromino{}, fmt.Errorf("blocks are not edge-connected")
	}

	blocks, width, height := normalize(blocks)

	return tetromino{
		blocks: blocks,
		width:  width,
		height: height,
		letter: letter,
	}, nil
}
