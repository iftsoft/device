package linker

import "github.com/iftsoft/device/config"

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
