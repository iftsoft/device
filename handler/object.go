package handler

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
	"github.com/iftsoft/device/proxy"
	"time"
)

type ObjectProxy struct {
	server duplex.ServerManager
	system *proxy.SystemServer
	device *proxy.DeviceServer
	objMap map[string]*ObjectState
	log    *core.LogAgent
}

func NewObjectProxy() *ObjectProxy {
	op := ObjectProxy{
		server: nil,
		system: proxy.NewSystemServer(),
		device: proxy.NewDeviceServer(),
		objMap: make(map[string]*ObjectState),
		log:    core.GetLogAgent(core.LogLevelTrace, "Object"),
	}
	return &op
}

func (op *ObjectProxy) Init(server duplex.ServerManager) {
	op.server = server
	op.system.Init(op.server, op, op.log)
	op.device.Init(op.server, op, op.log)
}

func (op *ObjectProxy) GetObjectState(name string) *ObjectState {
	if name == "" {
		return nil
	}
	obj, ok := op.objMap[name]
	if ok {
		return obj
	}
	obj = NewObjectState()
	obj.Init(name, op.log)
	op.objMap[name] = obj
	return obj
}

func (op *ObjectProxy) DelObjectState(name string) {
	if name == "" {
		return
	}
	delete(op.objMap, name)
}

// Implementation of duplex.ClientManager
func (op *ObjectProxy) OnClientStarted(name string) {
	if name == "" {
		return
	}
	op.log.Trace("ObjectProxy.OnClientStarted dev:%s", name)
	obj := op.GetObjectState(name)
	go op.runDeviceTask(obj.DevName)
}

func (op *ObjectProxy) runDeviceTask(name string) {
	query := &common.SystemQuery{}
	query.DevName = name
	value := &common.DeviceQuery{}
	time.Sleep(time.Millisecond)
	op.log.Trace("ObjectProxy.runDeviceTask dev:%s", name)
	_ = op.Inform(name, query)
	_ = op.Config(name, query)
	_ = op.Status(name, value)
}

func (op *ObjectProxy) OnClientStopped(name string) {
	if name == "" {
		return
	}
	op.log.Trace("ObjectProxy.OnClientStopped dev:%s", name)
}

// Implementation of common.SystemCallback
func (op *ObjectProxy) SystemReply(name string, reply *common.SystemReply) error {
	object := op.GetObjectState(name)
	if object != nil {
		return object.SystemReply(name, reply)
	}
	return nil
}

// Implementation of common.SystemManager
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

// Implementation of common.DeviceCallback
func (op *ObjectProxy) DeviceReply(name string, reply *common.DeviceReply) error {
	object := op.GetObjectState(name)
	if object != nil {
		return object.DeviceReply(name, reply)
	}
	return nil
}

func (op *ObjectProxy) ExecuteError(name string, reply *common.DeviceError) error {
	object := op.GetObjectState(name)
	if object != nil {
		return object.ExecuteError(name, reply)
	}
	return nil
}

func (op *ObjectProxy) StateChanged(name string, reply *common.DeviceState) error {
	object := op.GetObjectState(name)
	if object != nil {
		return object.StateChanged(name, reply)
	}
	return nil
}

func (op *ObjectProxy) ActionPrompt(name string, reply *common.DevicePrompt) error {
	object := op.GetObjectState(name)
	if object != nil {
		return object.ActionPrompt(name, reply)
	}
	return nil
}

// Implementation of common.DeviceManager
func (op *ObjectProxy) Cancel(name string, query *common.DeviceQuery) error {
	return op.device.SendDeviceCommand(name, common.CmdDeviceCancel, query)
}

func (op *ObjectProxy) Reset(name string, query *common.DeviceQuery) error {
	return op.device.SendDeviceCommand(name, common.CmdDeviceReset, query)
}

func (op *ObjectProxy) Status(name string, query *common.DeviceQuery) error {
	return op.device.SendDeviceCommand(name, common.CmdDeviceStatus, query)
}
