package terminal

import "syscall"

type Terminal interface {
	EnableRawMode() error
	DisableRawMode() error
	Size() (rows int, cols int, err error)
}

type ioctlCaller interface {
	Ioctl(fd uintptr, request uintptr, arg uintptr) error
}

type syscallIoctl struct{}

func (syscallIoctl) Ioctl(fd uintptr, request uintptr, arg uintptr) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, request, arg)
	if errno != 0 {
		return errno
	}
	return nil
}

type TermiosTerminal struct {
	fd       uintptr
	ioctl    ioctlCaller
	original *syscall.Termios
	rawMode  bool
}

func NewTermiosTerminal(fd int) *TermiosTerminal {
	return &TermiosTerminal{
		fd:    uintptr(fd),
		ioctl: syscallIoctl{},
	}
}

func (t *TermiosTerminal) EnableRawMode() error {
	if t.rawMode {
		return nil
	}
	current, err := getTermios(t.ioctl, t.fd)
	if err != nil {
		return err
	}
	original := *current
	makeRaw(current)
	if err := setTermios(t.ioctl, t.fd, current); err != nil {
		return err
	}
	t.original = &original
	t.rawMode = true
	return nil
}

func (t *TermiosTerminal) DisableRawMode() error {
	if !t.rawMode || t.original == nil {
		return nil
	}
	if err := setTermios(t.ioctl, t.fd, t.original); err != nil {
		return err
	}
	t.rawMode = false
	return nil
}

func (t *TermiosTerminal) Size() (rows int, cols int, err error) {
	ws, err := getWinSize(t.ioctl, t.fd)
	if err != nil {
		return 0, 0, err
	}
	return int(ws.Row), int(ws.Col), nil
}
