package proxy

import (
	"encoding/json"
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type SystemServer struct {
	server    duplex.ServerManager
	callback  common.SystemCallback
	log       *core.LogAgent
}

func NewSystemServer() *SystemServer {
	ss := SystemServer{
		server:    nil,
		callback:  nil,
		log:       nil,
	}
	return &ss
}

func (ss *SystemServer) Init(server duplex.ServerManager, callback common.SystemCallback, log *core.LogAgent) {
	ss.log = log
	ss.server = server
	ss.callback = callback
	if ss.server != nil {
		ss.server.AddDispatcher(duplex.ScopeSystem, ss)
	}
}

func (ss *SystemServer) EvalPacket(pack *duplex.Packet) error {
	if pack == nil {
		return errors.New("duplex Packet is nil")
	}
	switch pack.Command {
	case common.CmdSystemReply:
		reply := &common.SystemReply{}
		err := ss.decodeReply(pack.DevName, pack.Command, pack.Content, reply)
		if err == nil && ss.callback != nil {
			err = ss.callback.SystemReply(pack.DevName, reply)
		}
		return err

	case common.CmdSystemHealth:
		reply := &common.SystemHealth{}
		err := ss.decodeReply(pack.DevName, pack.Command, pack.Content, reply)
		if err == nil && ss.callback != nil {
			err = ss.callback.SystemHealth(pack.DevName, reply)
		}
		return err

	default:
		ss.log.Warn("SystemServer EvalPacket: Unknown  command - %s", pack.Command)
		return errors.New("duplex Packet unknown command")
	}
}

func (ss *SystemServer) decodeReply(name string, cmd string, dump []byte, reply interface{}) error {
	if ss.log != nil {
		ss.log.Dump("SystemServer dev:%s get cmd:%s, pack:%s", name, cmd, string(dump))
	}
	return json.Unmarshal(dump, reply)
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
		ss.log.Dump("SystemServer dev:%s send cmd:%s, pack:%s", name, cmd, string(dump))
	}
	pack := duplex.NewPacket(duplex.ScopeSystem, name, cmd, dump)
	err = transport.SendPacket(pack)
	return err
}
