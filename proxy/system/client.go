package system

import (
	"encoding/json"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type SystemClient struct {
	scopeId   duplex.PacketScope
	scopeItem *duplex.ScopeItem
	transport duplex.Transporter
	//	callback     common.SystemCallback
	commands common.SystemManager
	log      *core.LogAgent
}

func (sc *SystemClient) Init(trans duplex.Transporter,
	command common.SystemManager, log *core.LogAgent) {
	sc.log = log
	sc.transport = trans
	sc.commands = command
	sc.scopeId = duplex.ScopeSystem
	sc.scopeItem = duplex.NewScopeItem(sc.scopeId)
	sc.scopeItem.SetScopeFunc("Config", func(dump []byte) {
		if sc.log != nil {
			sc.log.Trace("SystemClient get cmd:Config, pack:%s", string(dump))
		}
		query := &common.SystemQuery{}
		err := json.Unmarshal(dump, query)
		if err == nil && sc.commands != nil {
			err = sc.commands.Config(query)
		}
	})
	sc.scopeItem.SetScopeFunc("Inform", func(dump []byte) {
		if sc.log != nil {
			sc.log.Trace("SystemClient get cmd:Inform, pack:%s", string(dump))
		}
		query := &common.SystemQuery{}
		err := json.Unmarshal(dump, query)
		if err == nil && sc.commands != nil {
			err = sc.commands.Inform(query)
		}
	})
	sc.scopeItem.SetScopeFunc("Start", func(dump []byte) {
		if sc.log != nil {
			sc.log.Trace("SystemClient get cmd:Start, pack:%s", string(dump))
		}
		query := &common.SystemQuery{}
		err := json.Unmarshal(dump, query)
		if err == nil && sc.commands != nil {
			err = sc.commands.Start(query)
		}
	})
	sc.scopeItem.SetScopeFunc("Stop", func(dump []byte) {
		if sc.log != nil {
			sc.log.Trace("SystemClient get cmd:Stop, pack:%s", string(dump))
		}
		query := &common.SystemQuery{}
		err := json.Unmarshal(dump, query)
		if err == nil && sc.commands != nil {
			err = sc.commands.Stop(query)
		}
	})
	sc.scopeItem.SetScopeFunc("Restart", func(dump []byte) {
		if sc.log != nil {
			sc.log.Trace("SystemClient get cmd:Restart, pack:%s", string(dump))
		}
		query := &common.SystemQuery{}
		err := json.Unmarshal(dump, query)
		if err == nil && sc.commands != nil {
			err = sc.commands.Restart(query)
		}
	})
}

// Implemetation of common.SystemCallback

func (sc *SystemClient) CommandReply(reply *common.SystemReply) error {
	dump, err := json.Marshal(reply)
	if err != nil {
		return err
	}
	pack := duplex.NewRequest(sc.scopeId)
	pack.Command = "CommandReply"
	pack.Content = dump
	err = sc.transport.SendPacket(pack, "")
	return err
}
