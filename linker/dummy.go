package linker

import (
	"github.com/iftsoft/device/core"
)

type DummyLink struct {
	log    *core.LogAgent
	reader PortReader
	isOpen bool
}

func NewDummyLink(call PortReader) *DummyLink {
	d := DummyLink{
		log:    core.GetLogAgent(core.LogLevelTrace, "Dummy"),
		reader: call,
		isOpen: false,
	}
	return &d
}

func (d *DummyLink) Open() error {
	d.isOpen = true
	d.log.Trace("DummyLink run cmd:Open")
	return nil
}

func (d *DummyLink) Close() error {
	d.isOpen = false
	d.log.Trace("DummyLink run cmd:Close")
	return nil
}

func (d *DummyLink) Flash() error {
	d.log.Trace("DummyLink run cmd:Flash")
	return nil
}

func (d *DummyLink) IsOpen() bool {
	return d.isOpen
}

func (d *DummyLink) Write(data []byte) (int, error) {
	if d.isOpen == false {
		return 0, errPortNotOpen
	}
	d.log.Dump("DummyLink write data : %v", data)
	go func(dump []byte) {
		if d.reader != nil {
			d.log.Dump("DummyLink read data : %v", dump)
			d.reader.OnRead(dump)
		}
	}(data)
	return len(data), nil
}
