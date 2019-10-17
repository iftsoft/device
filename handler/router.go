package handler

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
	"sync"
)

type HandlerRouter struct {
	config     config.HandlerList
	proxy      interface{}
	handlerMap map[string]*DeviceHandler
	reflexMap  map[string]ReflexCreator
	reflexLog  *core.LogAgent
	log        *core.LogAgent
	wg         sync.WaitGroup
}

func (hr *HandlerRouter) InitRouter(cfg config.HandlerList, proxy interface{}) {
	hr.config     = cfg
	hr.proxy      = proxy
	hr.handlerMap = make(map[string]*DeviceHandler)
	hr.reflexMap  = make(map[string]ReflexCreator)
	hr.reflexLog  = core.GetLogAgent(core.LogLevelTrace, "Reflex")
	hr.log        = core.GetLogAgent(core.LogLevelTrace, "Router")
}

func (hr *HandlerRouter) Cleanup() {
	for name, obj := range hr.handlerMap {
		obj.StopObject()
		delete(hr.handlerMap, name)
	}
	hr.wg.Wait()
}

func (hr *HandlerRouter) RegisterReflexFactory(fact ReflexCreator) {
	if fact == nil { return }
	info := fact.GetReflexInfo()
	if info == nil { return }
	hr.log.Debug("HandlerProxy.RegisterReflexFactory for reflex:%s", info.ReflexName)
	hr.reflexMap[info.ReflexName] = fact
}

func (hr *HandlerRouter) CreateNewHandler(name string) *DeviceHandler {
	if name == "" {
		return nil
	}
	obj := NewDeviceHandler(name, hr.log)
	obj.StartObject(&hr.wg)
	hr.handlerMap[name] = obj
	return obj
}

func (hr *HandlerRouter) HandlerReflexes(handler *DeviceHandler, greet *duplex.GreetingInfo) {
	// Attach mandatory reflexes
	for refName, factory := range hr.reflexMap {
		info := factory.GetReflexInfo()
		if info.Mandatory && info.IsMatched(greet) {
			hr.log.Debug("HandlerProxy.CreateNewHandler is attaching reflex:%s to device:%s",
				refName, handler.devName)
			err, reflex := factory.CreateReflex(handler.devName, hr.proxy, hr.reflexLog)
			if err == nil {
				err = handler.AttachReflex(refName, reflex)
			}
		}
	}
	return
}

func (hr *HandlerRouter) GetDeviceHandler(name string) *DeviceHandler {
	if name == "" {
		return nil
	}
	obj, ok := hr.handlerMap[name]
	if ok {
		return obj
	}
	return nil
}

func (hr *HandlerRouter) DelDeviceHandler(name string) {
	if name == "" {
		return
	}
	obj, ok := hr.handlerMap[name]
	if ok {
		obj.StopObject()
	}
	delete(hr.handlerMap, name)
}

// Implementation of duplex.ClientManager
func (hr *HandlerRouter) OnClientStarted(name string, info *duplex.GreetingInfo) {
	if name == "" {
		return
	}
	hr.log.Trace("HandlerProxy.OnClientStarted dev:%s, sup:%X, req:%X",
		name, info.Supported, info.Required)
	handler := hr.GetDeviceHandler(name)
	if handler == nil {
		handler = hr.CreateNewHandler(name)
	}
	hr.HandlerReflexes(handler, info)
	handler.OnClientStarted(name)
}

func (hr *HandlerRouter) OnClientStopped(name string) {
	if name == "" {
		return
	}
	hr.log.Trace("HandlerProxy.OnClientStopped dev:%s", name)
	handler := hr.GetDeviceHandler(name)
	if handler != nil {
		handler.OnClientStopped(name)
	}
}

