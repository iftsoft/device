package proxy

import (
	"encoding/json"
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type PinPadServer struct {
	scopeItem *duplex.ScopeItem
	server    duplex.ServerManager
	callback  common.PinPadCallback
	log       *core.LogAgent
}

func NewPinPadServer() *PinPadServer {
	pps := PinPadServer{
		scopeItem: duplex.NewScopeItem(duplex.ScopePinPad),
		server:    nil,
		callback:  nil,
		log:       nil,
	}
	return &pps
}

func (pps *PinPadServer) GetScopeItem() *duplex.ScopeItem {
	return pps.scopeItem
}

func (pps *PinPadServer) Init(server duplex.ServerManager, callback common.PinPadCallback, log *core.LogAgent) {
	pps.log = log
	pps.server = server
	pps.callback = callback
	if pps.scopeItem != nil {
		pps.scopeItem.SetScopeFunc(common.CmdPinPadReply, func(name string, dump []byte) {
			reply := &common.ReaderPinReply{}
			err := pps.decodeReply(name, common.CmdPinPadReply, dump, reply)
			if err == nil && pps.callback != nil {
				err = pps.callback.PinPadReply(name, reply)
			}
		})
		if pps.server != nil {
			pps.server.AddScopeItem(pps.scopeItem)
		}
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
