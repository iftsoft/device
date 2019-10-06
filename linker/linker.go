package linker

import (
	"errors"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
)

var (
	errPortNotOpen   = errors.New("port is not open")
)
const (
	linkerBufferSize = 1024
)

type PortReader interface {
	OnRead(data []byte) int
}

type PortLinker interface {
	Open() error
	Close() error
	Flash() error
	IsOpen() bool
	Write(data []byte) (int, error)
}

func GetPortLinker(cfg *config.LinkerConfig, call PortReader) PortLinker {
	dummy := NewDummyLink(call)
	if cfg == nil {
		return dummy
	}
	switch cfg.LinkType {
	case config.LinkTypeNone:
		return dummy
	case config.LinkTypeSerial:
		return NewSerialLink(cfg.Serial, call)
	case config.LinkTypeHidUsb:
		return NewDummyLinker(cfg.HidUsb, call)
	}
	return dummy
}

func GetLinkerPorts(out *core.LogAgent) error {
	_, err := EnumerateSerialPorts(out)
	if err != nil {
		out.Error("Serial port error: %s", err.Error())
	}
	_, err = EnumerateHidUsbPorts(out)
	if err != nil {
		out.Error("HidUsb port error: %s", err.Error())
	}
	return err
}

