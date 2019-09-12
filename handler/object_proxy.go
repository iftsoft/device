package handler

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
	"github.com/iftsoft/device/proxy"
)

type ObjectProxy struct {
	server    duplex.ServerManager
	system    *proxy.SystemServer
	device    *proxy.DeviceServer
	reader    *proxy.ReaderServer
	validator *proxy.ValidatorServer
	pinpad    *proxy.PinPadServer
	router    *ObjectRouter
	log       *core.LogAgent
}

func NewObjectProxy() *ObjectProxy {
	op := ObjectProxy{
		server:    nil,
		system:    proxy.NewSystemServer(),
		device:    proxy.NewDeviceServer(),
		reader:    proxy.NewReaderServer(),
		validator: proxy.NewValidatorServer(),
		pinpad:    proxy.NewPinPadServer(),
		router:    NewObjectRouter(),
		log:       core.GetLogAgent(core.LogLevelTrace, "Object"),
	}
	return &op
}

func (op *ObjectProxy) Init(server duplex.ServerManager) {
	op.server = server
	op.system.Init(server, op.router, op.log)
	op.device.Init(server, op.router, op.log)
	op.reader.Init(server, op.router, op.log)
	op.validator.Init(server, op.router, op.log)
	op.pinpad.Init(server, op.router, op.log)
	op.router.InitRouter(op.log, op)
}

func (op *ObjectProxy) Cleanup() {
	op.router.Cleanup()
}

func (op *ObjectProxy) GetClientManager() duplex.ClientManager {
	return op.router
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
func (op *ObjectProxy) RunAction(name string, query *common.DeviceQuery) error {
	return op.device.SendDeviceCommand(name, common.CmdRunAction, query)
}
func (op *ObjectProxy) StopAction(name string, query *common.DeviceQuery) error {
	return op.device.SendDeviceCommand(name, common.CmdStopAction, query)
}

// Implementation of common.ReaderManager
func (op *ObjectProxy) EnterCard(name string, query *common.DeviceQuery) error {
	return op.reader.SendReaderCommand(name, common.CmdEnterCard, query)
}
func (op *ObjectProxy) EjectCard(name string, query *common.DeviceQuery) error {
	return op.reader.SendReaderCommand(name, common.CmdEjectCard, query)
}
func (op *ObjectProxy) CaptureCard(name string, query *common.DeviceQuery) error {
	return op.reader.SendReaderCommand(name, common.CmdCaptureCard, query)
}
func (op *ObjectProxy) ReadCard(name string, query *common.DeviceQuery) error {
	return op.reader.SendReaderCommand(name, common.CmdReadCard, query)
}
func (op *ObjectProxy) ChipGetATR(name string, query *common.DeviceQuery) error {
	return op.reader.SendReaderCommand(name, common.CmdChipGetATR, query)
}
func (op *ObjectProxy) ChipPowerOff(name string, query *common.DeviceQuery) error {
	return op.reader.SendReaderCommand(name, common.CmdChipPowerOff, query)
}
func (op *ObjectProxy) ChipCommand(name string, query *common.ReaderChipQuery) error {
	return op.reader.SendReaderCommand(name, common.CmdChipCommand, query)
}

// Implementation of common.ValidatorManager
func (op *ObjectProxy) InitValidator(name string, query *common.ValidatorQuery) error {
	return op.validator.SendValidatorCommand(name, common.CmdInitValidator, query)
}
func (op *ObjectProxy) DoValidate(name string, query *common.ValidatorQuery) error {
	return op.validator.SendValidatorCommand(name, common.CmdDoValidate, query)
}
func (op *ObjectProxy) NoteAccept(name string, query *common.ValidatorQuery) error {
	return op.validator.SendValidatorCommand(name, common.CmdNoteAccept, query)
}
func (op *ObjectProxy) NoteReject(name string, query *common.ValidatorQuery) error {
	return op.validator.SendValidatorCommand(name, common.CmdNoteReject, query)
}
func (op *ObjectProxy) StopValidate(name string, query *common.ValidatorQuery) error {
	return op.validator.SendValidatorCommand(name, common.CmdStopValidate, query)
}
func (op *ObjectProxy) CheckValidator(name string, query *common.ValidatorQuery) error {
	return op.validator.SendValidatorCommand(name, common.CmdCheckValidator, query)
}
func (op *ObjectProxy) ClearValidator(name string, query *common.ValidatorQuery) error {
	return op.validator.SendValidatorCommand(name, common.CmdClearValidator, query)
}

// Implementation of common.PinPadManager
func (op *ObjectProxy) ReadPIN(name string, query *common.ReaderPinQuery) error {
	return op.pinpad.SendPinPadCommand(name, common.CmdReadPIN, query)
}
func (op *ObjectProxy) LoadMasterKey(name string, query *common.ReaderPinQuery) error {
	return op.pinpad.SendPinPadCommand(name, common.CmdLoadMasterKey, query)
}
func (op *ObjectProxy) LoadWorkKey(name string, query *common.ReaderPinQuery) error {
	return op.pinpad.SendPinPadCommand(name, common.CmdLoadWorkKey, query)
}
func (op *ObjectProxy) TestMasterKey(name string, query *common.ReaderPinQuery) error {
	return op.pinpad.SendPinPadCommand(name, common.CmdTestMasterKey, query)
}
func (op *ObjectProxy) TestWorkKey(name string, query *common.ReaderPinQuery) error {
	return op.pinpad.SendPinPadCommand(name, common.CmdTestWorkKey, query)
}
