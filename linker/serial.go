package linker

import (
	"errors"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"go.bug.st/serial.v1"
)

var (
	errPortNotOpen = errors.New("serial port is not open")
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

func EnumerateSerialPorts(out *core.LogAgent) (list []string, err error) {
	out.Debug("Serial port enumeration")
	list, err = serial.GetPortsList()
	for i, ser := range list {
		out.Dump("   Port#%d - %s", i, ser)
	}
	return list, err
}

func (s *SerialLink) Open() (err error) {
	if s.config == nil {
		return errors.New("serial config is not set")
	}
	m := &serial.Mode{
		BaudRate: int(s.config.BaudRate),
		DataBits: int(s.config.DataBits),
		Parity:   serial.OddParity,  //serial.Parity(s.config.Parity),
		StopBits: serial.OneStopBit, //.StopBits(s.config.StopBits),
	}
	s.port, err = serial.Open(s.config.PortName, m)
	if s.port == nil && err == nil {
		err = errPortNotOpen
	}
	s.log.Debug("Open serial port %s return %s", s.config.PortName, core.GetErrorText(err))
	return err
}

func (s *SerialLink) Close() (err error) {
	if s.port == nil {
		return err
	}
	err = s.port.Close()
	s.log.Debug("Close serial port %s return %s", s.config.PortName, core.GetErrorText(err))
	if err == nil {
		s.port = nil
	}
	return err
}

func (s *SerialLink) Flash() error {
	var err error
	if s.port == nil {
		err = errPortNotOpen
	}
	if err == nil {
		err = s.port.ResetInputBuffer()
	}
	if err == nil {
		err = s.port.ResetOutputBuffer()
	}
	s.log.Debug("Flash serial port %s return %s", s.config.PortName, core.GetErrorText(err))
	return err
}

func (s *SerialLink) Write(data []byte) (n int, err error) {
	if s.port == nil {
		return 0, errPortNotOpen
	}
	n, err = s.port.Write(data)
	s.log.Debug("Write to serial port %s return %s", s.config.PortName, core.GetErrorText(err))
	return n, err
}

func (s *SerialLink) Read(data []byte) (n int, err error) {
	if s.port == nil {
		return 0, errPortNotOpen
	}
	n, err = s.port.Read(data)
	s.log.Debug("Read from serial port %s return %s", s.config.PortName, core.GetErrorText(err))
	return n, err
}
