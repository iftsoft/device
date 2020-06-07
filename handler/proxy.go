package handler

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/duplex"
	"github.com/iftsoft/device/proxy"
)

type HandlerProxy struct {
	serverMng    duplex.ServerManager
	systemSrv    *proxy.SystemServer
	deviceSrv    *proxy.DeviceServer
	printerSrv   *proxy.PrinterServer
	readerSrv    *proxy.ReaderServer
	validatorSrv *proxy.ValidatorServer
	pinpadSrv    *proxy.PinPadServer
}


func (hp *HandlerProxy) initProxy() {
	hp.systemSrv    = proxy.NewSystemServer()
	hp.deviceSrv    = proxy.NewDeviceServer()
	hp.printerSrv   = proxy.NewPrinterServer()
	hp.readerSrv    = proxy.NewReaderServer()
	hp.validatorSrv = proxy.NewValidatorServer()
	hp.pinpadSrv    = proxy.NewPinPadServer()
}

func (hp *HandlerProxy) setupProxy(server duplex.ServerManager, hr *HandlerRouter) {
	hp.serverMng = server
	hp.systemSrv.Init(server, hr, hr.log)
	hp.deviceSrv.Init(server, hr, hr.log)
	hp.printerSrv.Init(server, hr, hr.log)
	hp.readerSrv.Init(server, hr, hr.log)
	hp.validatorSrv.Init(server, hr, hr.log)
	hp.pinpadSrv.Init(server, hr, hr.log)
}



// Implementation of common.SystemManager
func (hp *HandlerProxy) Terminate(name string, query *common.SystemQuery) error {
	return hp.systemSrv.SendSystemCommand(name, common.CmdSystemTerminate, query)
}
func (hp *HandlerProxy) SysInform(name string, query *common.SystemQuery) error {
	return hp.systemSrv.SendSystemCommand(name, common.CmdSystemInform, query)
}
func (hp *HandlerProxy) SysStart(name string, query *common.SystemQuery) error {
	return hp.systemSrv.SendSystemCommand(name, common.CmdSystemStart, query)
}
func (hp *HandlerProxy) SysStop(name string, query *common.SystemQuery) error {
	return hp.systemSrv.SendSystemCommand(name, common.CmdSystemStop, query)
}
func (hp *HandlerProxy) SysRestart(name string, query *common.SystemQuery) error {
	return hp.systemSrv.SendSystemCommand(name, common.CmdSystemRestart, query)
}

// Implementation of common.DeviceManager
func (hp *HandlerProxy) Cancel(name string, query *common.DeviceQuery) error {
	return hp.deviceSrv.SendDeviceCommand(name, common.CmdDeviceCancel, query)
}
func (hp *HandlerProxy) Reset(name string, query *common.DeviceQuery) error {
	return hp.deviceSrv.SendDeviceCommand(name, common.CmdDeviceReset, query)
}
func (hp *HandlerProxy) Status(name string, query *common.DeviceQuery) error {
	return hp.deviceSrv.SendDeviceCommand(name, common.CmdDeviceStatus, query)
}
func (hp *HandlerProxy) RunAction(name string, query *common.DeviceQuery) error {
	return hp.deviceSrv.SendDeviceCommand(name, common.CmdRunAction, query)
}
func (hp *HandlerProxy) StopAction(name string, query *common.DeviceQuery) error {
	return hp.deviceSrv.SendDeviceCommand(name, common.CmdStopAction, query)
}

// Implementation of common.PrinterManager
func (hp *HandlerProxy) InitPrinter(name string, query *common.PrinterSetup) error {
	return hp.printerSrv.SendPrinterCommand(name, common.CmdInitPrinter, query)
}
func (hp *HandlerProxy) PrintText(name string, query *common.PrinterQuery) error {
	return hp.printerSrv.SendPrinterCommand(name, common.CmdPrintText, query)
}