// Implementation of common.SystemCallback
func (hr *HandlerRouter) SystemReply(name string, reply *common.SystemReply) error {
	handler := hr.GetDeviceHandler(name)
	if handler != nil {
		return handler.SystemReply(name, reply)
	}
	return nil
}
func (hr *HandlerRouter) SystemHealth(name string, reply *common.SystemHealth) error {
	handler := hr.GetDeviceHandler(name)
	if handler != nil {
		return handler.SystemHealth(name, reply)
	}
	return nil
}

// Implementation of common.DeviceCallback
func (hr *HandlerRouter) DeviceReply(name string, reply *common.DeviceReply) error {
	handler := hr.GetDeviceHandler(name)
	if handler != nil {
		return handler.DeviceReply(name, reply)
	}
	return nil
}
func (hr *HandlerRouter) ExecuteError(name string, reply *common.DeviceError) error {
	handler := hr.GetDeviceHandler(name)
	if handler != nil {
		return handler.ExecuteError(name, reply)
	}
	return nil
}
func (hr *HandlerRouter) StateChanged(name string, reply *common.DeviceState) error {
	handler := hr.GetDeviceHandler(name)
	if handler != nil {
		return handler.StateChanged(name, reply)
	}
	return nil
}
func (hr *HandlerRouter) ActionPrompt(name string, reply *common.DevicePrompt) error {
	handler := hr.GetDeviceHandler(name)
	if handler != nil {
		return handler.ActionPrompt(name, reply)
	}
	return nil
}
func (hr *HandlerRouter) ReaderReturn(name string, reply *common.DeviceInform) error {
	handler := hr.GetDeviceHandler(name)
	if handler != nil {
		return handler.ReaderReturn(name, reply)
	}
	return nil
}

// Implementation of common.PrinterCallback
func (hr *HandlerRouter) PrinterProgress(name string, reply *common.PrinterProgress) error {
	handler := hr.GetDeviceHandler(name)
	if handler != nil {
		return handler.PrinterProgress(name, reply)
	}
	return nil
}

// Implementation of common.ReaderCallback
func (hr *HandlerRouter) CardPosition(name string, reply *common.ReaderCardPos) error {
	handler := hr.GetDeviceHandler(name)
	if handler != nil {
		return handler.CardPosition(name, reply)
	}
	return nil
}
func (hr *HandlerRouter) CardDescription(name string, reply *common.ReaderCardInfo) error {
	handler := hr.GetDeviceHandler(name)
	if handler != nil {
		return handler.CardDescription(name, reply)
	}
	return nil
}
func (hr *HandlerRouter) ChipResponse(name string, reply *common.ReaderChipReply) error {
	handler := hr.GetDeviceHandler(name)
	if handler != nil {
		return handler.ChipResponse(name, reply)
	}
	return nil
}

// Implementation of common.ValidatorCallback
func (hr *HandlerRouter) NoteAccepted(name string, reply *common.ValidatorAccept) error {
	handler := hr.GetDeviceHandler(name)
	if handler != nil {
		return handler.NoteAccepted(name, reply)
	}
	return nil
}
func (hr *HandlerRouter) CashIsStored(name string, reply *common.ValidatorAccept) error {
	handler := hr.GetDeviceHandler(name)
	if handler != nil {
		return handler.CashIsStored(name, reply)
	}
	return nil
}
func (hr *HandlerRouter) CashReturned(name string, reply *common.ValidatorAccept) error {
	handler := hr.GetDeviceHandler(name)
	if handler != nil {
		return handler.CashReturned(name, reply)
	}
	return nil
}
func (hr *HandlerRouter) ValidatorStore(name string, reply *common.ValidatorStore) error {
	handler := hr.GetDeviceHandler(name)
	if handler != nil {
		return handler.ValidatorStore(name, reply)
	}
	return nil
}

// Implementation of common.ReaderCallback
func (hr *HandlerRouter) PinPadReply(name string, reply *common.ReaderPinReply) error {
	handler := hr.GetDeviceHandler(name)
	if handler != nil {
		return handler.PinPadReply(name, reply)
	}
	return nil
}
