package goash

import (
	"os"
)

// IndentWriter is a writer that wraps an existing file (like os.Stdout) for
// the purpose of adding indentation.
type IndentWriter struct {
	file   *os.File
	Output string
}

// NewIndentWriter returns a new IndentWriter.
func NewIndentWriter(file *os.File) *IndentWriter {
	return &IndentWriter{file, ""}
}

// Write adds indentation to text and writes it to the wrapped file.
func (iw *IndentWriter) Write(b []byte) (n int, err error) {
	str := string(b)
	str = wrapindent(str)
	iw.Output += str

	nw, err := iw.file.Write([]byte(str))
	if err != nil {
		return nw, err
	}

	return len(b), nil
}

// Reset empties the writer's output state.
func (iw *IndentWriter) Reset() {
	iw.Output = ""
}
