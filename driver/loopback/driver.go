package loopback

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/dbase"
	"github.com/iftsoft/device/driver"
	"time"
)

type LoopbackDriver struct {
	config    *config.DeviceConfig
	storage   dbase.DBaseLinker
	device    common.DeviceCallback
	printer   common.PrinterCallback
	reader    common.ReaderCallback
	validator common.ValidatorCallback
	pinpad    common.PinPadCallback
	protocol  *LoopbackProtocol
	log       *core.LogAgent
	devState  common.EnumDevState
	devError  common.EnumDevError
	errorText string
	begTime   int64
}

func NewDummyDriver() *LoopbackDriver {
	dd := LoopbackDriver{
		config:    nil,
		storage:   nil,
		device:    nil,
		printer:   nil,
		reader:    nil,
		validator: nil,
		pinpad:    nil,
		protocol:  nil,
		log:       core.GetLogAgent(core.LogLevelTrace, "Driver"),
		devState:  common.DevStateUndefined,
		devError:  common.DevErrorSuccess,
		errorText: "",
		begTime:   time.Now().Unix(),
	}
	return &dd
}

// Implementation of DeviceDriver interface
func (dd *LoopbackDriver) InitDevice(context *driver.Context) error {
	dd.log.Debug("LoopbackDriver run cmd:%s", "InitDevice")
	dd.config = context.Config
	dd.storage = context.Storage
	dd.protocol = GetLoopbackProtocol(dd.config.Linker)

	mask := common.ScopeFlagSystem
	if device, ok := context.Manager.(common.DeviceCallback); ok {
		dd.device = device
		mask |= common.ScopeFlagDevice
	}
	if printer, ok := context.Manager.(common.PrinterCallback); ok {
		dd.printer = printer
		mask |= common.ScopeFlagPrinter
	}
	if reader, ok := context.Manager.(common.ReaderCallback); ok {
		dd.reader = reader
		mask |= common.ScopeFlagReader
	}
	if validator, ok := context.Manager.(common.ValidatorCallback); ok {
		dd.validator = validator
		mask |= common.ScopeFlagValidator
	}
	if pinpad, ok := context.Manager.(common.PinPadCallback); ok {
		dd.pinpad = pinpad
		mask |= common.ScopeFlagPinPad
	}
	if context.Greeting != nil {
		context.Greeting.DevType = common.DevTypeCommon
		context.Greeting.Required = mask
	}
	return nil
}

func (dd *LoopbackDriver) StartDevice(query *common.SystemConfig) error {
	dd.log.Debug("LoopbackDriver run cmd:%s", "StartDeviceLoop")
	if dd.config != nil && query != nil {
		dd.config.OverwriteConfig(query)
	}
	err := dd.protocol.OpenLink()
	if err == nil && dd.storage != nil {
		err = dd.storage.Open()
	}
	return err
}
func (dd *LoopbackDriver) DeviceTimer(unix int64) error {
	dd.log.Debug("LoopbackDriver run cmd:%s", "DeviceTimer")
	return nil
}
func (dd *LoopbackDriver) StopDevice() error {
	dd.log.Debug("LoopbackDriver run cmd:%s", "StopDeviceLoop")
	err := dd.protocol.CloseLink()
	if dd.storage != nil {
		err = dd.storage.Close()
	}
	return err
}
func (dd *LoopbackDriver) CheckDevice(metrics *common.SystemMetrics) error {
	dd.log.Debug("LoopbackDriver run cmd:%s", "CheckDevice")
	if metrics != nil {
		metrics.Uptime = time.Now().Unix() - dd.begTime
		metrics.DevState = dd.devState
		metrics.DevError = dd.devError
	}
	return nil
}

// Implementation of common.DeviceManager
//
func (dd *LoopbackDriver) Cancel(name string, query *common.DeviceQuery) error {
	err := dd.protocol.CheckLink()
	dd.devError, dd.errorText = common.CheckError(err)
	return dd.dummyDeviceReply(name, common.CmdDeviceCancel, query)
}
func (dd *LoopbackDriver) Reset(name string, query *common.DeviceQuery) error {
	err := dd.protocol.CheckLink()
	dd.devError, dd.errorText = common.CheckError(err)
	return dd.dummyDeviceReply(name, common.CmdDeviceReset, query)
}
func (dd *LoopbackDriver) Status(name string, query *common.DeviceQuery) error {
	err := dd.protocol.CheckLink()
	dd.devError, dd.errorText = common.CheckError(err)
	return dd.dummyDeviceReply(name, common.CmdDeviceStatus, query)
}
func (dd *LoopbackDriver) RunAction(name string, query *common.DeviceQuery) error {
	err := dd.protocol.CheckLink()
	dd.devError, dd.errorText = common.CheckError(err)
	return dd.dummyDeviceReply(name, common.CmdRunAction, query)
}
func (dd *LoopbackDriver) StopAction(name string, query *common.DeviceQuery) error {
	err := dd.protocol.CheckLink()
	dd.devError, dd.errorText = common.CheckError(err)
	return dd.dummyDeviceReply(name, common.CmdStopAction, query)
}

