package driver

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/linker"
)

type DummyDriver struct {
	config    *config.DeviceConfig
	device    common.DeviceCallback
	printer   common.PrinterCallback
	reader    common.ReaderCallback
	validator common.ValidatorCallback
	pinpad    common.PinPadCallback
	loopback  *linker.Loopback
	log       *core.LogAgent
	devState  common.EnumDevState
}

func NewDummyDriver() *DummyDriver {
	dd := DummyDriver{
		config:    nil,
		device:    nil,
		printer:   nil,
		reader:    nil,
		validator: nil,
		pinpad:    nil,
		loopback:  nil,
		log:       core.GetLogAgent(core.LogLevelTrace, "Dummy"),
		devState:  common.DevStateUndefined,
	}
	return &dd
}

// Implementation of DeviceDriver interface
func (dd *DummyDriver) InitDevice(manager interface{}, cfg *config.DeviceConfig) common.DevScopeMask {
	dd.log.Debug("DummyDriver run cmd:%s", "InitDevice")
	dd.config = cfg
	dd.loopback = linker.GetLoopback(dd.config, dd.log)

	mask := common.ScopeFlagSystem
	if device, ok := manager.(common.DeviceCallback); ok {
		dd.device = device
		mask |= common.ScopeFlagDevice
	}
	if printer, ok := manager.(common.PrinterCallback); ok {
		dd.printer = printer
		mask |= common.ScopeFlagPrinter
	}
	if reader, ok := manager.(common.ReaderCallback); ok {
		dd.reader = reader
		mask |= common.ScopeFlagReader
	}
	if validator, ok := manager.(common.ValidatorCallback); ok {
		dd.validator = validator
		mask |= common.ScopeFlagValidator
	}
	if pinpad, ok := manager.(common.PinPadCallback); ok {
		dd.pinpad = pinpad
		mask |= common.ScopeFlagPinPad
	}
	return mask
}

func (dd *DummyDriver) StartDevice() error {
	dd.log.Debug("DummyDriver run cmd:%s", "StartDevice")
	err := dd.loopback.OpenLink()
	return err
}
func (dd *DummyDriver) DeviceTimer(unix int64) error {
	dd.log.Debug("DummyDriver run cmd:%s", "DeviceTimer")
	return nil
}
func (dd *DummyDriver) StopDevice() error {
	dd.log.Debug("DummyDriver run cmd:%s", "StopDevice")
	err := dd.loopback.CloseLink()
	return err
}
func (dd *DummyDriver) CheckDevice(metrics *common.SystemMetrics) error {
	dd.log.Debug("DummyDriver run cmd:%s", "CheckDevice")
	return nil
}

// Implementation of common.DeviceManager
//
func (dd *DummyDriver) Cancel(name string, query *common.DeviceQuery) error {
	return dd.dummyDeviceReply(name, common.CmdDeviceCancel, query)
}
func (dd *DummyDriver) Reset(name string, query *common.DeviceQuery) error {
	return dd.dummyDeviceReply(name, common.CmdDeviceReset, query)
}
func (dd *DummyDriver) Status(name string, query *common.DeviceQuery) error {
	return dd.dummyDeviceReply(name, common.CmdDeviceStatus, query)
}
func (dd *DummyDriver) RunAction(name string, query *common.DeviceQuery) error {
	return dd.dummyDeviceReply(name, common.CmdRunAction, query)
}
func (dd *DummyDriver) StopAction(name string, query *common.DeviceQuery) error {
	return dd.dummyDeviceReply(name, common.CmdStopAction, query)
}

func (dd *DummyDriver) dummyDeviceReply(name string, cmd string, query interface{}) error {
	if dd.log != nil {
		dd.log.Debug("DummyDriver dev:%s run cmd:%s", name, cmd)
	}
	err := dd.loopback.CheckLink()
	reply := &common.DeviceReply{}
	reply.Command = cmd
	reply.DevState = dd.devState
	reply.ErrText = core.GetErrorText(err)
	if dd.device != nil {
		err = dd.device.DeviceReply(name, reply)
	}
	return err
}

// Implementation of common.PrinterManager
//
func (dd *DummyDriver) InitPrinter(name string, query *common.PrinterSetup) error {
	return dd.dummyDeviceReply(name, common.CmdInitPrinter, query)
}
func (dd *DummyDriver) PrintText(name string, query *common.PrinterQuery) error {
	return dd.dummyDeviceReply(name, common.CmdPrintText, query)
}

