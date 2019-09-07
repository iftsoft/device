package proxy

import (
	"encoding/json"
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type ReaderServer struct {
	scopeItem *duplex.ScopeItem
	server    duplex.ServerManager
	callback  common.ReaderCallback
	log       *core.LogAgent
}

func NewReaderServer() *ReaderServer {
	rs := ReaderServer{
		scopeItem: duplex.NewScopeItem(duplex.ScopeReader),
		server:    nil,
		callback:  nil,
		log:       nil,
	}
	return &rs
}

func (rs *ReaderServer) GetScopeItem() *duplex.ScopeItem {
	return rs.scopeItem
}

func (rs *ReaderServer) Init(server duplex.ServerManager, callback common.ReaderCallback, log *core.LogAgent) {
	rs.log = log
	rs.server = server
	rs.callback = callback
	if rs.scopeItem != nil {
		rs.scopeItem.SetScopeFunc(common.CmdCardDescription, func(name string, dump []byte) {
			reply := &common.ReaderCardInfo{}
			err := rs.decodeReply(name, common.CmdCardDescription, dump, reply)
			if err == nil && rs.callback != nil {
				err = rs.callback.CardDescription(name, reply)
			}
		})
		rs.scopeItem.SetScopeFunc(common.CmdChipResponse, func(name string, dump []byte) {
			reply := &common.ReaderChipReply{}
			err := rs.decodeReply(name, common.CmdChipResponse, dump, reply)
			if err == nil && rs.callback != nil {
				err = rs.callback.ChipResponse(name, reply)
			}
		})
		rs.scopeItem.SetScopeFunc(common.CmdPinPadReply, func(name string, dump []byte) {
			reply := &common.ReaderPinReply{}
			err := rs.decodeReply(name, common.CmdPinPadReply, dump, reply)
			if err == nil && rs.callback != nil {
				err = rs.callback.PinPadReply(name, reply)
			}
		})
		if rs.server != nil {
			rs.server.AddScopeItem(rs.scopeItem)
		}
	}
}

func (rs *ReaderServer) decodeReply(name string, cmd string, dump []byte, reply interface{}) (err error) {
	if rs.log != nil {
		rs.log.Dump("ReaderServer dev:%s take cmd:%s, pack:%s", name, cmd, string(dump))
	}
	err = json.Unmarshal(dump, reply)
	return err
}

func (rs *ReaderServer) SendReaderCommand(name string, cmd string, query interface{}) error {
	if rs.server == nil {
		return errors.New("ServerManager is not set for ReaderServer")
	}
	transport := rs.server.GetTransporter(name)
	if transport == nil {
		return errors.New("ReaderServer can't get transport to device")
	}
	dump, err := json.Marshal(query)
	if err != nil {
		return err
	}
	if rs.log != nil {
		rs.log.Dump("ReaderServer dev:%s send cmd:%s, pack:%s", name, cmd, string(dump))
	}
	pack := duplex.NewPacket(duplex.ScopeReader, name, cmd, dump)
	err = transport.SendPacket(pack)
	return err
}
