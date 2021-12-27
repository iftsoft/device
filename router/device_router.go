package router

import (
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/dbase"
	"github.com/iftsoft/device/driver"
)

type DeviceRouter struct {
	config     config.AppConfig
	storage    *dbase.DBaseStore
	log        *core.LogAgent
	handlerMap map[string]*SystemDevice
}

func (dr *DeviceRouter) initRouter(config config.AppConfig) {
	dr.config = config
	dr.handlerMap = make(map[string]*SystemDevice)
	dr.log = core.GetLogAgent(core.LogLevelTrace, "Router")
}

func (dr *DeviceRouter) terminateAll() {
	for name, obj := range dr.handlerMap {
		_ = obj.Terminate(name, &common.SystemQuery{})
	}
}

func (dr *DeviceRouter) cleanupRouter() {
	for name, obj := range dr.handlerMap {
		obj.StopDeviceLoop()
		delete(dr.handlerMap, name)
	}
}

func (dr *DeviceRouter) getDeviceConfig(name string) *config.DeviceConfig {
	for _, cfg := range dr.config.Devices {
		if cfg.DevName == name {
			return cfg
		}
	}
	return nil
}

func (dr *DeviceRouter) createSystemDevice(cfg *config.DeviceConfig) (*SystemDevice, error) {
	if cfg == nil {
		return nil, errors.New("device config is not set")
	}
	if cfg.DevName == "" {
		return nil, errors.New("device name is not set in device config")
	}
	ctx := &driver.Context{
		DevName: cfg.DevName,
		Complex: dr,
		Storage: dr.storage,
		Config:  cfg,
	}
	drv := driver.GetDeviceDriver(cfg.Common.Model)
	if drv == nil {
		return nil, errors.New("device model is not supported")
	}
	obj := NewSystemDevice(cfg.DevName)
	obj.InitDevice(drv, ctx)

	return obj, nil
}

//
//func (dr *DeviceRouter) createNewDevice(name string) *SystemDevice {
//	if name == "" {
//		return nil
//	}
//	cfg := dr.getDeviceConfig(name)
//	obj := NewSystemDevice(cfg)
//
//	obj.StartDeviceLoop()
//	dr.handlerMap[name] = obj
//	return obj
//}

func (dr *DeviceRouter) getDeviceHandler(name string) *SystemDevice {
	if name == "" {
		return nil
	}
	obj, ok := dr.handlerMap[name]
	if ok {
		return obj
	}
	return nil
}

func (dr *DeviceRouter) delDeviceHandler(name string) {
	if name == "" {
		return
	}
	obj, ok := dr.handlerMap[name]
	if ok {
		obj.StopDeviceLoop()
	}
	delete(dr.handlerMap, name)
}

// Implementation of common.ComplexCallback

func (dr *DeviceRouter) GetSystemCallback() common.SystemCallback {
	return dr
}
func (dr *DeviceRouter) GetDeviceCallback() common.DeviceCallback {
	return dr
}
func (dr *DeviceRouter) GetPrinterCallback() common.PrinterCallback {
	return dr
}
func (dr *DeviceRouter) GetReaderCallback() common.ReaderCallback {
	return dr
}
func (dr *DeviceRouter) GetPinPadCallback() common.PinPadCallback {
	return dr
}
func (dr *DeviceRouter) GetValidatorCallback() common.ValidatorCallback {
	return dr
}

// Implementation of common.SystemCallback

func (dr *DeviceRouter) SystemReply(name string, reply *common.SystemReply) error {
	//handler := dr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.SystemReply(name, reply)
	//}
	return nil
}
func (dr *DeviceRouter) SystemHealth(name string, reply *common.SystemHealth) error {
	//handler := dr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.SystemHealth(name, reply)
	//}
	return nil
}

// Implementation of common.DeviceCallback
func (dr *DeviceRouter) DeviceReply(name string, reply *common.DeviceReply) error {
	//handler := dr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.DeviceReply(name, reply)
	//}
	return nil
}
func (dr *DeviceRouter) ExecuteError(name string, reply *common.DeviceError) error {
	//handler := dr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.ExecuteError(name, reply)
	//}
	return nil
}
func (dr *DeviceRouter) StateChanged(name string, reply *common.DeviceState) error {
	//handler := dr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.StateChanged(name, reply)
	//}
	return nil
}
func (dr *DeviceRouter) ActionPrompt(name string, reply *common.DevicePrompt) error {
	//handler := dr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.ActionPrompt(name, reply)
	//}
	return nil
}
func (dr *DeviceRouter) ReaderReturn(name string, reply *common.DeviceInform) error {
	//handler := dr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.ReaderReturn(name, reply)
	//}
	return nil
}

// Implementation of common.PrinterCallback
func (dr *DeviceRouter) PrinterProgress(name string, reply *common.PrinterProgress) error {
	//handler := dr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.PrinterProgress(name, reply)
	//}
	return nil
}

// Implementation of common.ReaderCallback
func (dr *DeviceRouter) CardPosition(name string, reply *common.ReaderCardPos) error {
	//handler := dr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.CardPosition(name, reply)
	//}
	return nil
}
func (dr *DeviceRouter) CardDescription(name string, reply *common.ReaderCardInfo) error {
	//handler := dr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.CardDescription(name, reply)
	//}
	return nil
}
func (dr *DeviceRouter) ChipResponse(name string, reply *common.ReaderChipReply) error {
	//handler := dr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.ChipResponse(name, reply)
	//}
	return nil
}

// Implementation of common.ValidatorCallback
func (dr *DeviceRouter) NoteAccepted(name string, reply *common.ValidatorAccept) error {
	//handler := dr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.NoteAccepted(name, reply)
	//}
	return nil
}
func (dr *DeviceRouter) CashIsStored(name string, reply *common.ValidatorAccept) error {
	//handler := dr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.CashIsStored(name, reply)
	//}
	return nil
}
func (dr *DeviceRouter) CashReturned(name string, reply *common.ValidatorAccept) error {
	//handler := dr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.CashReturned(name, reply)
	//}
	return nil
}
func (dr *DeviceRouter) ValidatorStore(name string, reply *common.ValidatorStore) error {
	//handler := dr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.ValidatorStore(name, reply)
	//}
	return nil
}

// Implementation of common.ReaderCallback
func (dr *DeviceRouter) PinPadReply(name string, reply *common.ReaderPinReply) error {
	//handler := dr.getDeviceHandler(name)
	//if handler != nil {
	//	return handler.PinPadReply(name, reply)
	//}
	return nil
}
