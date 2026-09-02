package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestExecutePrintsSolution(t *testing.T) {
	path := filepath.Join(t.TempDir(), "pieces.txt")
	if err := os.WriteFile(path, []byte("....\n.##.\n.##.\n....\n"), 0o600); err != nil {
		t.Fatalf("write input: %v", err)
	}

	var output bytes.Buffer
	execute([]string{path}, &output)

	const want = "AA\nAA\n"
	if output.String() != want {
		t.Fatalf("execute() output = %q, want %q", output.String(), want)
	}
}

func TestExecutePrintsErrorForInvalidInvocation(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "no argument"},
		{name: "too many arguments", args: []string{"first", "second"}},
		{name: "missing file", args: []string{filepath.Join(t.TempDir(), "missing.txt")}},
	}

	invalidPath := filepath.Join(t.TempDir(), "invalid.txt")
	if err := os.WriteFile(invalidPath, []byte("####\n...#\n....\n....\n"), 0o600); err != nil {
		t.Fatalf("write invalid input: %v", err)
	}
	tests = append(tests, struct {
		name string
		args []string
	}{name: "invalid tetromino", args: []string{invalidPath}})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer
			execute(test.args, &output)

			const want = "ERROR\n"
			if output.String() != want {
				t.Fatalf("execute() output = %q, want %q", output.String(), want)
			}
		})
	}
}
