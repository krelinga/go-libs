package exam

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Loc represents a line of code in a file.
type Loc struct {
	File string
	Line int
}

// String returns a readable string representation of the location.
func (l Loc) String() string {
	if l.File == "" && l.Line == 0 {
		return "<uninitialized>"
	}
	return fmt.Sprintf("%s:%d", l.File, l.Line)
}

// Here returns a Loc instance with the file and line number where Here() was called from.
func Here() Loc {
	return hereOffset(1)
}

func hereOffset(skip int) Loc {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return Loc{File: "unknown", Line: 0}
	}

	// Convert absolute path to relative path from current working directory
	cwd, err := os.Getwd()
	if err != nil {
		// If we can't get the current working directory, return the absolute path
		return Loc{File: file, Line: line}
	}

	relPath, err := filepath.Rel(cwd, file)
	if err != nil {
		// If we can't make it relative, return the absolute path
		return Loc{File: file, Line: line}
	}

	return Loc{File: relPath, Line: line}
}