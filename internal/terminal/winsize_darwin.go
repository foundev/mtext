//go:build darwin

package terminal

import "syscall"

const ioctlWindowSize = syscall.TIOCGWINSZ
