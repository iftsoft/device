package handler

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
	"github.com/iftsoft/device/proxy"
)

type HandlerProxy struct {
	HandlerRouter
	serverMng    duplex.ServerManager
	systemSrv    *proxy.SystemServer
	deviceSrv    *proxy.DeviceServer
	printerSrv   *proxy.PrinterServer
	readerSrv    *proxy.ReaderServer
	validatorSrv *proxy.ValidatorServer
	pinpadSrv    *proxy.PinPadServer
}

func NewHandlerProxy(config config.HandlerList) *HandlerProxy {
	op := &HandlerProxy{
		HandlerRouter: HandlerRouter{
			config:    config,
			proxy:  nil,
			objMap: make(map[string]*DeviceHandler),
			log:    core.GetLogAgent(core.LogLevelTrace, "Router"),
		},
		serverMng:    nil,
		systemSrv:    proxy.NewSystemServer(),
		deviceSrv:    proxy.NewDeviceServer(),
		printerSrv:   proxy.NewPrinterServer(),
		readerSrv:    proxy.NewReaderServer(),
		validatorSrv: proxy.NewValidatorServer(),
		pinpadSrv:    proxy.NewPinPadServer(),
	}
	op.HandlerRouter.proxy = op
	return op
}

func (op *HandlerProxy) Init(server duplex.ServerManager) {
	op.serverMng = server
	op.systemSrv.Init(server, &op.HandlerRouter, op.log)
	op.deviceSrv.Init(server, &op.HandlerRouter, op.log)
	op.printerSrv.Init(server, &op.HandlerRouter, op.log)
	op.readerSrv.Init(server, &op.HandlerRouter, op.log)
	op.validatorSrv.Init(server, &op.HandlerRouter, op.log)
	op.pinpadSrv.Init(server, &op.HandlerRouter, op.log)
}



// Implementation of common.SystemManager
func (op *HandlerProxy) Terminate(name string, query *common.SystemQuery) error {
	return op.systemSrv.SendSystemCommand(name, common.CmdSystemTerminate, query)
}
func (op *HandlerProxy) SysInform(name string, query *common.SystemQuery) error {
	return op.systemSrv.SendSystemCommand(name, common.CmdSystemInform, query)
}
func (op *HandlerProxy) SysStart(name string, query *common.SystemQuery) error {
	return op.systemSrv.SendSystemCommand(name, common.CmdSystemStart, query)
}
func (op *HandlerProxy) SysStop(name string, query *common.SystemQuery) error {
	return op.systemSrv.SendSystemCommand(name, common.CmdSystemStop, query)
}
func (op *HandlerProxy) SysRestart(name string, query *common.SystemQuery) error {
	return op.systemSrv.SendSystemCommand(name, common.CmdSystemRestart, query)
}

// Implementation of common.DeviceManager
func (op *HandlerProxy) Cancel(name string, query *common.DeviceQuery) error {
	return op.deviceSrv.SendDeviceCommand(name, common.CmdDeviceCancel, query)
}
func (op *HandlerProxy) Reset(name string, query *common.DeviceQuery) error {
	return op.deviceSrv.SendDeviceCommand(name, common.CmdDeviceReset, query)
}
func (op *HandlerProxy) Status(name string, query *common.DeviceQuery) error {
	return op.deviceSrv.SendDeviceCommand(name, common.CmdDeviceStatus, query)
}
func (op *HandlerProxy) RunAction(name string, query *common.DeviceQuery) error {
	return op.deviceSrv.SendDeviceCommand(name, common.CmdRunAction, query)
}
func (op *HandlerProxy) StopAction(name string, query *common.DeviceQuery) error {
	return op.deviceSrv.SendDeviceCommand(name, common.CmdStopAction, query)
}

// Implementation of common.PrinterManager
func (op *HandlerProxy) InitPrinter(name string, query *common.PrinterSetup) error {
	return op.printerSrv.SendPrinterCommand(name, common.CmdInitPrinter, query)
}
func (op *HandlerProxy) PrintText(name string, query *common.PrinterQuery) error {
	return op.printerSrv.SendPrinterCommand(name, common.CmdPrintText, query)
}

