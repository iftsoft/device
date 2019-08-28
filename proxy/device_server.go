package proxy

import (
	"encoding/json"
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type DeviceServer struct {
	scopeItem *duplex.ScopeItem
	server    duplex.ServerManager
	callback  common.DeviceCallback
	log       *core.LogAgent
}

func NewDeviceServer() *DeviceServer {
	ss := DeviceServer{
		scopeItem: duplex.NewScopeItem(duplex.ScopeDevice),
		server:    nil,
		callback:  nil,
		log:       nil,
	}
	return &ss
}

func (ss *DeviceServer) GetScopeItem() *duplex.ScopeItem {
	return ss.scopeItem
}

func (ss *DeviceServer) Init(server duplex.ServerManager, callback common.DeviceCallback, log *core.LogAgent) {
	ss.log = log
	ss.server = server
	ss.callback = callback
	if ss.scopeItem != nil {
		ss.scopeItem.SetScopeFunc(common.CmdDeviceReply, func(name string, dump []byte) {
			reply := &common.DeviceReply{}
			err := ss.decodeReply(name, common.CmdDeviceReply, dump, reply)
			if err == nil && ss.callback != nil {
				err = ss.callback.DeviceReply(name, reply)
			}
		})
		ss.scopeItem.SetScopeFunc(common.CmdExecuteError, func(name string, dump []byte) {
			reply := &common.DeviceError{}
			err := ss.decodeReply(name, common.CmdExecuteError, dump, reply)
			if err == nil && ss.callback != nil {
				err = ss.callback.ExecuteError(name, reply)
			}
		})
		ss.scopeItem.SetScopeFunc(common.CmdStateChanged, func(name string, dump []byte) {
			reply := &common.DeviceState{}
			err := ss.decodeReply(name, common.CmdStateChanged, dump, reply)
			if err == nil && ss.callback != nil {
				err = ss.callback.StateChanged(name, reply)
			}
		})
		ss.scopeItem.SetScopeFunc(common.CmdActionPrompt, func(name string, dump []byte) {
			reply := &common.DevicePrompt{}
			err := ss.decodeReply(name, common.CmdActionPrompt, dump, reply)
			if err == nil && ss.callback != nil {
				err = ss.callback.ActionPrompt(name, reply)
			}
		})
		if ss.server != nil {
			ss.server.AddScopeItem(ss.scopeItem)
		}
	}
}

func (ss *DeviceServer) decodeReply(name string, cmd string, dump []byte, reply interface{}) (err error) {
	if ss.log != nil {
		ss.log.Dump("DeviceServer for dev:%s get cmd:%s, pack:%s", name, cmd, string(dump))
	}
	err = json.Unmarshal(dump, reply)
	return err
}

func (ss *DeviceServer) SendDeviceCommand(name string, cmd string, query interface{}) error {
	if ss.server == nil {
		return errors.New("ServerManager is not set for DeviceServer")
	}
	transport := ss.server.GetTransporter(name)
	if transport == nil {
		return errors.New("DeviceServer can't get transport to device")
	}
	dump, err := json.Marshal(query)
	if err != nil {
		return err
	}
	if ss.log != nil {
		ss.log.Dump("DeviceServer dev:%s run cmd:%s, pack:%s", name, cmd, string(dump))
	}
	pack := duplex.NewPacket(duplex.ScopeDevice, name, cmd, dump)
	err = transport.SendPacket(pack)
	return err
}
