package general

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
)

type DeviceStub struct {
	callback common.DeviceCallback
	log      *core.LogAgent
}

func NewDeviceStub() *DeviceStub {
	ss := DeviceStub{
		callback: nil,
		log:      nil,
	}
	return &ss
}

func (ss *DeviceStub) Init(callback common.DeviceCallback, log *core.LogAgent) {
	ss.log = log
	ss.callback = callback
}

func (ss *DeviceStub) dummyCommandReply(name string, cmd string, query *common.DeviceQuery) error {
	if ss.log != nil {
		ss.log.Trace("SystemStub dev:%s run cmd:%s, timeout:%d", name, cmd, query.Timeout)
	}
	reply := &common.DeviceReply{}
	reply.Command = cmd
	reply.DevState = common.DevStateUndefined
	var err error
	if ss.callback != nil {
		err = ss.callback.CommandReply(name, reply)
	}
	return err
}

// Implementation of common.DeviceManager
//
func (ss *DeviceStub) Cancel(name string, query *common.DeviceQuery) error {
	return ss.dummyCommandReply(name, common.CmdDeviceCancel, query)
}

func (ss *DeviceStub) Reset(name string, query *common.DeviceQuery) error {
	return ss.dummyCommandReply(name, common.CmdDeviceReset, query)
}

func (ss *DeviceStub) Status(name string, query *common.DeviceQuery) error {
	return ss.dummyCommandReply(name, common.CmdDeviceStatus, query)
}
