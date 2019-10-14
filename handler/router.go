package handler

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"sync"
)

type HandlerRouter struct {
	config config.HandlerList
	proxy  interface{}
	objMap map[string]*DeviceHandler
	log    *core.LogAgent
	wg     sync.WaitGroup
}


func (or *HandlerRouter) Cleanup() {
	for name, obj := range or.objMap {
		obj.StopObject()
		delete(or.objMap, name)
	}
	or.wg.Wait()
}

func (or *HandlerRouter) GetDeviceHandler(name string) *DeviceHandler {
	if name == "" {
		return nil
	}
	obj, ok := or.objMap[name]
	if ok {
		return obj
	}
	obj = NewDeviceHandler(name, or.log)
	obj.StartObject(&or.wg)
	or.objMap[name] = obj
	return obj
}

func (or *HandlerRouter) DelDeviceHandler(name string) {
	if name == "" {
		return
	}
	obj, ok := or.objMap[name]
	if ok {
		obj.StopObject()
	}
	delete(or.objMap, name)
}

// Implementation of duplex.ClientManager
func (or *HandlerRouter) OnClientStarted(name string) {
	if name == "" {
		return
	}
	or.log.Trace("HandlerProxy.OnClientStarted dev:%s", name)
	obj := or.GetDeviceHandler(name)
	obj.OnClientStarted(name)
}

func (or *HandlerRouter) OnClientStopped(name string) {
	if name == "" {
		return
	}
	or.log.Trace("HandlerProxy.OnClientStopped dev:%s", name)
	obj := or.GetDeviceHandler(name)
	obj.OnClientStopped(name)
}

// Implementation of common.SystemCallback
func (or *HandlerRouter) SystemReply(name string, reply *common.SystemReply) error {
	object := or.GetDeviceHandler(name)
	if object != nil {
		return object.SystemReply(name, reply)
	}
	return nil
}
func (or *HandlerRouter) SystemHealth(name string, reply *common.SystemHealth) error {
	object := or.GetDeviceHandler(name)
	if object != nil {
		return object.SystemHealth(name, reply)
	}
	return nil
}

// Implementation of common.DeviceCallback
func (or *HandlerRouter) DeviceReply(name string, reply *common.DeviceReply) error {
	object := or.GetDeviceHandler(name)
	if object != nil {
		return object.DeviceReply(name, reply)
	}
	return nil
}
func (or *HandlerRouter) ExecuteError(name string, reply *common.DeviceError) error {
	object := or.GetDeviceHandler(name)
	if object != nil {
		return object.ExecuteError(name, reply)
	}
	return nil
}
func (or *HandlerRouter) StateChanged(name string, reply *common.DeviceState) error {
	object := or.GetDeviceHandler(name)
	if object != nil {
		return object.StateChanged(name, reply)
	}
	return nil
}
func (or *HandlerRouter) ActionPrompt(name string, reply *common.DevicePrompt) error {
	object := or.GetDeviceHandler(name)
	if object != nil {
		return object.ActionPrompt(name, reply)
	}
	return nil
}
func (or *HandlerRouter) ReaderReturn(name string, reply *common.DeviceInform) error {
	object := or.GetDeviceHandler(name)
	if object != nil {
		return object.ReaderReturn(name, reply)
	}
	return nil
}

// Implementation of common.PrinterCallback
func (or *HandlerRouter) PrinterProgress(name string, reply *common.PrinterProgress) error {
	object := or.GetDeviceHandler(name)
	if object != nil {
		return object.PrinterProgress(name, reply)
	}
	return nil
}

// Implementation of common.ReaderCallback
func (or *HandlerRouter) CardPosition(name string, reply *common.ReaderCardPos) error {
	object := or.GetDeviceHandler(name)
	if object != nil {
		return object.CardPosition(name, reply)
	}
	return nil
}
func (or *HandlerRouter) CardDescription(name string, reply *common.ReaderCardInfo) error {
	object := or.GetDeviceHandler(name)
	if object != nil {
		return object.CardDescription(name, reply)
	}
	return nil
}
func (or *HandlerRouter) ChipResponse(name string, reply *common.ReaderChipReply) error {
	object := or.GetDeviceHandler(name)
	if object != nil {
		return object.ChipResponse(name, reply)
	}
	return nil
}

// Implementation of common.ValidatorCallback
func (or *HandlerRouter) NoteAccepted(name string, reply *common.ValidatorAccept) error {
	object := or.GetDeviceHandler(name)
	if object != nil {
		return object.NoteAccepted(name, reply)
	}
	return nil
}
func (or *HandlerRouter) CashIsStored(name string, reply *common.ValidatorAccept) error {
	object := or.GetDeviceHandler(name)
	if object != nil {
		return object.CashIsStored(name, reply)
	}
	return nil
}
func (or *HandlerRouter) CashReturned(name string, reply *common.ValidatorAccept) error {
	object := or.GetDeviceHandler(name)
	if object != nil {
		return object.CashReturned(name, reply)
	}
	return nil
}
func (or *HandlerRouter) ValidatorStore(name string, reply *common.ValidatorStore) error {
	object := or.GetDeviceHandler(name)
	if object != nil {
		return object.ValidatorStore(name, reply)
	}
	return nil
}

// Implementation of common.ReaderCallback
func (or *HandlerRouter) PinPadReply(name string, reply *common.ReaderPinReply) error {
	object := or.GetDeviceHandler(name)
	if object != nil {
		return object.PinPadReply(name, reply)
	}
	return nil
}
