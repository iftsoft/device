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

// Implementation of DeviceDriver interface
func (dd *DummyDriver) InitDevice(manager interface{}) error {
	dd.log.Debug("DummyDriver run cmd:%s", "InitDevice")
	if device, ok := manager.(common.DeviceCallback); ok {
		dd.device = device
	}
	return nil
}

func (dd *DummyDriver) StartDevice(cfg *config.DeviceConfig) error {
	dd.config = cfg
	dd.log.Debug("DummyDriver run cmd:%s", "StartDevice")
	return nil
}
func (dd *DummyDriver) DeviceTimer(unix int64) error {
	dd.log.Debug("DummyDriver run cmd:%s", "DeviceTimer")
	return nil
}
func (dd *DummyDriver) StopDevice() error {
	dd.log.Debug("DummyDriver run cmd:%s", "StopDevice")
	return nil
}
func (dd *DummyDriver) CheckDevice(metrics *common.SystemMetrics) error {
	dd.log.Debug("DummyDriver run cmd:%s", "CheckDevice")
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
