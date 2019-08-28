package proxy

import (
	"encoding/json"
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type SystemServer struct {
	scopeItem *duplex.ScopeItem
	server    duplex.ServerManager
	callback  common.SystemCallback
	log       *core.LogAgent
}

func NewSystemServer() *SystemServer {
	ss := SystemServer{
		scopeItem: duplex.NewScopeItem(duplex.ScopeSystem),
		server:    nil,
		callback:  nil,
		log:       nil,
	}
	return &ss
}

func (ss *SystemServer) GetScopeItem() *duplex.ScopeItem {
	return ss.scopeItem
}

func (ss *SystemServer) Init(server duplex.ServerManager, callback common.SystemCallback, log *core.LogAgent) {
	ss.log = log
	ss.server = server
	ss.callback = callback
	if ss.scopeItem != nil {
		ss.scopeItem.SetScopeFunc(common.CmdSystemReply, func(name string, dump []byte) {
			reply, err := ss.decodeReply(name, common.CmdSystemReply, dump)
			if err == nil && ss.callback != nil {
				err = ss.callback.SystemReply(name, reply)
			}
		})
		if ss.server != nil {
			ss.server.AddScopeItem(ss.scopeItem)
		}
	}
}

func (ss *SystemServer) decodeReply(name string, cmd string, dump []byte) (reply *common.SystemReply, err error) {
	if ss.log != nil {
		ss.log.Dump("SystemServer for dev:%s get cmd:%s, pack:%s", name, cmd, string(dump))
	}
	reply = &common.SystemReply{}
	err = json.Unmarshal(dump, reply)
	return reply, err
}

func (ss *SystemServer) SendSystemCommand(name string, cmd string, query interface{}) error {
	if ss.server == nil {
		return errors.New("ServerManager is not set for SystemServer")
	}
	transport := ss.server.GetTransporter(name)
	if transport == nil {
		return errors.New("SystemServer can't get transport to device")
	}
	dump, err := json.Marshal(query)
	if err != nil {
		return err
	}
	if ss.log != nil {
		ss.log.Dump("SystemServer dev:%s run cmd:%s, pack:%s", name, cmd, string(dump))
	}
	pack := duplex.NewPacket(duplex.ScopeSystem, name, cmd, dump)
	err = transport.SendPacket(pack)
	return err
}
