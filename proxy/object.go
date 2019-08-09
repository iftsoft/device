package proxy

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
	"github.com/iftsoft/device/proxy/system"
)

type ObjectProxy struct {
	server duplex.ServerManager
	system *system.SystemServer
	log    *core.LogAgent
}

func NewObjectProxy() *ObjectProxy {
	op := ObjectProxy{
		server: nil,
		system: system.NewSystemServer(),
		log:    nil, //core.GetLogAgent(core.LogLevelTrace, "Object"),
	}
	return &op
}

func (op *ObjectProxy) Init(server duplex.ServerManager, log *core.LogAgent) {
	op.server = server
	op.log = log
}

func (op *ObjectProxy) SetSystemCallback(callback common.SystemCallback) {
	op.system.Init(op.server, callback, op.log)
}

// Implemetation of common.SystemManager
func (op *ObjectProxy) Config(name string, query *common.SystemQuery) error {
	return op.system.SendSystemCommand(name, common.CmdSystemConfig, query)
}

func (op *ObjectProxy) Inform(name string, query *common.SystemQuery) error {
	return op.system.SendSystemCommand(name, common.CmdSystemInform, query)
}

func (op *ObjectProxy) Start(name string, query *common.SystemQuery) error {
	return op.system.SendSystemCommand(name, common.CmdSystemStart, query)
}

func (op *ObjectProxy) Stop(name string, query *common.SystemQuery) error {
	return op.system.SendSystemCommand(name, common.CmdSystemStop, query)
}

func (op *ObjectProxy) Restart(name string, query *common.SystemQuery) error {
	return op.system.SendSystemCommand(name, common.CmdSystemRestart, query)
}
