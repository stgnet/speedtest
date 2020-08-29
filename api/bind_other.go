// +build !linux
package api

import (
	"errors"
)

// only linux has BindToDevice call

func BindToDevice(fd int, ifName string) error {
	return errors.New("BindtoDevice is not implemented")
}
