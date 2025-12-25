package terminal

import (
	"syscall"
	"unsafe"
)

func getTermios(ioctl ioctlCaller, fd uintptr) (*syscall.Termios, error) {
	termios := &syscall.Termios{}
	if err := ioctl.Ioctl(fd, ioctlReadTermios, uintptr(unsafe.Pointer(termios))); err != nil {
		return nil, err
	}
	return termios, nil
}

func setTermios(ioctl ioctlCaller, fd uintptr, termios *syscall.Termios) error {
	return ioctl.Ioctl(fd, ioctlWriteTermios, uintptr(unsafe.Pointer(termios)))
}