func (dd *LoopbackDriver) dummyDeviceReply(name string, cmd string, query interface{}) error {
	if dd.log != nil {
		dd.log.Debug("LoopbackDriver dev:%s run cmd:%s with result: (%d) %s",
			name, cmd, dd.devError, dd.errorText)
	}
	reply := &common.DeviceReply{}
	reply.Command = cmd
	reply.DevState = dd.devState
	reply.ErrCode = dd.devError
	reply.ErrText = dd.errorText
	if dd.device != nil {
		return dd.device.DeviceReply(name, reply)
	}
	return nil
}

// Implementation of common.PrinterManager
//
func (dd *LoopbackDriver) InitPrinter(name string, query *common.PrinterSetup) error {
	err := dd.protocol.CheckLink()
	dd.devError, dd.errorText = common.CheckError(err)
	return dd.dummyDeviceReply(name, common.CmdInitPrinter, query)
}
func (dd *LoopbackDriver) PrintText(name string, query *common.PrinterQuery) error {
	err := dd.protocol.CheckLink()
	dd.devError, dd.errorText = common.CheckError(err)
	return dd.dummyDeviceReply(name, common.CmdPrintText, query)
}

// Implementation of common.ReaderManager
//
func (dd *LoopbackDriver) EnterCard(name string, query *common.DeviceQuery) error {
	err := dd.protocol.CheckLink()
	dd.devError, dd.errorText = common.CheckError(err)
	return dd.dummyDeviceReply(name, common.CmdEnterCard, query)
}
func (dd *LoopbackDriver) EjectCard(name string, query *common.DeviceQuery) error {
	err := dd.protocol.CheckLink()
	dd.devError, dd.errorText = common.CheckError(err)
	return dd.dummyDeviceReply(name, common.CmdEjectCard, query)
}
func (dd *LoopbackDriver) CaptureCard(name string, query *common.DeviceQuery) error {
	err := dd.protocol.CheckLink()
	dd.devError, dd.errorText = common.CheckError(err)
	return dd.dummyDeviceReply(name, common.CmdCaptureCard, query)
}
func (dd *LoopbackDriver) ReadCard(name string, query *common.DeviceQuery) error {
	err := dd.protocol.CheckLink()
	dd.devError, dd.errorText = common.CheckError(err)
	return dd.dummyDeviceReply(name, common.CmdReadCard, query)
}
func (dd *LoopbackDriver) ChipGetATR(name string, query *common.DeviceQuery) error {
	err := dd.protocol.CheckLink()
	dd.devError, dd.errorText = common.CheckError(err)
	return dd.dummyDeviceReply(name, common.CmdChipGetATR, query)
}
func (dd *LoopbackDriver) ChipPowerOff(name string, query *common.DeviceQuery) error {
	err := dd.protocol.CheckLink()
	dd.devError, dd.errorText = common.CheckError(err)
	return dd.dummyDeviceReply(name, common.CmdChipPowerOff, query)
}
func (dd *LoopbackDriver) ChipCommand(name string, query *common.ReaderChipQuery) error {
	err := dd.protocol.CheckLink()
	dd.devError, dd.errorText = common.CheckError(err)
	return dd.dummyDeviceReply(name, common.CmdChipCommand, query)
}

