package proxy

import (
	"encoding/json"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type SystemClient struct {
	scopeItem *duplex.ScopeItem
	commands  common.SystemManager
	log       *core.LogAgent
}

func NewSystemClient() *SystemClient {
	sc := SystemClient{
		scopeItem: duplex.NewScopeItem(duplex.ScopeSystem),
		commands:  nil,
		log:       nil,
	}
	return &sc
}

func (sc *SystemClient) Init(command common.SystemManager, log *core.LogAgent) {
	sc.log = log
	sc.commands = command
	// init scope functions
	if sc.scopeItem != nil {
		sc.scopeItem.SetScopeFunc(common.CmdSystemTerminate, func(name string, dump []byte) {
			query := &common.SystemQuery{}
			err := sc.decodeQuery(name, common.CmdSystemTerminate, dump, query)
			if err == nil && sc.commands != nil {
				err = sc.commands.Terminate(name, query)
			}
		})
		sc.scopeItem.SetScopeFunc(common.CmdSystemInform, func(name string, dump []byte) {
			query := &common.SystemQuery{}
			err := sc.decodeQuery(name, common.CmdSystemInform, dump, query)
			if err == nil && sc.commands != nil {
				err = sc.commands.SysInform(name, query)
			}
		})
		sc.scopeItem.SetScopeFunc(common.CmdSystemStart, func(name string, dump []byte) {
			query := &common.SystemConfig{}
			err := sc.decodeQuery(name, common.CmdSystemStart, dump, query)
			if err == nil && sc.commands != nil {
				err = sc.commands.SysStart(name, query)
			}
		})
		sc.scopeItem.SetScopeFunc(common.CmdSystemStop, func(name string, dump []byte) {
			query := &common.SystemQuery{}
			err := sc.decodeQuery(name, common.CmdSystemStop, dump, query)
			if err == nil && sc.commands != nil {
				err = sc.commands.SysStop(name, query)
			}
		})
		sc.scopeItem.SetScopeFunc(common.CmdSystemRestart, func(name string, dump []byte) {
			query := &common.SystemConfig{}
			err := sc.decodeQuery(name, common.CmdSystemRestart, dump, query)
			if err == nil && sc.commands != nil {
				err = sc.commands.SysRestart(name, query)
			}
		})
	}
}

func (sc *SystemClient) decodeQuery(name string, cmd string, dump []byte, query interface{}) error {
	if sc.log != nil {
		sc.log.Dump("SystemClient dev:%s take cmd:%s, pack:%s", name, cmd, string(dump))
	}
	return json.Unmarshal(dump, query)
}

func (sc *SystemClient) GetScopeItem() *duplex.ScopeItem {
	return sc.scopeItem
}
