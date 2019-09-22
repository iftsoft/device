package linker

import (
	"errors"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"go.bug.st/serial.v1"
)

type SerialLink struct {
	config *config.SerialConfig
	log    *core.LogAgent
	port   serial.Port
}

func NewSerialLink(cfg *config.SerialConfig) *SerialLink {
	s := SerialLink{
		config: cfg,
		log:    core.GetLogAgent(core.LogLevelTrace, "Serial"),
		port:   nil,
	}
	return &s
}

func (s SerialLink) Open() (err error) {
	if s.config == nil {
		return errors.New("serial config is not set")
	}
	m := &serial.Mode{
		BaudRate: int(s.config.BaudRate),
		DataBits: int(s.config.DataBits),
		Parity:   serial.Parity(s.config.Parity),
		StopBits: serial.StopBits(s.config.StopBits),
	}
	s.port, err = serial.Open(s.config.PortName, m)
	return err
}

func (s SerialLink) Close() (err error) {
	if s.port == nil {
		return err
	}
	err = s.port.Close()
	if err == nil {
		s.port = nil
	}
	return err
}

func (s SerialLink) Reset(ResetMode) error {
	var err error
	if s.port == nil {
		err = errors.New("serial port is not open")
	}
	if err == nil {
		err = s.port.ResetInputBuffer()
	}
	if err == nil {
		err = s.port.ResetOutputBuffer()
	}
	return err
}

func (s SerialLink) Write(data []byte) (n int, err error) {
	if s.port == nil {
		err = errors.New("serial port is not open")
	}
	n, err = s.port.Write(data)
	return n, err
}

func (s SerialLink) Read(data []byte) (n int, err error) {
	if s.port == nil {
		err = errors.New("serial port is not open")
	}
	n, err = s.port.Read(data)
	return n, err
}
