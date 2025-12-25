//go:build linux

package terminal

import "syscall"

const ioctlWindowSize = syscall.TIOCGWINSZ
