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
	callback  common.ComplexCallback
	config    *config.AppConfig
	storage   *dbase.DBaseStore
	log       *core.LogAgent
	deviceMap map[string]*SystemDevice
}

func (dr *DeviceRouter) initRouter(log *core.LogAgent, config *config.AppConfig, callback common.ComplexCallback) {
	dr.log = log
	dr.config = config
	dr.callback = callback
	dr.deviceMap = make(map[string]*SystemDevice)
	dr.storage = dbase.GetNewDBaseStore(config.Storage)
}

func (dr *DeviceRouter) cleanupRouter() {
	for name, obj := range dr.deviceMap {
		obj.StopDeviceLoop()
		delete(dr.deviceMap, name)
	}
	dr.storage.Close()
}

func (dr *DeviceRouter) startupRouter() {
	dr.storage.Open()
	for _, cfg := range dr.config.Devices {
		obj, err := dr.createSystemDevice(cfg)
		if err == nil {
			obj.StartDeviceLoop()
			dr.deviceMap[cfg.DevName] = obj
		}
	}
}

//func (dr *DeviceRouter) getDeviceConfig(name string) *config.DeviceConfig {
//	for _, cfg := range dr.config.Devices {
//		if cfg.DevName == name {
//			return cfg
//		}
//	}
//	return nil
//}

func (dr *DeviceRouter) createSystemDevice(cfg *config.DeviceConfig) (*SystemDevice, error) {
	if cfg == nil {
		return nil, errors.New("device config is not set")
	}
	if cfg.DevName == "" {
		return nil, errors.New("device name is not set in device config")
	}
	ctx := &driver.Context{
		DevName: cfg.DevName,
		Complex: dr.callback,
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

func (dr *DeviceRouter) getSystemDevice(name string) *SystemDevice {
	if name == "" {
		return nil
	}
	obj, ok := dr.deviceMap[name]
	if ok {
		return obj
	}
	return nil
}

func (dr *DeviceRouter) delSystemDevice(name string) {
	if name == "" {
		return
	}
	obj, ok := dr.deviceMap[name]
	if ok {
		obj.StopDeviceLoop()
	}
	delete(dr.deviceMap, name)
}

// Implementation of common.ComplexManager

func (dr *DeviceRouter) GetSystemManager() common.SystemManager {
	return dr
}
func (dr *DeviceRouter) GetDeviceManager() common.DeviceManager {
	return dr
}
func (dr *DeviceRouter) GetPrinterManager() common.PrinterManager {
	return dr
}
func (dr *DeviceRouter) GetReaderManager() common.ReaderManager {
	return dr
}
func (dr *DeviceRouter) GetPinPadManager() common.PinPadManager {
	return dr
}
func (dr *DeviceRouter) GetValidatorManager() common.ValidatorManager {
	return dr
}

// Implementation of common.SystemManager

func (dr *DeviceRouter) Terminate(name string, query *common.SystemQuery) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdSystemTerminate, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}

func (dr *DeviceRouter) SysInform(name string, query *common.SystemQuery) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdSystemInform, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}

func (dr *DeviceRouter) SysStart(name string, query *common.SystemConfig) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdSystemStart, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}

func (dr *DeviceRouter) SysStop(name string, query *common.SystemQuery) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdSystemStop, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}

func (dr *DeviceRouter) SysRestart(name string, query *common.SystemConfig) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdSystemRestart, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}

// Implementation of common.DeviceManager

func (dr *DeviceRouter) Cancel(name string, query *common.DeviceQuery) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdDeviceCancel, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (dr *DeviceRouter) Reset(name string, query *common.DeviceQuery) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdDeviceReset, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (dr *DeviceRouter) Status(name string, query *common.DeviceQuery) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdDeviceStatus, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (dr *DeviceRouter) RunAction(name string, query *common.DeviceQuery) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdRunAction, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (dr *DeviceRouter) StopAction(name string, query *common.DeviceQuery) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdStopAction, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}

// Implementation of common.PrinterManager

func (dr *DeviceRouter) InitPrinter(name string, query *common.PrinterSetup) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdInitPrinter, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (dr *DeviceRouter) PrintText(name string, query *common.PrinterQuery) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdPrintText, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}

// Implementation of common.ReaderManager

func (dr *DeviceRouter) EnterCard(name string, query *common.DeviceQuery) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdEnterCard, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (dr *DeviceRouter) EjectCard(name string, query *common.DeviceQuery) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdEjectCard, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (dr *DeviceRouter) CaptureCard(name string, query *common.DeviceQuery) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdCaptureCard, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (dr *DeviceRouter) ReadCard(name string, query *common.DeviceQuery) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdReadCard, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (dr *DeviceRouter) ChipGetATR(name string, query *common.DeviceQuery) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdChipGetATR, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (dr *DeviceRouter) ChipPowerOff(name string, query *common.DeviceQuery) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdChipPowerOff, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (dr *DeviceRouter) ChipCommand(name string, query *common.ReaderChipQuery) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdChipCommand, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}

// Implementation of common.ValidatorManager

func (dr *DeviceRouter) InitValidator(name string, query *common.ValidatorQuery) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdInitValidator, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (dr *DeviceRouter) DoValidate(name string, query *common.ValidatorQuery) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdDoValidate, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (dr *DeviceRouter) NoteAccept(name string, query *common.ValidatorQuery) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdNoteAccept, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (dr *DeviceRouter) NoteReturn(name string, query *common.ValidatorQuery) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdNoteReturn, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (dr *DeviceRouter) StopValidate(name string, query *common.ValidatorQuery) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdStopValidate, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (dr *DeviceRouter) CheckValidator(name string, query *common.ValidatorQuery) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdCheckValidator, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (dr *DeviceRouter) ClearValidator(name string, query *common.ValidatorQuery) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdClearValidator, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}

// Implementation of common.PinPadManager

func (dr *DeviceRouter) ReadPIN(name string, query *common.ReaderPinQuery) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdReadPIN, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (dr *DeviceRouter) LoadMasterKey(name string, query *common.ReaderPinQuery) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdLoadMasterKey, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (dr *DeviceRouter) LoadWorkKey(name string, query *common.ReaderPinQuery) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdLoadWorkKey, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (dr *DeviceRouter) TestMasterKey(name string, query *common.ReaderPinQuery) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdTestMasterKey, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (dr *DeviceRouter) TestWorkKey(name string, query *common.ReaderPinQuery) error {
	device := dr.getSystemDevice(name)
	if device != nil {
		return device.RunCommand(common.CmdTestWorkKey, name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
