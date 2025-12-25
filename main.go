package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"mtext/internal/editor"
	"mtext/internal/terminal"
)

func main() {
	flag.Parse()
	term := terminal.NewTermiosTerminal(int(os.Stdin.Fd()))
	ed, err := editor.New(term, os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if flag.NArg() > 0 {
		if err := ed.OpenFile(flag.Arg(0)); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	if err := ed.Run(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
