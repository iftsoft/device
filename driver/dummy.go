package driver

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
)

type DummyDriver struct {
	config *config.DeviceConfig
	device common.DeviceCallback
	log    *core.LogAgent
}

func NewDummyDriver() *DummyDriver {
	dd := DummyDriver{
		config: nil,
		device: nil,
		log:    core.GetLogAgent(core.LogLevelTrace, "Dummy"),
	}
	return &dd
}

func (dd *DummyDriver) InitDevice(manager interface{}) error {
	device, okDev := manager.(common.DeviceCallback)
	if okDev {
		dd.device = device
	}
	return nil
}

func (dd *DummyDriver) StartDevice(cfg *config.DeviceConfig) error {
	dd.config = cfg
	return nil
}
func (dd *DummyDriver) DeviceTimer(unix int64) error {
	return nil
}
func (dd *DummyDriver) StopDevice() error {
	return nil
}
func (dd *DummyDriver) ClearDevice() error {
	return nil
}

// Implementation of common.DeviceManager
//
func (dd *DummyDriver) Cancel(name string, query *common.DeviceQuery) error {
	return dd.dummyCommandReply(name, common.CmdDeviceCancel, query)
}

func (dd *DummyDriver) Reset(name string, query *common.DeviceQuery) error {
	return dd.dummyCommandReply(name, common.CmdDeviceReset, query)
}

func (dd *DummyDriver) Status(name string, query *common.DeviceQuery) error {
	return dd.dummyCommandReply(name, common.CmdDeviceStatus, query)
}

func (dd *DummyDriver) dummyCommandReply(name string, cmd string, query interface{}) error {
	if dd.log != nil {
		dd.log.Debug("DummyDriver dev:%s run cmd:%s", name, cmd)
	}
	reply := &common.DeviceReply{}
	reply.Command = cmd
	reply.DevState = common.DevStateUndefined
	var err error
	if dd.device != nil {
		err = dd.device.DeviceReply(name, reply)
	}
	return err
}
