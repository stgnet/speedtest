// +build linux

package api

import (
	"syscall"
)

func BindToDevice(fd int, ifName string) error {
	return syscall.SetsockoptString(fd, syscall.SOL_SOCKET, syscall.SO_BINDTODEVICE, ifName)
}
