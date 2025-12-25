//go:build linux

package terminal

import "syscall"

const (
	ioctlReadTermios  = syscall.TCGETS
	ioctlWriteTermios = syscall.TCSETS
)

func makeRaw(termios *syscall.Termios) {
	termios.Iflag &^= syscall.BRKINT | syscall.ICRNL | syscall.INPCK | syscall.ISTRIP | syscall.IXON
	termios.Oflag &^= syscall.OPOST
	termios.Cflag |= syscall.CS8
	termios.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN | syscall.ISIG
	termios.Cc[syscall.VMIN] = 0
	termios.Cc[syscall.VTIME] = 1
}