// Implementation of common.ReaderManager
func (hp *HandlerProxy) EnterCard(name string, query *common.DeviceQuery) error {
	return hp.readerSrv.SendReaderCommand(name, common.CmdEnterCard, query)
}
func (hp *HandlerProxy) EjectCard(name string, query *common.DeviceQuery) error {
	return hp.readerSrv.SendReaderCommand(name, common.CmdEjectCard, query)
}
func (hp *HandlerProxy) CaptureCard(name string, query *common.DeviceQuery) error {
	return hp.readerSrv.SendReaderCommand(name, common.CmdCaptureCard, query)
}
func (hp *HandlerProxy) ReadCard(name string, query *common.DeviceQuery) error {
	return hp.readerSrv.SendReaderCommand(name, common.CmdReadCard, query)
}
func (hp *HandlerProxy) ChipGetATR(name string, query *common.DeviceQuery) error {
	return hp.readerSrv.SendReaderCommand(name, common.CmdChipGetATR, query)
}
func (hp *HandlerProxy) ChipPowerOff(name string, query *common.DeviceQuery) error {
	return hp.readerSrv.SendReaderCommand(name, common.CmdChipPowerOff, query)
}
func (hp *HandlerProxy) ChipCommand(name string, query *common.ReaderChipQuery) error {
	return hp.readerSrv.SendReaderCommand(name, common.CmdChipCommand, query)
}

// Implementation of common.ValidatorManager
func (hp *HandlerProxy) InitValidator(name string, query *common.ValidatorQuery) error {
	return hp.validatorSrv.SendValidatorCommand(name, common.CmdInitValidator, query)
}
func (hp *HandlerProxy) DoValidate(name string, query *common.ValidatorQuery) error {
	return hp.validatorSrv.SendValidatorCommand(name, common.CmdDoValidate, query)
}
func (hp *HandlerProxy) NoteAccept(name string, query *common.ValidatorQuery) error {
	return hp.validatorSrv.SendValidatorCommand(name, common.CmdNoteAccept, query)
}
func (hp *HandlerProxy) NoteReturn(name string, query *common.ValidatorQuery) error {
	return hp.validatorSrv.SendValidatorCommand(name, common.CmdNoteReturn, query)
}
func (hp *HandlerProxy) StopValidate(name string, query *common.ValidatorQuery) error {
	return hp.validatorSrv.SendValidatorCommand(name, common.CmdStopValidate, query)
}
func (hp *HandlerProxy) CheckValidator(name string, query *common.ValidatorQuery) error {
	return hp.validatorSrv.SendValidatorCommand(name, common.CmdCheckValidator, query)
}
func (hp *HandlerProxy) ClearValidator(name string, query *common.ValidatorQuery) error {
	return hp.validatorSrv.SendValidatorCommand(name, common.CmdClearValidator, query)
}

// Implementation of common.PinPadManager
func (hp *HandlerProxy) ReadPIN(name string, query *common.ReaderPinQuery) error {
	return hp.pinpadSrv.SendPinPadCommand(name, common.CmdReadPIN, query)
}
func (hp *HandlerProxy) LoadMasterKey(name string, query *common.ReaderPinQuery) error {
	return hp.pinpadSrv.SendPinPadCommand(name, common.CmdLoadMasterKey, query)
}
func (hp *HandlerProxy) LoadWorkKey(name string, query *common.ReaderPinQuery) error {
	return hp.pinpadSrv.SendPinPadCommand(name, common.CmdLoadWorkKey, query)
}
func (hp *HandlerProxy) TestMasterKey(name string, query *common.ReaderPinQuery) error {
	return hp.pinpadSrv.SendPinPadCommand(name, common.CmdTestMasterKey, query)
}
func (hp *HandlerProxy) TestWorkKey(name string, query *common.ReaderPinQuery) error {
	return hp.pinpadSrv.SendPinPadCommand(name, common.CmdTestWorkKey, query)
}
