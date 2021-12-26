package router

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
)

type DeviceRouter struct {
	config     config.AppConfig
	handlerMap map[string]*DeviceItem
	log        *core.LogAgent
}

func (hr *DeviceRouter) initRouter(config config.AppConfig) {
	hr.config = config
	hr.handlerMap = make(map[string]*DeviceItem)
	hr.log = core.GetLogAgent(core.LogLevelTrace, "Router")
}

func (hr *DeviceRouter) terminateAll() {
	for name, obj := range hr.handlerMap {
		_ = obj.Terminate(name, &common.SystemQuery{})
	}
}

func (hr *DeviceRouter) cleanupRouter() {
	for name, obj := range hr.handlerMap {
		obj.StopDeviceLoop()
		delete(hr.handlerMap, name)
	}
}

func (hr *DeviceRouter) getDeviceConfig(name string) *config.DeviceConfig {
	for _, cfg := range hr.config.Devices {
		if cfg.DevName == name {
			return cfg
		}
	}
	return nil
}

//
func (hr *DeviceRouter) createNewDevice(name string) *DeviceItem {
	if name == "" {
		return nil
	}
	cfg := hr.getDeviceConfig(name)
	obj := NewDeviceItem(cfg)
	obj.StartDeviceLoop()
	hr.handlerMap[name] = obj
	return obj
}

func (hr *DeviceRouter) getDeviceHandler(name string) *DeviceItem {
	if name == "" {
		return nil
	}
	obj, ok := hr.handlerMap[name]
	if ok {
		return obj
	}
	return nil
}

func (hr *DeviceRouter) delDeviceHandler(name string) {
	if name == "" {
		return
	}
	obj, ok := hr.handlerMap[name]
	if ok {
		obj.StopDeviceLoop()
	}
	delete(hr.handlerMap, name)
}

// Implementation of common.SystemCallback
func (hr *DeviceRouter) SystemReply(name string, reply *common.SystemReply) error {
	//handler := hr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.SystemReply(name, reply)
	//}
	return nil
}
func (hr *DeviceRouter) SystemHealth(name string, reply *common.SystemHealth) error {
	//handler := hr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.SystemHealth(name, reply)
	//}
	return nil
}

// Implementation of common.DeviceCallback
func (hr *DeviceRouter) DeviceReply(name string, reply *common.DeviceReply) error {
	//handler := hr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.DeviceReply(name, reply)
	//}
	return nil
}
func (hr *DeviceRouter) ExecuteError(name string, reply *common.DeviceError) error {
	//handler := hr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.ExecuteError(name, reply)
	//}
	return nil
}
func (hr *DeviceRouter) StateChanged(name string, reply *common.DeviceState) error {
	//handler := hr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.StateChanged(name, reply)
	//}
	return nil
}
func (hr *DeviceRouter) ActionPrompt(name string, reply *common.DevicePrompt) error {
	//handler := hr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.ActionPrompt(name, reply)
	//}
	return nil
}
func (hr *DeviceRouter) ReaderReturn(name string, reply *common.DeviceInform) error {
	//handler := hr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.ReaderReturn(name, reply)
	//}
	return nil
}

// Implementation of common.PrinterCallback
func (hr *DeviceRouter) PrinterProgress(name string, reply *common.PrinterProgress) error {
	//handler := hr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.PrinterProgress(name, reply)
	//}
	return nil
}

// Implementation of common.ReaderCallback
func (hr *DeviceRouter) CardPosition(name string, reply *common.ReaderCardPos) error {
	//handler := hr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.CardPosition(name, reply)
	//}
	return nil
}
func (hr *DeviceRouter) CardDescription(name string, reply *common.ReaderCardInfo) error {
	//handler := hr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.CardDescription(name, reply)
	//}
	return nil
}
func (hr *DeviceRouter) ChipResponse(name string, reply *common.ReaderChipReply) error {
	//handler := hr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.ChipResponse(name, reply)
	//}
	return nil
}

// Implementation of common.ValidatorCallback
func (hr *DeviceRouter) NoteAccepted(name string, reply *common.ValidatorAccept) error {
	//handler := hr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.NoteAccepted(name, reply)
	//}
	return nil
}
func (hr *DeviceRouter) CashIsStored(name string, reply *common.ValidatorAccept) error {
	//handler := hr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.CashIsStored(name, reply)
	//}
	return nil
}
func (hr *DeviceRouter) CashReturned(name string, reply *common.ValidatorAccept) error {
	//handler := hr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.CashReturned(name, reply)
	//}
	return nil
}
func (hr *DeviceRouter) ValidatorStore(name string, reply *common.ValidatorStore) error {
	//handler := hr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.ValidatorStore(name, reply)
	//}
	return nil
}

// Implementation of common.ReaderCallback
func (hr *DeviceRouter) PinPadReply(name string, reply *common.ReaderPinReply) error {
	//handler := hr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.PinPadReply(name, reply)
	//}
	return nil
}
