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

func EnumerateHidUsbPorts(out *core.LogAgent) (list []*config.HidUsbConfig, err error) {
	out.Debug("HidUsb port enumeration")
	units := hid.Enumerate(0, 0)
	if units == nil {
		return nil, errors.New("hidapi library is not working")
	}
	for i, unit := range units {
		out.Dump("   Port#%d - %d:%d/%s (%s - %s, %s)", i,
			unit.VendorID, unit.ProductID, unit.Serial, unit.Manufacturer, unit.Product, unit.Path)
		item := &config.HidUsbConfig{
			VendorID:  unit.VendorID,
			ProductID: unit.ProductID,
			Serial:    unit.Serial,
		}
		list = append(list, item)
	}
	return list, err
}

func (h *HidUsbLink) Open() (err error) {
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

func (h *HidUsbLink) Close() (err error) {
	if h.link == nil {
		return err
	}
	err = h.link.Close()
	if err == nil {
		h.link = nil
	}
	return err
}

func (h *HidUsbLink) Flash() error {
	return nil
}

func (h *HidUsbLink) Write(data []byte) (n int, err error) {
	if h.link == nil {
		err = errors.New("serial port is not open")
	}
	n, err = h.link.Write(data)
	return n, err
}

func (h *HidUsbLink) Read(data []byte) (n int, err error) {
	if h.link == nil {
		err = errors.New("serial port is not open")
	}
	n, err = h.link.Read(data)
	return n, err
}
