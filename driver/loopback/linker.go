package loopback

import (
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/linker"
)

type LoopbackLinker struct {
	config  *config.LinkerConfig
	log     *core.LogAgent
	port    linker.PortLinker
	reply   chan []byte
}

func (lbl *LoopbackLinker) InitLinker(cfg *config.LinkerConfig) {
	lbl.config = cfg
	lbl.port   = linker.GetPortLinker(cfg, lbl)
	lbl.reply  = make(chan []byte)
	lbl.log    = core.GetLogAgent(core.LogLevelDump, "Linker")
}

func (lb *LoopbackLinker) GetReplyChan() chan []byte {
	if lb == nil {
		return nil
	}
	return lb.reply
}

func (lbl *LoopbackLinker) OpenLink() error {
	if lbl.port == nil {
		return common.NewError(common.DevErrorConfigFault, "port not set")
	}
	err := lbl.port.Open()
	lbl.log.Trace("LoopbackLinker OpenLink return : %s", core.GetErrorText(err))
	return common.ExtendError(common.DevErrorLinkerFault, err)
}

func (lbl *LoopbackLinker) CloseLink() error {
	if lbl.port == nil {
		return common.NewError(common.DevErrorConfigFault, "port not set")
	}
	err := lbl.port.Close()
	lbl.log.Trace("LoopbackLinker CloseLink return : %s", core.GetErrorText(err))
	return common.ExtendError(common.DevErrorLinkerFault, err)
}

////////////////////////////////////////////////////////////////
// Data flow:  STX, LEN, []DATA, LRC, ETX

func (lbl *LoopbackLinker) writeToPort(data []byte) error {
	if lbl.port == nil {
		return errors.New("port not set")
	}
	pack := []byte{linker.STX, 0}
	pack[1] = byte(len(data))
	pack = append(pack, data...)
	pack = append(pack, linker.CalcLRC(pack), linker.ETX)

	lbl.log.Dump("LoopbackLinker writeToPort data : %s", core.GetBinaryDump(pack))
	n, err := lbl.port.Write(pack)
	if n != len(pack) {
		return common.NewError(common.DevErrorLinkerFault, "wrong byte count")
	}
	return common.ExtendError(common.DevErrorLinkerFault, err)
}

// implementation of PotReader interface

func (lbl *LoopbackLinker) OnRead(dump []byte) int {
	lbl.log.Dump("LoopbackLinker OnRead data : %s", core.GetBinaryDump(dump))
	// Checking header size
	size := len(dump)
	if size < 4 {
		return 0
	}
	// Looking for for STX
	if dump[0] != linker.STX {
		for i := 1; i < size; i++ {
			if dump[i] == linker.STX {
				return i
			}
		}
		return size
	}
	// Checking full size
	sz := int(dump[1])
	if size < sz+4 {
		return 0
	}
	// Checking CRC
	lrc1 := linker.CalcLRC(dump[0 : sz+2])
	lrc2 := dump[sz+2]
	if lrc1 != lrc2 {
		lbl.log.Warn("LoopbackLinker OnRead CRC mismatch - come:%Xd, calc:%xd", lrc1, lrc2)
	}
	lbl.reply <- dump[2 : sz+2]
	return sz
}

