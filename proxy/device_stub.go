package proxy

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

func (ss *DeviceStub) dummyCommandReply(name string, cmd string, query interface{}) error {
	if ss.log != nil {
		ss.log.Debug("SystemStub dev:%s run cmd:%s", name, cmd)
	}
	reply := &common.DeviceReply{}
	reply.Command = cmd
	reply.DevState = common.DevStateUndefined
	var err error
	if ss.callback != nil {
		err = ss.callback.DeviceReply(name, reply)
	}
	return err
}