// Implementation of common.ReaderManager
func (op *HandlerProxy) EnterCard(name string, query *common.DeviceQuery) error {
	return op.readerSrv.SendReaderCommand(name, common.CmdEnterCard, query)
}
func (op *HandlerProxy) EjectCard(name string, query *common.DeviceQuery) error {
	return op.readerSrv.SendReaderCommand(name, common.CmdEjectCard, query)
}
func (op *HandlerProxy) CaptureCard(name string, query *common.DeviceQuery) error {
	return op.readerSrv.SendReaderCommand(name, common.CmdCaptureCard, query)
}
func (op *HandlerProxy) ReadCard(name string, query *common.DeviceQuery) error {
	return op.readerSrv.SendReaderCommand(name, common.CmdReadCard, query)
}
func (op *HandlerProxy) ChipGetATR(name string, query *common.DeviceQuery) error {
	return op.readerSrv.SendReaderCommand(name, common.CmdChipGetATR, query)
}
func (op *HandlerProxy) ChipPowerOff(name string, query *common.DeviceQuery) error {
	return op.readerSrv.SendReaderCommand(name, common.CmdChipPowerOff, query)
}
func (op *HandlerProxy) ChipCommand(name string, query *common.ReaderChipQuery) error {
	return op.readerSrv.SendReaderCommand(name, common.CmdChipCommand, query)
}

// Implementation of common.ValidatorManager
func (op *HandlerProxy) InitValidator(name string, query *common.ValidatorQuery) error {
	return op.validatorSrv.SendValidatorCommand(name, common.CmdInitValidator, query)
}
func (op *HandlerProxy) DoValidate(name string, query *common.ValidatorQuery) error {
	return op.validatorSrv.SendValidatorCommand(name, common.CmdDoValidate, query)
}
func (op *HandlerProxy) NoteAccept(name string, query *common.ValidatorQuery) error {
	return op.validatorSrv.SendValidatorCommand(name, common.CmdNoteAccept, query)
}
func (op *HandlerProxy) NoteReject(name string, query *common.ValidatorQuery) error {
	return op.validatorSrv.SendValidatorCommand(name, common.CmdNoteReject, query)
}
func (op *HandlerProxy) StopValidate(name string, query *common.ValidatorQuery) error {
	return op.validatorSrv.SendValidatorCommand(name, common.CmdStopValidate, query)
}
func (op *HandlerProxy) CheckValidator(name string, query *common.ValidatorQuery) error {
	return op.validatorSrv.SendValidatorCommand(name, common.CmdCheckValidator, query)
}
func (op *HandlerProxy) ClearValidator(name string, query *common.ValidatorQuery) error {
	return op.validatorSrv.SendValidatorCommand(name, common.CmdClearValidator, query)
}

// Implementation of common.PinPadManager
func (op *HandlerProxy) ReadPIN(name string, query *common.ReaderPinQuery) error {
	return op.pinpadSrv.SendPinPadCommand(name, common.CmdReadPIN, query)
}
func (op *HandlerProxy) LoadMasterKey(name string, query *common.ReaderPinQuery) error {
	return op.pinpadSrv.SendPinPadCommand(name, common.CmdLoadMasterKey, query)
}
func (op *HandlerProxy) LoadWorkKey(name string, query *common.ReaderPinQuery) error {
	return op.pinpadSrv.SendPinPadCommand(name, common.CmdLoadWorkKey, query)
}
func (op *HandlerProxy) TestMasterKey(name string, query *common.ReaderPinQuery) error {
	return op.pinpadSrv.SendPinPadCommand(name, common.CmdTestMasterKey, query)
}
func (op *HandlerProxy) TestWorkKey(name string, query *common.ReaderPinQuery) error {
	return op.pinpadSrv.SendPinPadCommand(name, common.CmdTestWorkKey, query)
}
