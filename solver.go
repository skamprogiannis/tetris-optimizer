package main

import (
	"sort"
	"strings"
)

type solverPiece struct {
	tetromino
	signature string
	duplicate bool
}

func solveTetrominoes(tetrominoes []tetromino) ([]byte, int) {
	minimumSize := 1
	for minimumSize*minimumSize < len(tetrominoes)*4 {
		minimumSize++
	}
	for _, piece := range tetrominoes {
		minimumSize = max(minimumSize, max(piece.width, piece.height))
	}

	for size := minimumSize; ; size++ {
		pieces := orderPieces(tetrominoes, size)
		board := make([]byte, size*size)
		for index := range board {
			board[index] = '.'
		}

		positions := make([]int, len(pieces))
		if placePieces(board, size, pieces, positions, 0, size*size-len(pieces)*4) {
			return board, size
		}
	}
}

func orderPieces(tetrominoes []tetromino, size int) []solverPiece {
	pieces := make([]solverPiece, len(tetrominoes))
	for index, piece := range tetrominoes {
		pieces[index] = solverPiece{
			tetromino: piece,
			signature: shapeSignature(piece),
		}
	}

	sort.SliceStable(pieces, func(left, right int) bool {
		leftPositions := (size - pieces[left].width + 1) * (size - pieces[left].height + 1)
		rightPositions := (size - pieces[right].width + 1) * (size - pieces[right].height + 1)
		if leftPositions != rightPositions {
			return leftPositions < rightPositions
		}
		return pieces[left].signature < pieces[right].signature
	})

	for index := 1; index < len(pieces); index++ {
		pieces[index].duplicate = pieces[index-1].signature == pieces[index].signature
	}

	return pieces
}

func shapeSignature(piece tetromino) string {
	var signature strings.Builder
	signature.Grow(len(piece.blocks) * 2)
	for _, block := range piece.blocks {
		signature.WriteByte(byte(block.x))
		signature.WriteByte(byte(block.y))
	}
	return signature.String()
}

func placePieces(board []byte, size int, pieces []solverPiece, positions []int, pieceIndex, allowedHoles int) bool {
	if pieceIndex == len(pieces) {
		return true
	}

	piece := pieces[pieceIndex]
	startPosition := 0
	if piece.duplicate {
		startPosition = positions[pieceIndex-1] + 1
	}

	lastX := size - piece.width
	lastY := size - piece.height
	for position := startPosition; position < size*size; position++ {
		x := position % size
		y := position / size
		if x > lastX || y > lastY || !canPlace(board, size, piece.tetromino, x, y) {
			continue
		}

		setPiece(board, size, piece.tetromino, x, y, piece.letter)
		positions[pieceIndex] = position

		if unavoidableHoles(board, size) <= allowedHoles &&
			placePieces(board, size, pieces, positions, pieceIndex+1, allowedHoles) {
			return true
		}

		setPiece(board, size, piece.tetromino, x, y, '.')
	}

	return false
}

func canPlace(board []byte, size int, piece tetromino, x, y int) bool {
	for _, block := range piece.blocks {
		if board[(y+block.y)*size+x+block.x] != '.' {
			return false
		}
	}
	return true
}

func setPiece(board []byte, size int, piece tetromino, x, y int, value byte) {
	for _, block := range piece.blocks {
		board[(y+block.y)*size+x+block.x] = value
	}
}

// unavoidableHoles returns the minimum number of cells that must remain empty.
// A connected tetromino cannot span separate empty regions, so each region's
// size modulo four is a lower bound on the holes it leaves behind.
func unavoidableHoles(board []byte, size int) int {
	visited := make([]bool, len(board))
	queue := make([]int, 0, len(board))
	holes := 0

	for start, cell := range board {
		if cell != '.' || visited[start] {
			continue
		}

		visited[start] = true
		queue = append(queue[:0], start)
		regionSize := 0

		for head := 0; head < len(queue); head++ {
			position := queue[head]
			regionSize++
			x := position % size
			y := position / size

			if x > 0 {
				queue = visitEmpty(board, visited, queue, position-1)
			}
			if x+1 < size {
				queue = visitEmpty(board, visited, queue, position+1)
			}
			if y > 0 {
				queue = visitEmpty(board, visited, queue, position-size)
			}
			if y+1 < size {
				queue = visitEmpty(board, visited, queue, position+size)
			}
		}

		holes += regionSize % 4
	}

	return holes
}

func visitEmpty(board []byte, visited []bool, queue []int, position int) []int {
	if board[position] == '.' && !visited[position] {
		visited[position] = true
		queue = append(queue, position)
	}
	return queue
}

func formatBoard(board []byte, size int) string {
	var output strings.Builder
	output.Grow(len(board) + size)

	for row := 0; row < size; row++ {
		output.Write(board[row*size : (row+1)*size])
		output.WriteByte('\n')
	}

	return output.String()
}
