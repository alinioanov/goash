package goash

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/creack/pty"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	Shell = "sh"

	Logo          = "[>]"
	LogoOnSuccess = "[.]"
	LogoOnError   = "[!]"
	Indent        = "    "
)

var (
	ptty *os.File
	tty  *os.File
	iout *IndentWriter
)

// Sh runs a shell command written as a Go formatted string.
func Sh(cmds string, vars ...interface{}) error {
	return Shd(fmt.Sprintf(cmds, vars...), cmds, vars...)
}

// Sh runs a shell command written as a Go formatted string and prints to the console
// a command description instead of the actual command.
func Shd(desc, cmds string, vars ...interface{}) error {
	var err error
	if ptty == nil {
		// Initialize the pseudo-terminal.
		ptty, tty, err = pty.Open()
		if err != nil {
			panic(err)
		}
		// Indent the command output.
		iout = NewIndentWriter(os.Stdout)
		// Keep the pseudo-terminals syncronized.
		go io.Copy(iout, ptty)
		go io.Copy(ptty, os.Stdin)
	}
	// Create and set up the command.
	cmd := exec.Command(Shell, "-c", fmt.Sprintf(cmds, vars...))
	cmd.Stdout = tty
	cmd.Stderr = tty
	cmd.Stdin = tty
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Setctty = true
	cmd.SysProcAttr.Setsid = true
	// Print the command or its description.
	fmt.Printf("%s%s %s\n%s", strings.Repeat("\b", len(Indent)), Logo, wrapindent(desc), Indent)
	// Set terminal in raw mode.
	oldTermState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	// Reset the indent writer state before each command.
	iout.Reset()
	// Run the command
	errc := cmd.Run()
	// Revert terminal to previous state.
	_ = terminal.Restore(int(os.Stdin.Fd()), oldTermState)
	// Move  the terminal cursor back to the beginning of the command output.
	for i := 1; i < countLines(iout.Output)+countLines(wrapindent(Indent+desc)); i++ {
		fmt.Print("\033[F\r")
	}
	// Rewrite the command with the appropriate status color.
	if errc == nil {
		fmt.Print(color(fmt.Sprintf("%s%s %s\r\n%s", strings.Repeat("\b", len(Indent)), LogoOnSuccess, wrapindent(desc), Indent), "green"))
	} else {
		fmt.Print(color(fmt.Sprintf("%s%s %s\r\n%s", strings.Repeat("\b", len(Indent)), LogoOnError, wrapindent(desc), Indent), "red"))
	}
	// Rewrite the command output.
	fmt.Printf("%s\r\n", iout.Output)
	//
	return errc
}
