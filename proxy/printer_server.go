package proxy

import (
	"encoding/json"
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type PrinterServer struct {
	scopeItem *duplex.ScopeItem
	server    duplex.ServerManager
	callback  common.PrinterCallback
	log       *core.LogAgent
}

func NewPrinterServer() *PrinterServer {
	ps := PrinterServer{
		scopeItem: duplex.NewScopeItem(duplex.ScopePrinter),
		server:    nil,
		callback:  nil,
		log:       nil,
	}
	return &ps
}

func (ps *PrinterServer) GetScopeItem() *duplex.ScopeItem {
	return ps.scopeItem
}

func (ps *PrinterServer) Init(server duplex.ServerManager, callback common.PrinterCallback, log *core.LogAgent) {
	ps.log = log
	ps.server = server
	ps.callback = callback
	if ps.scopeItem != nil {
		ps.scopeItem.SetScopeFunc(common.CmdPrinterProgress, func(name string, dump []byte) {
			reply := &common.PrinterProgress{}
			err := ps.decodeReply(name, common.CmdPrinterProgress, dump, reply)
			if err == nil && ps.callback != nil {
				err = ps.callback.PrinterProgress(name, reply)
			}
		})
		if ps.server != nil {
			ps.server.AddScopeItem(ps.scopeItem)
		}
	}
}

func (rs *PrinterServer) decodeReply(name string, cmd string, dump []byte, reply interface{}) (err error) {
	if rs.log != nil {
		rs.log.Dump("PrinterServer dev:%s take cmd:%s, pack:%s", name, cmd, string(dump))
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