// Implementation of common.ValidatorManager
//
func (dd *LoopbackDriver) InitValidator(name string, query *common.ValidatorQuery) error {
	err := dd.protocol.CheckLink()
	dd.devError, dd.errorText = common.CheckError(err)
	return dd.dummyValidatorStore(name, common.CmdInitValidator, query)
}
func (dd *LoopbackDriver) DoValidate(name string, query *common.ValidatorQuery) error {
	err := dd.protocol.CheckLink()
	dd.devError, dd.errorText = common.CheckError(err)
	return dd.dummyValidatorAccept(name, common.CmdDoValidate, query)
}
func (dd *LoopbackDriver) NoteAccept(name string, query *common.ValidatorQuery) error {
	err := dd.protocol.CheckLink()
	dd.devError, dd.errorText = common.CheckError(err)
	return dd.dummyValidatorAccept(name, common.CmdNoteAccept, query)
}
func (dd *LoopbackDriver) NoteReturn(name string, query *common.ValidatorQuery) error {
	err := dd.protocol.CheckLink()
	dd.devError, dd.errorText = common.CheckError(err)
	return dd.dummyValidatorAccept(name, common.CmdNoteReturn, query)
}
func (dd *LoopbackDriver) StopValidate(name string, query *common.ValidatorQuery) error {
	err := dd.protocol.CheckLink()
	dd.devError, dd.errorText = common.CheckError(err)
	return dd.dummyValidatorStore(name, common.CmdStopValidate, query)
}
func (dd *LoopbackDriver) CheckValidator(name string, query *common.ValidatorQuery) error {
	err := dd.protocol.CheckLink()
	dd.devError, dd.errorText = common.CheckError(err)
	return dd.dummyValidatorStore(name, common.CmdCheckValidator, query)
}
func (dd *LoopbackDriver) ClearValidator(name string, query *common.ValidatorQuery) error {
	err := dd.protocol.CheckLink()
	dd.devError, dd.errorText = common.CheckError(err)
	return dd.dummyValidatorStore(name, common.CmdClearValidator, query)
}

func (dd *LoopbackDriver) dummyValidatorStore(name string, cmd string, query *common.ValidatorQuery) error {
	if dd.log != nil {
		dd.log.Debug("LoopbackDriver dev:%s run cmd:%s with result: (%d) %s",
			name, cmd, dd.devError, dd.errorText)
	}
	reply := &common.ValidatorStore{}
	reply.Command = cmd
	reply.DevState = dd.devState
	reply.ErrCode = dd.devError
	reply.ErrText = dd.errorText
	var err error
	if dd.validator != nil {
		err = dd.validator.ValidatorStore(name, reply)
	}
	return err
}

func (dd *LoopbackDriver) dummyValidatorAccept(name string, cmd string, query *common.ValidatorQuery) error {
	if dd.log != nil {
		dd.log.Debug("LoopbackDriver dev:%s run cmd:%s", name, cmd)
	}
	reply := &common.ValidatorAccept{}
	reply.Currency = query.Currency
	var err error
	if dd.validator != nil {
		switch cmd {
		case common.CmdDoValidate:
			err = dd.validator.NoteAccepted(name, reply)
		case common.CmdNoteAccept:
			err = dd.validator.CashIsStored(name, reply)
		case common.CmdNoteReturn:
			err = dd.validator.CashReturned(name, reply)
		}
	}
	return err
}

// Implementation of common.PinPadManager
//
func (dd *LoopbackDriver) ReadPIN(name string, query *common.ReaderPinQuery) error {
	err := dd.protocol.CheckLink()
	dd.devError, dd.errorText = common.CheckError(err)
	return dd.dummyDeviceReply(name, common.CmdReadPIN, query)
}
func (dd *LoopbackDriver) LoadMasterKey(name string, query *common.ReaderPinQuery) error {
	err := dd.protocol.CheckLink()
	dd.devError, dd.errorText = common.CheckError(err)
	return dd.dummyPinPadReply(name, common.CmdLoadMasterKey, query)
}
func (dd *LoopbackDriver) LoadWorkKey(name string, query *common.ReaderPinQuery) error {
	err := dd.protocol.CheckLink()
	dd.devError, dd.errorText = common.CheckError(err)
	return dd.dummyPinPadReply(name, common.CmdLoadWorkKey, query)
}
func (dd *LoopbackDriver) TestMasterKey(name string, query *common.ReaderPinQuery) error {
	err := dd.protocol.CheckLink()
	dd.devError, dd.errorText = common.CheckError(err)
	return dd.dummyPinPadReply(name, common.CmdTestMasterKey, query)
}
func (dd *LoopbackDriver) TestWorkKey(name string, query *common.ReaderPinQuery) error {
	err := dd.protocol.CheckLink()
	dd.devError, dd.errorText = common.CheckError(err)
	return dd.dummyPinPadReply(name, common.CmdTestWorkKey, query)
}

func (dd *LoopbackDriver) dummyPinPadReply(name string, cmd string, query interface{}) error {
	if dd.log != nil {
		dd.log.Debug("LoopbackDriver dev:%s run cmd:%s with result: (%d) %s",
			name, cmd, dd.devError, dd.errorText)
	}
	reply := &common.ReaderPinReply{}
	reply.Command = cmd
	reply.DevState = dd.devState
	reply.ErrCode = dd.devError
	reply.ErrText = dd.errorText
	var err error
	if dd.reader != nil {
		err = dd.pinpad.PinPadReply(name, reply)
	}
	return err
}
