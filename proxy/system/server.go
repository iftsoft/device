package system

import (
	"encoding/json"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type SystemServer struct {
	scopeId   duplex.PacketScope
	scopeItem *duplex.ScopeItem
	transport duplex.Transporter
	callback  common.SystemCallback
	//	commands     common.SystemManager
	log *core.LogAgent
}

func (ss *SystemServer) Init(trans duplex.Transporter,
	callback common.SystemCallback, log *core.LogAgent) {
	ss.log = log
	ss.transport = trans
	ss.callback = callback
	ss.scopeId = duplex.ScopeSystem
	ss.scopeItem = duplex.NewScopeItem(ss.scopeId)
	ss.scopeItem.SetScopeFunc("CommandReply", func(dump []byte) {
		if ss.log != nil {
			ss.log.Trace("SystemServer get cmd:CommandReply, pack:%s", string(dump))
		}
		reply := &common.SystemReply{}
		err := json.Unmarshal(dump, reply)
		if err == nil && ss.callback != nil {
			err = ss.callback.CommandReply(reply)
		}
	})
}

// Implemetation of common.SystemManager

func (ss *SystemServer) Config(query *common.SystemQuery) error {
	dump, err := json.Marshal(query)
	if err != nil {
		return err
	}
	if ss.log != nil {
		ss.log.Trace("SystemServer run cmd:Config, pack:%s", string(dump))
	}
	pack := duplex.NewRequest(ss.scopeId)
	pack.Command = "Config"
	pack.Content = dump
	err = ss.transport.SendPacket(pack, "")
	return err
}

func (ss *SystemServer) Inform(query *common.SystemQuery) error {
	dump, err := json.Marshal(query)
	if err != nil {
		return err
	}
	if ss.log != nil {
		ss.log.Trace("SystemServer run cmd:Inform, pack:%s", string(dump))
	}
	pack := duplex.NewRequest(ss.scopeId)
	pack.Command = "Inform"
	pack.Content = dump
	err = ss.transport.SendPacket(pack, "")
	return err
}

func (ss *SystemServer) Start(query *common.SystemQuery) error {
	dump, err := json.Marshal(query)
	if err != nil {
		return err
	}
	if ss.log != nil {
		ss.log.Trace("SystemServer run cmd:Start, pack:%s", string(dump))
	}
	pack := duplex.NewRequest(ss.scopeId)
	pack.Command = "Start"
	pack.Content = dump
	err = ss.transport.SendPacket(pack, "")
	return err
}

func (ss *SystemServer) Stop(query *common.SystemQuery) error {
	dump, err := json.Marshal(query)
	if err != nil {
		return err
	}
	if ss.log != nil {
		ss.log.Trace("SystemServer run cmd:Stop, pack:%s", string(dump))
	}
	pack := duplex.NewRequest(ss.scopeId)
	pack.Command = "Stop"
	pack.Content = dump
	err = ss.transport.SendPacket(pack, "")
	return err
}

func (ss *SystemServer) Restart(query *common.SystemQuery) error {
	dump, err := json.Marshal(query)
	if err != nil {
		return err
	}
	if ss.log != nil {
		ss.log.Trace("SystemServer run cmd:Restart, pack:%s", string(dump))
	}
	pack := duplex.NewRequest(ss.scopeId)
	pack.Command = "Restart"
	pack.Content = dump
	err = ss.transport.SendPacket(pack, "")
	return err
}
