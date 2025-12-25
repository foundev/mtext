package terminal

import "unsafe"

type winSize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func getWinSize(ioctl ioctlCaller, fd uintptr) (*winSize, error) {
	ws := &winSize{}
	if err := ioctl.Ioctl(fd, ioctlWindowSize, uintptr(unsafe.Pointer(ws))); err != nil {
		return nil, err
	}
	return ws, nil
}
