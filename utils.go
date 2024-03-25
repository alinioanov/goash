package goash

import (
	"os"
	"strings"

	"github.com/creack/pty"
)

// countLines counts the number of terminal lines for a string. This means counting
// not just the line ends but also the wrapped lines of strings longer than the number
// of terminal columns.
func countLines(s string) int {
	count := 1
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			count += 1
		}
	}
	return count
}

// wrap returns the wrapped and indented version of a string.
func wrapindent(s string) string {
	_, cols, _ := pty.Getsize(os.Stdout)
	cols = 80
	ret := ""
	l := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			l = 0
		} else {
			l += 1
		}
		if l > cols-len(Indent) {
			ret += "\n"
			l = 1
		}
		ret += string(s[i])
	}
	//
	ret = strings.ReplaceAll(ret, "\n", "\r\n"+Indent)
	return ret
}
