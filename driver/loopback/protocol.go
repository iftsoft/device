package loopback

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"time"
)


type LoopbackProtocol struct {
	LoopbackLinker
	timeout uint16
	DevState   common.EnumDevState
}

func GetLoopbackProtocol(cfg *config.LinkerConfig) *LoopbackProtocol {
	lbp := &LoopbackProtocol{
		LoopbackLinker: LoopbackLinker{},
		timeout:    cfg.Timeout,
		DevState:   0,
	}
	lbp.InitLinker(cfg)
	if lbp.timeout == 0 {
		lbp.timeout = 250
	}
	return lbp
}


////////////////////////////////////////////////////////////////

func (lbp *LoopbackProtocol) CheckLink() error {
	data := []byte{0xAA, 0x55, 0x00, 0xFF}
	back, err := lbp.exchange(data)
	if err == nil {
		err = lbp.checkReply(data, back)
	}
	lbp.logError("CheckLink", err)
	return err
}

////////////////////////////////////////////////////////////////

func (lbp *LoopbackProtocol) logError(cmd string, err error) {
	code, text := common.CheckError(err)
	lbp.log.Trace("LoopbackProtocol.%s return: %d - %s", cmd, code, text)
}

func (lbp *LoopbackProtocol) checkReply(data, back []byte) error {
	if len(data) != len(back) {
		return common.NewError(common.DevErrorLinkerFault, "length mismatch")
	}
	for i:=0; i<len(back); i++ {
		if data[i] != back[i] {
			return common.NewError(common.DevErrorLinkerFault, "data mismatch")
		}
	}
	return nil
}

func (lbp *LoopbackProtocol) exchange(data []byte) ([]byte, error) {
	err := lbp.writeData(data)
	if err != nil {
		return nil, err
	}
	return lbp.readData(lbp.timeout)
}

func (lbp *LoopbackProtocol) writeData(data []byte) error {
	lbp.log.Dump("LoopbackProtocol writeData data : %s", core.GetBinaryDump(data))
	err := lbp.writeToPort(data)
	return err
}

func (lbp *LoopbackProtocol) readData(timeout uint16) ([]byte, error) {
	tick := time.NewTicker(time.Millisecond)
	defer tick.Stop()

	dump := []byte{}
	start := time.Now()
	for {
		select {
		case dump = <-lbp.reply:
			lbp.log.Dump("LoopbackProtocol check data : %s", core.GetBinaryDump(dump))
			return dump, nil
		case tm := <-tick.C:
			delta := uint16(tm.Sub(start) / time.Millisecond)
			if delta > timeout {
				lbp.log.Warn("LoopbackProtocol timeout (ms): %d", delta)
				return nil, common.NewError(common.DevErrorLinkerTimeout, "linker timeout")
			}
		}
	}
}
