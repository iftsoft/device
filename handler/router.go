package handler

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"sync"
)

type HandlerRouter struct {
	config     config.HandlerList
	handlerMap map[string]*DeviceHandler
	log        *core.LogAgent
	wg         sync.WaitGroup
}

func (hr *HandlerRouter) initRouter(config config.HandlerList) {
	hr.config     = config
	hr.handlerMap = make(map[string]*DeviceHandler)
	hr.log        = core.GetLogAgent(core.LogLevelTrace, "Router")
}

func (hr *HandlerRouter) terminateAll(hp *HandlerProxy) {
	for name, _ := range hr.handlerMap {
		hp.Terminate(name, &common.SystemQuery{})
	}
}

func (hr *HandlerRouter) cleanupRouter() {
	for name, obj := range hr.handlerMap {
		obj.StopObject()
		delete(hr.handlerMap, name)
	}
	hr.wg.Wait()
}


//
func (hr *HandlerRouter) createNewHandler(name string) *DeviceHandler {
	if name == "" {
		return nil
	}
	obj := NewDeviceHandler(name, hr.log)
	obj.StartObject(&hr.wg)
	hr.handlerMap[name] = obj
	return obj
}


func (hr *HandlerRouter) getDeviceHandler(name string) *DeviceHandler {
	if name == "" {
		return nil
	}
	obj, ok := hr.handlerMap[name]
	if ok {
		return obj
	}
	return nil
}

func (hr *HandlerRouter) delDeviceHandler(name string) {
	if name == "" {
		return
	}
	obj, ok := hr.handlerMap[name]
	if ok {
		obj.StopObject()
	}
	delete(hr.handlerMap, name)
}


func (hr *HandlerRouter) onClientStarted(name string) *DeviceHandler {
	hr.log.Trace("HandlerRouter.OnClientStarted device:%s", name)
	handler := hr.getDeviceHandler(name)
	if handler == nil {
		handler = hr.createNewHandler(name)
	}
	return handler
}

func (hr *HandlerRouter) onClientStopped(name string) {
	hr.log.Trace("HandlerRouter.OnClientStopped dev:%s", name)
	handler := hr.getDeviceHandler(name)
	if handler != nil {
		handler.OnClientStopped(name)
	}
}


// Implementation of common.SystemCallback
func (hr *HandlerRouter) SystemReply(name string, reply *common.SystemReply) error {
	handler := hr.getDeviceHandler(name)
	if handler != nil {
		return handler.SystemReply(name, reply)
	}
	return nil
}
func (hr *HandlerRouter) SystemHealth(name string, reply *common.SystemHealth) error {
	handler := hr.getDeviceHandler(name)
	if handler != nil {
		return handler.SystemHealth(name, reply)
	}
	return nil
}

// Implementation of common.DeviceCallback
func (hr *HandlerRouter) DeviceReply(name string, reply *common.DeviceReply) error {
	handler := hr.getDeviceHandler(name)
	if handler != nil {
		return handler.DeviceReply(name, reply)
	}
	return nil
}
func (hr *HandlerRouter) ExecuteError(name string, reply *common.DeviceError) error {
	handler := hr.getDeviceHandler(name)
	if handler != nil {
		return handler.ExecuteError(name, reply)
	}
	return nil
}
func (hr *HandlerRouter) StateChanged(name string, reply *common.DeviceState) error {
	handler := hr.getDeviceHandler(name)
	if handler != nil {
		return handler.StateChanged(name, reply)
	}
	return nil
}
func (hr *HandlerRouter) ActionPrompt(name string, reply *common.DevicePrompt) error {
	handler := hr.getDeviceHandler(name)
	if handler != nil {
		return handler.ActionPrompt(name, reply)
	}
	return nil
}
func (hr *HandlerRouter) ReaderReturn(name string, reply *common.DeviceInform) error {
	handler := hr.getDeviceHandler(name)
	if handler != nil {
		return handler.ReaderReturn(name, reply)
	}
	return nil
}

// Implementation of common.PrinterCallback
func (hr *HandlerRouter) PrinterProgress(name string, reply *common.PrinterProgress) error {
	handler := hr.getDeviceHandler(name)
	if handler != nil {
		return handler.PrinterProgress(name, reply)
	}
	return nil
}

// Implementation of common.ReaderCallback
func (hr *HandlerRouter) CardPosition(name string, reply *common.ReaderCardPos) error {
	handler := hr.getDeviceHandler(name)
	if handler != nil {
		return handler.CardPosition(name, reply)
	}
	return nil
}
func (hr *HandlerRouter) CardDescription(name string, reply *common.ReaderCardInfo) error {
	handler := hr.getDeviceHandler(name)
	if handler != nil {
		return handler.CardDescription(name, reply)
	}
	return nil
}
func (hr *HandlerRouter) ChipResponse(name string, reply *common.ReaderChipReply) error {
	handler := hr.getDeviceHandler(name)
	if handler != nil {
		return handler.ChipResponse(name, reply)
	}
	return nil
}

// Implementation of common.ValidatorCallback
func (hr *HandlerRouter) NoteAccepted(name string, reply *common.ValidatorAccept) error {
	handler := hr.getDeviceHandler(name)
	if handler != nil {
		return handler.NoteAccepted(name, reply)
	}
	return nil
}
func (hr *HandlerRouter) CashIsStored(name string, reply *common.ValidatorAccept) error {
	handler := hr.getDeviceHandler(name)
	if handler != nil {
		return handler.CashIsStored(name, reply)
	}
	return nil
}
func (hr *HandlerRouter) CashReturned(name string, reply *common.ValidatorAccept) error {
	handler := hr.getDeviceHandler(name)
	if handler != nil {
		return handler.CashReturned(name, reply)
	}
	return nil
}
func (hr *HandlerRouter) ValidatorStore(name string, reply *common.ValidatorStore) error {
	handler := hr.getDeviceHandler(name)
	if handler != nil {
		return handler.ValidatorStore(name, reply)
	}
	return nil
}

// Implementation of common.ReaderCallback
func (hr *HandlerRouter) PinPadReply(name string, reply *common.ReaderPinReply) error {
	handler := hr.getDeviceHandler(name)
	if handler != nil {
		return handler.PinPadReply(name, reply)
	}
	return nil
}
