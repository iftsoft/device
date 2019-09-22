package linker

import (
	"errors"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"github.com/karalabe/hid"
)

type HidUsbLink struct {
	config *config.HidUsbConfig
	log    *core.LogAgent
	link   *hid.Device
}

func NewDummyLinker(cfg *config.HidUsbConfig) *HidUsbLink {
	h := HidUsbLink{
		config: cfg,
		log:    core.GetLogAgent(core.LogLevelTrace, "HidUsb"),
		link:   nil,
	}
	return &h
}

func (h HidUsbLink) Open() (err error) {
	if h.config == nil {
		return errors.New("HidUsb config is not set")
	}
	info := hid.DeviceInfo{
		VendorID:  h.config.VendorID,
		ProductID: h.config.ProductID,
		Serial:    h.config.Serial,
	}
	h.link, err = info.Open()
	return err
}

func (h HidUsbLink) Close() (err error) {
	if h.link == nil {
		return err
	}
	err = h.link.Close()
	if err == nil {
		h.link = nil
	}
	return err
}

func (h HidUsbLink) Reset(ResetMode) error {
	return nil
}

func (h HidUsbLink) Write(data []byte) (n int, err error) {
	if h.link == nil {
		err = errors.New("serial port is not open")
	}
	n, err = h.link.Write(data)
	return n, err
}

func (h HidUsbLink) Read(data []byte) (n int, err error) {
	if h.link == nil {
		err = errors.New("serial port is not open")
	}
	n, err = h.link.Read(data)
	return n, err
}
