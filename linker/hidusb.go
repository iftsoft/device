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
	reader PortReader
	isOpen bool
}

func NewDummyLinker(cfg *config.HidUsbConfig, call PortReader) *HidUsbLink {
	h := HidUsbLink{
		config: cfg,
		log:    core.GetLogAgent(core.LogLevelTrace, "HidUsb"),
		link:   nil,
		reader: call,
		isOpen: false,
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
	if err == nil {
		go h.readingLoop()
	}
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
	if h.link == nil {
		return errPortNotOpen
	}
	return nil
}

func (h *HidUsbLink) IsOpen() bool {
	return h.isOpen
}

func (h *HidUsbLink) Write(data []byte) (n int, err error) {
	if h.link == nil {
		return 0, errPortNotOpen
	}
	n, err = h.link.Write(data)
	h.log.Trace("Write to hidapi port %d:%d return %s",
		h.config.VendorID, h.config.ProductID, core.GetErrorText(err))
	return n, err
}

func (h *HidUsbLink) readData(data []byte) (n int, err error) {
	if h.link == nil {
		return 0, errPortNotOpen
	}
	n, err = h.link.Read(data)
	h.log.Trace("Read from hidapi port %d:%d return %s",
		h.config.VendorID, h.config.ProductID, core.GetErrorText(err))
	return n, err
}

func (h *HidUsbLink) readingLoop() {
	h.isOpen = true
	defer func() { h.isOpen = false }()
	h.log.Trace("HidUsb reading loop is started")
	defer h.log.Trace("HidUsb reading loop is stopped")

	rest := []byte{}
	buff := make([]byte, linkerBufferSize)
	for {
		n, err := h.readData(buff)
		if n > 0 {
			dump := buff[0:n]
			data := append(rest, dump...)
			rest = h.processData(data)
		}
		if err != nil {
			h.log.Warn("HidUsb ReadData error: %s", err)
			return
		}
	}
}

func (h *HidUsbLink) processData(data []byte) (out []byte) {
	if h.reader == nil {
		return nil
	}
	k := h.reader.OnRead(data)
	if k == 0 {
		return data
	}
	if k > 0 && k < len(data) {
		return h.processData(data[k:])
	}
	return nil
}
