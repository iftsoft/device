package loopback

import (
	"errors"
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

func GetLoopbackLinker(cfg *config.LinkerConfig, log *core.LogAgent) *LoopbackLinker {
	lb := &LoopbackLinker{
		config:  cfg,
		log:     log,
		port:    nil,
		reply:   make(chan []byte),
	}
	lb.port = linker.GetPortLinker(cfg, lb)
	return lb
}

func (lb *LoopbackLinker) GetReplyChan() chan []byte {
	if lb == nil {
		return nil
	}
	return lb.reply
}

func (lb *LoopbackLinker) OpenLink() error {
	if lb.port == nil {
		return errors.New("port not set")
	}
	err := lb.port.Open()
	lb.log.Trace("LoopbackLinker OpenLink return : %s", core.GetErrorText(err))
	return err
}

func (lb *LoopbackLinker) CloseLink() error {
	if lb.port == nil {
		return errors.New("port not set")
	}
	err := lb.port.Close()
	lb.log.Trace("LoopbackLinker CloseLink return : %s", core.GetErrorText(err))
	return err
}

////////////////////////////////////////////////////////////////
// Data flow:  STX, LEN, []DATA, LRC, ETX

func (lb *LoopbackLinker) Write(data []byte) error {
	pack := []byte{linker.STX, 0}
	pack[1] = byte(len(data))
	pack = append(pack, data...)
	pack = append(pack, linker.CalcLRC(pack), linker.ETX)

	lb.log.Dump("LoopbackLinker Write data : %s", core.GetBinaryDump(pack))
	n, err := lb.port.Write(pack)
	if n != len(pack) {
		return errors.New("wrong byte count")
	}
	return err
}

// implementation of PotReader interface

func (lb *LoopbackLinker) OnRead(dump []byte) int {
	lb.log.Dump("LoopbackLinker OnRead data : %s", core.GetBinaryDump(dump))
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
		lb.log.Warn("LoopbackLinker OnRead CRC mismatch - come:%Xd, calc:%xd", lrc1, lrc2)
	}
	lb.reply <- dump[2 : sz+2]
	return sz
}

