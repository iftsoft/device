package linker

import (
	"errors"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
)

type PortLinker interface {
	Open() error
	Close() error
	Flash() error
	Write(data []byte) (int, error)
	Read(data []byte) (int, error)
}

func GetPortLinker(cfg *config.LinkerConfig) PortLinker {
	dummy := NewDummyLink()
	if cfg == nil {
		return dummy
	}
	switch cfg.LinkType {
	case config.LinkTypeNone:
		return dummy
	case config.LinkTypeSerial:
		return NewSerialLink(cfg.Serial)
	case config.LinkTypeHidUsb:
		return NewDummyLinker(cfg.HidUsb)
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

func CheckPortLoopback(port PortLinker) error {
	if port == nil {
		return errors.New("wrong linker pointer")
	}
	err := port.Flash()
	if err != nil {
		return err
	}

	var nw, nr int
	data := []byte{0xAA, 0x55, 0x00, 0xFF}
	dump := make([]byte, 8)
	nw, err = port.Write(data)
	if err != nil {
		return err
	}
	nr, err = port.Read(dump)
	if err != nil {
		return err
	}
	if nw != nr {
		return errors.New("wrong byte count")
	}
	if data[0] != dump[0] &&
		data[1] != dump[1] &&
		data[2] != dump[2] &&
		data[3] != dump[3] {
		return errors.New("wrong dump data")
	}
	return err
}
