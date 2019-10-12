package loopback

import (
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"time"
)


type LoopbackProtocol struct {
	config  *config.DeviceConfig
	log     *core.LogAgent
	linker  *LoopbackLinker
	reply   chan []byte
	timeout uint16
	DevState   common.EnumDevState
}

func GetLoopbackProtocol(linker *LoopbackLinker, cfg *config.DeviceConfig, log *core.LogAgent) *LoopbackProtocol {
	lb := &LoopbackProtocol{
		config:     cfg,
		log:        log,
		linker:     linker,
		reply:      linker.GetReplyChan(),
		timeout:    cfg.Linker.Timeout,
		DevState:   0,
	}
	if lb.timeout == 0 {
		lb.timeout = 250
	}
	return lb
}


////////////////////////////////////////////////////////////////

func (lb *LoopbackProtocol) CheckLink() common.DevReply {
	devErr := common.DevReply{}
	data := []byte{0xAA, 0x55, 0x00, 0xFF}
	back, err := lb.exchange(data)
	if err == nil {
		err = lb.checkReply(data, back)
	}
	if err != nil {
		devErr.Code   = common.DevErrorLinkerFault
		devErr.Reason = err
	}
	lb.logError("CheckLink", err)
	return devErr
}

////////////////////////////////////////////////////////////////

func (lb *LoopbackProtocol) logError(cmd string, err error) {
	if err == nil {
		lb.log.Trace("LoopbackProtocol.%s return: Success", cmd)
	} else {
		lb.log.Error("LoopbackProtocol.%s return: %s", cmd, core.GetErrorText(err))
	}
}

func (lb *LoopbackProtocol) checkReply(data, back []byte) error {
	if len(data) != len(back) {
		return errors.New("length mismatch")
	}
	for i:=0; i<len(back); i++ {
		if data[i] != back[i] {
			return errors.New("data mismatch")
		}
	}
	return nil
}

func (lb *LoopbackProtocol) exchange(data []byte) ([]byte, error) {
	err := lb.writeData(data)
	if err != nil {
		return nil, err
	}
	return lb.readData(lb.timeout)
}

func (lb *LoopbackProtocol) writeData(data []byte) error {
	if lb.linker == nil {
		return errors.New("linker not set")
	}
	lb.log.Dump("LoopbackProtocol write data : %v", data)
	err := lb.linker.Write(data)
	return err
}

func (lb *LoopbackProtocol) readData(timeout uint16) ([]byte, error) {
	tick := time.NewTicker(time.Millisecond)
	defer tick.Stop()

	dump := []byte{}
	start := time.Now()
	for {
		select {
		case dump = <-lb.reply:
			lb.log.Dump("LoopbackProtocol check data : %s", core.GetBinaryDump(dump))
			return dump, nil
		case tm := <-tick.C:
			delta := uint16(tm.Sub(start) / time.Millisecond)
			if delta > timeout {
				lb.log.Warn("LoopbackProtocol timeout (ms): %d", delta)
				return nil, errors.New("linker timeout")
			}
		}
	}
}
