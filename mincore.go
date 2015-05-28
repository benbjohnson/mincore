package mincore

import (
	"syscall"
	"unsafe"
)

func Mincore(addr unsafe.Pointer, size uint64, vec []byte) error {
	_, _, errno := syscall.Syscall(
		syscall.SYS_MINCORE,
		uintptr(addr),
		uintptr(size),
		uintptr(unsafe.Pointer(&vec[0])),
	)
	if errno == 0 {
		return nil
	}
	return errno
}