// Implementation of common.ReaderManager
//
func (dd *DummyDriver) EnterCard(name string, query *common.DeviceQuery) error {
	return dd.dummyDeviceReply(name, common.CmdEnterCard, query)
}
func (dd *DummyDriver) EjectCard(name string, query *common.DeviceQuery) error {
	return dd.dummyDeviceReply(name, common.CmdEjectCard, query)
}
func (dd *DummyDriver) CaptureCard(name string, query *common.DeviceQuery) error {
	return dd.dummyDeviceReply(name, common.CmdCaptureCard, query)
}
func (dd *DummyDriver) ReadCard(name string, query *common.DeviceQuery) error {
	return dd.dummyDeviceReply(name, common.CmdReadCard, query)
}
func (dd *DummyDriver) ChipGetATR(name string, query *common.DeviceQuery) error {
	return dd.dummyDeviceReply(name, common.CmdChipGetATR, query)
}
func (dd *DummyDriver) ChipPowerOff(name string, query *common.DeviceQuery) error {
	return dd.dummyDeviceReply(name, common.CmdChipPowerOff, query)
}
func (dd *DummyDriver) ChipCommand(name string, query *common.ReaderChipQuery) error {
	return dd.dummyDeviceReply(name, common.CmdChipCommand, query)
}

// Implementation of common.ValidatorManager
//
func (dd *DummyDriver) InitValidator(name string, query *common.ValidatorQuery) error {
	return dd.dummyValidatorStore(name, common.CmdInitValidator, query)
}
func (dd *DummyDriver) DoValidate(name string, query *common.ValidatorQuery) error {
	return dd.dummyValidatorAccept(name, common.CmdDoValidate, query)
}
func (dd *DummyDriver) NoteAccept(name string, query *common.ValidatorQuery) error {
	return dd.dummyValidatorAccept(name, common.CmdNoteAccept, query)
}
func (dd *DummyDriver) NoteReject(name string, query *common.ValidatorQuery) error {
	return dd.dummyValidatorAccept(name, common.CmdNoteReject, query)
}
func (dd *DummyDriver) StopValidate(name string, query *common.ValidatorQuery) error {
	return dd.dummyValidatorStore(name, common.CmdStopValidate, query)
}
func (dd *DummyDriver) CheckValidator(name string, query *common.ValidatorQuery) error {
	return dd.dummyValidatorStore(name, common.CmdCheckValidator, query)
}
func (dd *DummyDriver) ClearValidator(name string, query *common.ValidatorQuery) error {
	return dd.dummyValidatorStore(name, common.CmdClearValidator, query)
}

func (dd *DummyDriver) dummyValidatorStore(name string, cmd string, query *common.ValidatorQuery) error {
	if dd.log != nil {
		dd.log.Debug("DummyDriver dev:%s run cmd:%s", name, cmd)
	}
	reply := &common.ValidatorStore{}
	var err error
	if dd.validator != nil {
		err = dd.validator.ValidatorStore(name, reply)
	}
	return err
}

func (dd *DummyDriver) dummyValidatorAccept(name string, cmd string, query *common.ValidatorQuery) error {
	if dd.log != nil {
		dd.log.Debug("DummyDriver dev:%s run cmd:%s", name, cmd)
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
		case common.CmdNoteReject:
			err = dd.validator.CashReturned(name, reply)
		}
	}
	return err
}

// Implementation of common.PinPadManager
//
func (dd *DummyDriver) ReadPIN(name string, query *common.ReaderPinQuery) error {
	return dd.dummyDeviceReply(name, common.CmdReadPIN, query)
}
func (dd *DummyDriver) LoadMasterKey(name string, query *common.ReaderPinQuery) error {
	return dd.dummyPinPadReply(name, common.CmdLoadMasterKey, query)
}
func (dd *DummyDriver) LoadWorkKey(name string, query *common.ReaderPinQuery) error {
	return dd.dummyPinPadReply(name, common.CmdLoadWorkKey, query)
}
func (dd *DummyDriver) TestMasterKey(name string, query *common.ReaderPinQuery) error {
	return dd.dummyPinPadReply(name, common.CmdTestMasterKey, query)
}
func (dd *DummyDriver) TestWorkKey(name string, query *common.ReaderPinQuery) error {
	return dd.dummyPinPadReply(name, common.CmdTestWorkKey, query)
}

func (dd *DummyDriver) dummyPinPadReply(name string, cmd string, query interface{}) error {
	if dd.log != nil {
		dd.log.Debug("DummyDriver dev:%s run cmd:%s", name, cmd)
	}
	reply := &common.ReaderPinReply{}
	var err error
	if dd.reader != nil {
		err = dd.pinpad.PinPadReply(name, reply)
	}
	return err
}
