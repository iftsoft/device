package linker

import (
	"errors"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"time"
)

type Loopback struct {
	config  *config.DeviceConfig
	log     *core.LogAgent
	linker  PortLinker
	reply   chan []byte
	timeout uint16
}

func GetLoopback(cfg *config.DeviceConfig, log *core.LogAgent) *Loopback {
	lb := &Loopback{
		config:  cfg,
		log:     log,
		linker:  nil,
		reply:   make(chan []byte),
		timeout: cfg.Linker.Timeout,
	}
	lb.linker = GetPortLinker(cfg.Linker, lb)
	if lb.timeout == 0 {
		lb.timeout = 250
	}
	return lb
}

func (lb *Loopback) OnRead(data []byte) int {
	lb.log.Dump("Loopback OnRead data : %s", core.GetBinaryDump(data))
	if len(data) < 4 {
		return 0
	}
	lb.reply <- data[0:4]
	return 4
}

func (lb *Loopback) OpenLink() error {
	if lb.linker == nil {
		return errors.New("linker not set")
	}
	err := lb.linker.Open()
	lb.log.Trace("Loopback OpenLink return : %s", core.GetErrorText(err))
	return err
}

func (lb *Loopback) CloseLink() error {
	if lb.linker == nil {
		return errors.New("linker not set")
	}
	err := lb.linker.Close()
	lb.log.Trace("Loopback CloseLink return : %s", core.GetErrorText(err))
	return err
}

func (lb *Loopback) CheckLink() error {
	if lb.linker == nil {
		return errors.New("linker not set")
	}
	data := []byte{0xAA, 0x55, 0x00, 0xFF}
	dump := []byte{}
	n, err := lb.linker.Write(data)
	if err != nil {
		return err
	}

	tick := time.NewTicker(time.Millisecond)
	defer tick.Stop()

	start := time.Now()
	for {
		select {
		case dump = <-lb.reply:
			goto StopWait
		case tm := <-tick.C:
			delta := uint16(tm.Sub(start) / time.Millisecond)
			if delta > lb.timeout {
				lb.log.Warn("Loopback timeout (ms): %d", delta)
				return errors.New("linker timeout")
			}
		}
	}
StopWait:
	lb.log.Dump("Loopback check data : %s", core.GetBinaryDump(dump))

	if n != len(dump) {
		return errors.New("wrong byte count")
	}
	if data[0] != dump[0] &&
		data[1] != dump[1] &&
		data[2] != dump[2] &&
		data[3] != dump[3] {
		return errors.New("wrong dump data")
	}
	return nil
}
