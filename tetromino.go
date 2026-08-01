package main

type point struct {
	x int
	y int
}

type tetromino struct {
	blocks [4]point
	width  int
	height int
	letter byte
}

func blocksAreConnected(blocks [4]point) bool {
	var visited [4]bool
	var queue [4]int
	visited[0] = true
	queue[0] = 0
	visitedCount := 1

	for head := 0; head < visitedCount; head++ {
		current := blocks[queue[head]]

		for index, candidate := range blocks {
			if visited[index] {
				continue
			}

			deltaX := current.x - candidate.x
			deltaY := current.y - candidate.y
			adjacent := deltaX == 0 && (deltaY == 1 || deltaY == -1) ||
				deltaY == 0 && (deltaX == 1 || deltaX == -1)
			if !adjacent {
				continue
			}

			visited[index] = true
			queue[visitedCount] = index
			visitedCount++
		}
	}

	return visitedCount == len(blocks)
}

func normalize(blocks [4]point) ([4]point, int, int) {
	minX := blocks[0].x
	minY := blocks[0].y

	for _, block := range blocks[1:] {
		minX = min(minX, block.x)
		minY = min(minY, block.y)
	}

	width := 0
	height := 0

	for index := range blocks {
		blocks[index].x -= minX
		blocks[index].y -= minY

		width = max(width, blocks[index].x+1)
		height = max(height, blocks[index].y+1)
	}

	return blocks, width, height
}
