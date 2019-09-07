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
		sc.scopeItem.SetScopeFunc(common.CmdSystemConfig, func(name string, dump []byte) {
			query, err := sc.decodeQuery(name, common.CmdSystemRestart, dump)
			if err == nil && sc.commands != nil {
				err = sc.commands.Config(name, query)
			}
		})
		sc.scopeItem.SetScopeFunc(common.CmdSystemInform, func(name string, dump []byte) {
			query, err := sc.decodeQuery(name, common.CmdSystemInform, dump)
			if err == nil && sc.commands != nil {
				err = sc.commands.Inform(name, query)
			}
		})
		sc.scopeItem.SetScopeFunc(common.CmdSystemStart, func(name string, dump []byte) {
			query, err := sc.decodeQuery(name, common.CmdSystemStart, dump)
			if err == nil && sc.commands != nil {
				err = sc.commands.Start(name, query)
			}
		})
		sc.scopeItem.SetScopeFunc(common.CmdSystemStop, func(name string, dump []byte) {
			query, err := sc.decodeQuery(name, common.CmdSystemStop, dump)
			if err == nil && sc.commands != nil {
				err = sc.commands.Stop(name, query)
			}
		})
		sc.scopeItem.SetScopeFunc(common.CmdSystemRestart, func(name string, dump []byte) {
			query, err := sc.decodeQuery(name, common.CmdSystemRestart, dump)
			if err == nil && sc.commands != nil {
				err = sc.commands.Restart(name, query)
			}
		})
	}
}

func (sc *SystemClient) decodeQuery(name string, cmd string, dump []byte) (
	query *common.SystemQuery, err error) {
	if sc.log != nil {
		sc.log.Dump("SystemClient dev:%s take cmd:%s, pack:%s", name, cmd, string(dump))
	}
	query = &common.SystemQuery{}
	err = json.Unmarshal(dump, query)
	return query, err
}

func (sc *SystemClient) GetScopeItem() *duplex.ScopeItem {
	return sc.scopeItem
}