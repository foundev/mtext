//go:build linux || darwin

package terminal

import (
	"syscall"
	"testing"
	"unsafe"
)

type fakeIoctl struct {
	termios syscall.Termios
	winsize winSize
	calls   []uintptr
}

func (f *fakeIoctl) Ioctl(_ uintptr, request uintptr, arg uintptr) error {
	f.calls = append(f.calls, request)
	switch request {
	case ioctlReadTermios:
		ptr := (*syscall.Termios)(unsafe.Pointer(arg))
		*ptr = f.termios
	case ioctlWriteTermios:
		ptr := (*syscall.Termios)(unsafe.Pointer(arg))
		f.termios = *ptr
	case ioctlWindowSize:
		ptr := (*winSize)(unsafe.Pointer(arg))
		*ptr = f.winsize
	}
	return nil
}

func TestTermiosTerminalRawMode(t *testing.T) {
	fake := &fakeIoctl{}
	term := &TermiosTerminal{fd: 1, ioctl: fake}
	if err := term.EnableRawMode(); err != nil {
		t.Fatalf("enable raw mode: %v", err)
	}
	if !term.rawMode {
		t.Fatalf("expected raw mode to be enabled")
	}
	if fake.termios.Cc[syscall.VMIN] != 0 || fake.termios.Cc[syscall.VTIME] != 1 {
		t.Fatalf("expected raw mode control chars to be set")
	}
	if err := term.DisableRawMode(); err != nil {
		t.Fatalf("disable raw mode: %v", err)
	}
	if term.rawMode {
		t.Fatalf("expected raw mode to be disabled")
	}
}

func TestTermiosTerminalSize(t *testing.T) {
	fake := &fakeIoctl{winsize: winSize{Row: 24, Col: 80}}
	term := &TermiosTerminal{fd: 1, ioctl: fake}
	rows, cols, err := term.Size()
	if err != nil {
		t.Fatalf("size: %v", err)
	}
	if rows != 24 || cols != 80 {
		t.Fatalf("expected 24x80, got %dx%d", rows, cols)
	}
}
