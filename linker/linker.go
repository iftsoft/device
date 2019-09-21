package linker

import "github.com/iftsoft/device/config"

type PortLinker interface {
	Open() error
	Close() error
	Reset() error
	Write(data []byte) error
	Read(wait int) ([]byte, error)
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
