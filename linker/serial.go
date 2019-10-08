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
	reader PortReader
	isOpen bool
}

func NewSerialLink(cfg *config.SerialConfig, call PortReader) *SerialLink {
	s := SerialLink{
		config: cfg,
		log:    core.GetLogAgent(core.LogLevelDebug, "Serial"),
		port:   nil,
		reader: call,
		isOpen: false,
	}
	return &s
}

func EnumerateSerialPorts(out *core.LogAgent) (list []string, err error) {
	defer core.PanicRecover(&err, out)
	out.Debug("Serial port enumeration")
	list, err = serial.GetPortsList()
	for i, ser := range list {
		out.Debug("   Port#%d - %s", i, ser)
	}
	return list, err
}

func (s *SerialLink) Open() (err error) {
	defer core.PanicRecover(&err, s.log)
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
	if s.port == nil && err == nil {
		err = errPortNotOpen
	}
	if err == nil {
		go s.readingLoop()
	} else {
		s.port = nil
	}
	s.log.Trace("Open serial port %s return %s", s.config.PortName, core.GetErrorText(err))
	return err
}

func (s *SerialLink) Close() (err error) {
	defer core.PanicRecover(&err, s.log)
	if s.port == nil {
		return err
	}
	err = s.port.Close()
	s.log.Trace("Close serial port %s return %s", s.config.PortName, core.GetErrorText(err))
	if err == nil {
		s.port = nil
	}
	return err
}

func (s *SerialLink) Flash() (err error) {
	defer core.PanicRecover(&err, s.log)
	if s.port == nil {
		err = errPortNotOpen
	}
	if err == nil {
		err = s.port.ResetInputBuffer()
	}
	if err == nil {
		err = s.port.ResetOutputBuffer()
	}
	s.log.Trace("Flash serial port %s return %s", s.config.PortName, core.GetErrorText(err))
	return err
}

func (s *SerialLink) IsOpen() bool {
	return s.isOpen
}

func (s *SerialLink) Write(data []byte) (n int, err error) {
	defer core.PanicRecover(&err, s.log)
	if s.port == nil {
		return 0, errPortNotOpen
	}
	s.log.Dump("Serial write data : %s", core.GetBinaryDump(data))
	n, err = s.port.Write(data)
	s.log.Trace("Write to serial port %s return %s", s.config.PortName, core.GetErrorText(err))
	return n, err
}

func (s *SerialLink) readData(data []byte) (n int, err error) {
	defer core.PanicRecover(&err, s.log)
	if s.port == nil {
		return 0, errPortNotOpen
	}
	n, err = s.port.Read(data)
	s.log.Dump("Serial read data : %s", core.GetBinaryDump(data[0:n]))
	s.log.Trace("Read from serial port %s of %d bytes return %s",
		s.config.PortName, n, core.GetErrorText(err))
	return n, err
}

func (s *SerialLink) readingLoop() {
	s.isOpen = true
	defer func() { s.isOpen = false }()
	s.log.Trace("Serial reading loop is started")
	defer s.log.Trace("Serial reading loop is stopped")

	rest := []byte{}
	for {
		buff := make([]byte, linkerBufferSize)
		n, err := s.readData(buff)
		if n > 0 {
			dump := buff[0:n]
			data := append(rest, dump...)
			rest = s.processData(data)
		}
		if err != nil {
			s.log.Warn("Serial ReadData error: %s", err)
			return
		}
	}
}

func (s *SerialLink) processData(data []byte) (out []byte) {
	s.log.Dump("Process reply data : %s", core.GetBinaryDump(data))
	if s.reader == nil {
		return nil
	}
	if len(data) == 0 {
		return nil
	}
	k := s.reader.OnRead(data)
	if k == 0 {
		return data
	}
	if k > 0 && k < len(data) {
		return s.processData(data[k:])
	}
	return nil
}
