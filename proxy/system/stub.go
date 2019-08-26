package system

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
)

type SystemStub struct {
	callback common.SystemCallback
	log      *core.LogAgent
}

func NewSystemStub() *SystemStub {
	ss := SystemStub{
		callback: nil,
		log:      nil,
	}
	return &ss
}

func (ss *SystemStub) Init(callback common.SystemCallback, log *core.LogAgent) {
	ss.log = log
	ss.callback = callback
}

func (ss *SystemStub) dummyCommandReply(name string, cmd string, query *common.SystemQuery) error {
	if ss.log != nil {
		ss.log.Trace("SystemStub dev:%s run cmd:%s, data:%s", name, cmd, query.DevName)
	}
	reply := &common.SystemReply{}
	reply.DevName = query.DevName
	reply.Command = cmd
	var err error
	if ss.callback != nil {
		err = ss.callback.CommandReply(name, reply)
	}
	return err
}

// Implementation of common.SystemManager
//
func (ss *SystemStub) Config(name string, query *common.SystemQuery) error {
	return ss.dummyCommandReply(name, common.CmdSystemConfig, query)
}

func (ss *SystemStub) Inform(name string, query *common.SystemQuery) error {
	return ss.dummyCommandReply(name, common.CmdSystemInform, query)
}

func (ss *SystemStub) Start(name string, query *common.SystemQuery) error {
	return ss.dummyCommandReply(name, common.CmdSystemStart, query)
}

func (ss *SystemStub) Stop(name string, query *common.SystemQuery) error {
	return ss.dummyCommandReply(name, common.CmdSystemStop, query)
}

func (ss *SystemStub) Restart(name string, query *common.SystemQuery) error {
	return ss.dummyCommandReply(name, common.CmdSystemRestart, query)
}
