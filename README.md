# Tetris Optimizer

A standard-library Go CLI that packs fixed-orientation tetrominoes into the
smallest square that can contain them all.

The program validates the complete input before searching, preserves the order
of the submitted pieces through `A`–`Z` labels, and prints `ERROR` for malformed
input or an invalid invocation.

## Highlights

- Strict parsing of 4×4 tetromino grids and their blank-line separators
- Edge-connectivity validation for every four-block piece
- Minimum-size search rather than a merely compact first result
- Recursive backtracking with duplicate-shape symmetry breaking
- Empty-region pruning based on the number of cells each region must leave free
- Tests covering malformed inputs, all canonical tetromino families, and the
  official audit's expected board sizes

## Usage

The command accepts exactly one text-file path:

```console
go run . pieces.txt
```

Each piece uses `#` for its four blocks and `.` for empty cells in a 4×4 grid.
Separate pieces with one empty line and end the file with a newline.

```text
#...
#...
#...
#...

....
....
..##
..##
```

The solver keeps each piece's orientation and labels it by input order:

```text
ABB.
ABB.
A...
A...
```

Invalid input produces:

```text
ERROR
```

## How it works

1. The parser validates row dimensions, characters, block counts, separators,
   connectivity, and the final newline, then normalizes each piece to `(0, 0)`.
2. The solver starts at `ceil(sqrt(piece count × 4))`, also accounting for the
   widest and tallest piece.
3. Pieces with fewer candidate positions are placed first. Identical shapes use
   an ordered-placement constraint so equivalent permutations are explored only
   once.
4. After each placement, connected empty regions are inspected. Because one
   tetromino always removes four connected cells, a region's size modulo four is
   a lower bound on the holes it must leave. Branches exceeding the board's hole
   budget are discarded.
5. If no arrangement exists at the current size, the board grows by one cell in
   each dimension and the search repeats. The first solution is therefore
   minimal.

## Development

```console
go test ./...
go test -race ./...
go vet ./...
go build ./...
```

The suite includes the space-count expectations from the
[official 01 Edu audit](https://github.com/01-edu/public/tree/master/subjects/tetris-optimizer/audit).

## Project structure

```text
.
├── main.go          # CLI boundary and output contract
├── parser.go        # File-format and tetromino validation
├── solver.go        # Minimum-square backtracking search
├── tetromino.go     # Core shape data and normalization helpers
└── *_test.go        # Parser, CLI, solver, and audit regression tests
```

Built as part of the Zone01 Athens curriculum by
[Stefanos Kamprogiannis](https://github.com/skamprogiannis).
