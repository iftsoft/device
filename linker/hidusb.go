package linker

import (
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
)

type HidUsbLink struct {
	config *config.HidUsbConfig
	log    *core.LogAgent
}

func NewDummyLinker(cfg *config.HidUsbConfig) *HidUsbLink {
	h := HidUsbLink{
		config: cfg,
		log:    core.GetLogAgent(core.LogLevelTrace, "HidUsb"),
	}
	return &h
}

func (h HidUsbLink) Open() error {
	panic("implement me")
}

func (h HidUsbLink) Close() error {
	panic("implement me")
}

func (h HidUsbLink) Reset() error {
	panic("implement me")
}

func (h HidUsbLink) Write(data []byte) error {
	panic("implement me")
}

func (h HidUsbLink) Read(wait int) ([]byte, error) {
	panic("implement me")
}
