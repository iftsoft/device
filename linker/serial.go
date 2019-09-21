package linker

import (
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
)

type SerialLink struct {
	config *config.SerialConfig
	log    *core.LogAgent
}

func NewSerialLink(cfg *config.SerialConfig) *SerialLink {
	s := SerialLink{
		config: cfg,
		log:    core.GetLogAgent(core.LogLevelTrace, "Serial"),
	}
	return &s
}

func (s SerialLink) Open() error {
	panic("implement me")
}

func (s SerialLink) Close() error {
	panic("implement me")
}

func (s SerialLink) Reset() error {
	panic("implement me")
}

func (s SerialLink) Write(data []byte) error {
	panic("implement me")
}

func (s SerialLink) Read(wait int) ([]byte, error) {
	panic("implement me")
}
