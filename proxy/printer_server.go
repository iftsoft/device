package proxy

import (
	"encoding/json"
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type PrinterServer struct {
	server    duplex.ServerManager
	callback  common.PrinterCallback
	log       *core.LogAgent
}

func NewPrinterServer() *PrinterServer {
	ps := PrinterServer{
		server:    nil,
		callback:  nil,
		log:       nil,
	}
	return &ps
}


func (ps *PrinterServer) Init(server duplex.ServerManager, callback common.PrinterCallback, log *core.LogAgent) {
	ps.log = log
	ps.server = server
	ps.callback = callback
	if ps.server != nil {
		ps.server.AddDispatcher(duplex.ScopePrinter, ps)
	}
}

func (ps *PrinterServer) EvalPacket(pack *duplex.Packet) error {
	if pack == nil {
		return errors.New("duplex Packet is nil")
	}
	switch pack.Command {
	case common.CmdPrinterProgress:
		reply := &common.PrinterProgress{}
		err := ps.decodeReply(pack.DevName, pack.Command, pack.Content, reply)
		if err == nil && ps.callback != nil {
			err = ps.callback.PrinterProgress(pack.DevName, reply)
		}
		return err

	default:
		ps.log.Warn("PrinterServer EvalPacket: Unknown  command - %s", pack.Command)
		return errors.New("duplex Packet unknown command")
	}
}

func (ps *PrinterServer) decodeReply(name string, cmd string, dump []byte, reply interface{}) (err error) {
	if ps.log != nil {
		ps.log.Dump("PrinterServer dev:%s take cmd:%s, pack:%s", name, cmd, string(dump))
	}
	err = json.Unmarshal(dump, reply)
	return err
}

func (ps *PrinterServer) SendPrinterCommand(name string, cmd string, query interface{}) error {
	if ps.server == nil {
		return errors.New("ServerManager is not set for PrinterServer")
	}
	transport := ps.server.GetTransporter(name)
	if transport == nil {
		return errors.New("PrinterServer can't get transport to device")
	}
	dump, err := json.Marshal(query)
	if err != nil {
		return err
	}
	if ps.log != nil {
		ps.log.Dump("PrinterServer dev:%s send cmd:%s, pack:%s", name, cmd, string(dump))
	}
	pack := duplex.NewPacket(duplex.ScopePrinter, name, cmd, dump)
	err = transport.SendPacket(pack)
	return err
}
