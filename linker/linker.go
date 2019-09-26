package linker

import (
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
)

type ResetMode byte

type PortLinker interface {
	Open() error
	Close() error
	Reset(mode ResetMode) error
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
	out.Info("Serial ports")
	serList, err := EnumerateSerialPorts()
	if err == nil {
		for i, ser := range serList {
			out.Info("   Port#%d - %s", i, ser)
		}
	} else {
		out.Error("Serial port error: %s", err.Error())
	}
	out.Info("HID / USB ports")
	hidList, err := EnumerateHidUsbPorts()
	if err == nil {
		for i, hid := range hidList {
			out.Info("   Port#%d - %d:%d/%s", i, hid.VendorID, hid.ProductID, hid.Serial)
		}
	} else {
		out.Error("HidUsb port error: %s", err.Error())
	}
	return err
}
