package proxy

import (
	"encoding/json"
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type PinPadServer struct {
	server    duplex.ServerManager
	callback  common.PinPadCallback
	log       *core.LogAgent
}

func NewPinPadServer() *PinPadServer {
	pps := PinPadServer{
		server:    nil,
		callback:  nil,
		log:       nil,
	}
	return &pps
}

func (pps *PinPadServer) Init(server duplex.ServerManager, callback common.PinPadCallback, log *core.LogAgent) {
	pps.log = log
	pps.server = server
	pps.callback = callback
	if pps.server != nil {
		pps.server.AddDispatcher(duplex.ScopePinPad, pps)
	}
}

func (pps *PinPadServer) EvalPacket(pack *duplex.Packet) error {
	if pack == nil {
		return errors.New("duplex Packet is nil")
	}
	switch pack.Command {
	case common.CmdPinPadReply:
		reply := &common.ReaderPinReply{}
		err := pps.decodeReply(pack.DevName, pack.Command, pack.Content, reply)
		if err == nil && pps.callback != nil {
			err = pps.callback.PinPadReply(pack.DevName, reply)
		}
		return err

	default:
		pps.log.Warn("PinPadServer EvalPacket: Unknown  command - %s", pack.Command)
		return errors.New("duplex Packet unknown command")
	}
}

func (pps *PinPadServer) decodeReply(name string, cmd string, dump []byte, reply interface{}) (err error) {
	if pps.log != nil {
		pps.log.Dump("PinPadServer dev:%s take cmd:%s, pack:%s", name, cmd, string(dump))
	}
	err = json.Unmarshal(dump, reply)
	return err
}

func (pps *PinPadServer) SendPinPadCommand(name string, cmd string, query interface{}) error {
	if pps.server == nil {
		return errors.New("ServerManager is not set for PinPadServer")
	}
	transport := pps.server.GetTransporter(name)
	if transport == nil {
		return errors.New("PinPadServer can't get transport to device")
	}
	dump, err := json.Marshal(query)
	if err != nil {
		return err
	}
	if pps.log != nil {
		pps.log.Dump("ReaderServer dev:%s send cmd:%s, pack:%s", name, cmd, string(dump))
	}
	pack := duplex.NewPacket(duplex.ScopePinPad, name, cmd, dump)
	err = transport.SendPacket(pack)
	return err
}
