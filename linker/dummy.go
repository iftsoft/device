package linker

import (
	"github.com/iftsoft/device/core"
)

type DummyLink struct {
	log *core.LogAgent
}

func NewDummyLink() *DummyLink {
	d := DummyLink{
		log: core.GetLogAgent(core.LogLevelTrace, "Dummy"),
	}
	return &d
}

func (d DummyLink) Open() error {
	d.log.Debug("DummyLink run cmd:Open")
	return nil
}

func (d DummyLink) Close() error {
	d.log.Debug("DummyLink run cmd:Close")
	return nil
}

func (d DummyLink) Reset(ResetMode) error {
	d.log.Debug("DummyLink run cmd:Reset")
	return nil
}

func (d DummyLink) Write(data []byte) (int, error) {
	d.log.Debug("DummyLink run cmd:Write")
	return 0, nil
}

func (d DummyLink) Read(data []byte) (int, error) {
	d.log.Debug("DummyLink run cmd:Read")
	return 0, nil
}
