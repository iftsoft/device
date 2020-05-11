package validator

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/dbase"
	"github.com/iftsoft/device/dbase/dbvalid"
	"github.com/iftsoft/device/driver"
	"time"
)

type ValidatorDriver struct {
	ValidatorEngine
	storage dbase.DBaseLinker
	begTime int64
}

func NewValidatorDriver() *ValidatorDriver {
	ccd := ValidatorDriver{}
	return &ccd
}

// Implementation of DeviceDriver interface
func (vd *ValidatorDriver) InitDevice(context *driver.Context) error {
	vd.initEngine(context.Config)
	vd.devName = context.DevName
	vd.begTime = time.Now().Unix()
	vd.log.Debug("ValidatorDriver run cmd:%s", "InitDevice")

	mask := common.ScopeFlagSystem
	if device, ok := context.Manager.(common.DeviceCallback); ok {
		vd.device = device
		mask |= common.ScopeFlagDevice
	}
	if validator, ok := context.Manager.(common.ValidatorCallback); ok {
		vd.validator = validator
		mask |= common.ScopeFlagValidator
	}
	if context.Storage != nil {
		vd.storage = context.Storage
		vd.booker  = dbvalid.NewDBaseValidator(vd.storage, vd.devName)
	}
	if context.Greeting != nil {
		context.Greeting.DevType = common.DevTypeCashValidator
		context.Greeting.Required = mask
	}
	return nil
}

func (vd *ValidatorDriver) StartDevice() error {
	vd.log.Debug("ValidatorDriver run cmd:%s", "StartDevice")
	var err error
	if vd.storage != nil {
		err = vd.storage.Open()
	}
	if err == nil {
		err = vd.DevStartup()
	}
	return err
}
func (vd *ValidatorDriver) DeviceTimer(unix int64) error {
	vd.log.Trace("ValidatorDriver run cmd:%s", "DeviceTimer")
	vd.NextMimicStage()
	return nil
}
func (vd *ValidatorDriver) StopDevice() error {
	vd.log.Debug("ValidatorDriver run cmd:%s", "StopDevice")
	err := vd.DevCleanup()
	if vd.storage != nil {
		_ = vd.storage.Close()
	}
	return err
}
func (vd *ValidatorDriver) CheckDevice(metrics *common.SystemMetrics) error {
	vd.log.Debug("ValidatorDriver run cmd:%s", "CheckDevice")
	if metrics != nil {
		metrics.Uptime = time.Now().Unix() - vd.begTime
		metrics.DevState = vd.DevState
		metrics.DevError = vd.DevError
	}
	return nil
}

// Implementation of common.DeviceManager
//
func (vd *ValidatorDriver) Cancel(name string, query *common.DeviceQuery) error {
	err := vd.DevStatus()
	vd.DevError, vd.DevReply = common.CheckError(err)
	return vd.dummyDeviceReply(name, common.CmdDeviceCancel, query)
}
func (vd *ValidatorDriver) Reset(name string, query *common.DeviceQuery) error {
	err := vd.DevReset()
	if err == nil {
		err = vd.DevIdent()
	}
	if err == nil {
		err = vd.DevInitBillList()
	}
	vd.DevError, vd.DevReply = common.CheckError(err)
	return vd.dummyDeviceReply(name, common.CmdDeviceReset, query)
}
func (vd *ValidatorDriver) Status(name string, query *common.DeviceQuery) error {
	err := vd.DevStatus()
	vd.DevError, vd.DevReply = common.CheckError(err)
	return vd.dummyDeviceReply(name, common.CmdDeviceStatus, query)
}
func (vd *ValidatorDriver) RunAction(name string, query *common.DeviceQuery) error {
	err := vd.DevEnableBills(common.CurrencyUAH)
	if err == nil {
		err = vd.DevStatus()
	}
	vd.DevError, vd.DevReply = common.CheckError(err)
	return vd.dummyDeviceReply(name, common.CmdRunAction, query)
}
func (vd *ValidatorDriver) StopAction(name string, query *common.DeviceQuery) error {
	err := vd.DevDisableBills()
	if err == nil {
		err = vd.DevStatus()
	}
	vd.DevError, vd.DevReply = common.CheckError(err)
	return vd.dummyDeviceReply(name, common.CmdStopAction, query)
}

func (vd *ValidatorDriver) dummyDeviceReply(name string, cmd string, query interface{}) error {
	if vd.log != nil {
		vd.log.Debug("ValidatorDriver dev:%s run cmd:%s; Reply: (%d) %s",
			name, cmd, vd.DevError, vd.DevReply)
	}
	var err error
	reply := &common.DeviceReply{}
	reply.Command  = cmd
	reply.Action   = vd.DevAction
	reply.DevState = vd.DevState
	reply.ErrCode  = vd.DevError
	reply.ErrText  = vd.DevReply
	if vd.device != nil {
		err = vd.device.DeviceReply(name, reply)
	}
	return err
}

// Implementation of common.ValidatorManager
//
func (vd *ValidatorDriver) InitValidator(name string, query *common.ValidatorQuery) error {
	err := vd.DevReset()
	if err == nil {
		err = vd.DevInitBillList()
	}
	vd.DevError, vd.DevReply = common.CheckError(err)
	return vd.dummyValidatorStore(name, common.CmdInitValidator, query)
}
func (vd *ValidatorDriver) DoValidate(name string, query *common.ValidatorQuery) error {
	err := vd.DevEnableBills(query.Currency)
	if err == nil {
		err = vd.DevStatus()
	}
	vd.DevError, vd.DevReply = common.CheckError(err)
	return vd.dummyValidatorAccept(name, common.CmdDoValidate, query)
}
func (vd *ValidatorDriver) NoteAccept(name string, query *common.ValidatorQuery) error {
	err := vd.DevNoteAccept()
	vd.DevError, vd.DevReply = common.CheckError(err)
	return vd.dummyValidatorAccept(name, common.CmdNoteAccept, query)
}
func (vd *ValidatorDriver) NoteReturn(name string, query *common.ValidatorQuery) error {
	err := vd.DevNoteReturn()
	vd.DevError, vd.DevReply = common.CheckError(err)
	return vd.dummyValidatorAccept(name, common.CmdNoteReturn, query)
}
func (vd *ValidatorDriver) StopValidate(name string, query *common.ValidatorQuery) error {
	err := vd.DevDisableBills()
	if err == nil {
		err = vd.DevStatus()
	}
	vd.DevError, vd.DevReply = common.CheckError(err)
	return vd.dummyValidatorStore(name, common.CmdStopValidate, query)
}
func (vd *ValidatorDriver) CheckValidator(name string, query *common.ValidatorQuery) error {
	err := vd.DevCheckBatch()
	vd.DevError, vd.DevReply = common.CheckError(err)
	return vd.dummyValidatorStore(name, common.CmdCheckValidator, query)
}
func (vd *ValidatorDriver) ClearValidator(name string, query *common.ValidatorQuery) error {
	err := vd.DevClearBatch()
	vd.DevError, vd.DevReply = common.CheckError(err)
	err = vd.dummyValidatorStore(name, common.CmdClearValidator, query)
	err = vd.DevCheckBatch()
	return err
}

func (vd *ValidatorDriver) dummyValidatorStore(name string, cmd string, query *common.ValidatorQuery) error {
	if vd.log != nil {
		vd.log.Debug("ValidatorDriver dev:%s run cmd:%s; Reply: (%d) %s",
			name, cmd, vd.DevError, vd.DevReply)
	}
	var err error
	reply := &common.ValidatorStore{}
	reply.Command = cmd
	reply.DevState = vd.DevState
	reply.ErrCode = vd.DevError
	reply.ErrText = vd.DevReply
	reply.Notes   = vd.Batch.Notes
	reply.BatchId = vd.Batch.BatchId
	reply.State   = vd.Batch.State
	reply.Detail  = vd.Batch.Detail
	if vd.validator != nil {
		err = vd.validator.ValidatorStore(name, reply)
	}
	return err
}

func (vd *ValidatorDriver) dummyValidatorAccept(name string, cmd string, query *common.ValidatorQuery) error {
	if vd.log != nil {
		vd.log.Debug("ValidatorDriver dev:%s run cmd:%s; Reply: (%d) %s",
			name, cmd, vd.DevError, vd.DevReply)
	}
	reply := &common.ValidatorAccept{}
	reply.Currency = query.Currency
	var err error
	if vd.validator != nil {
		switch cmd {
		case common.CmdDoValidate:
			err = vd.validator.NoteAccepted(name, reply)
		case common.CmdNoteAccept:
			err = vd.validator.CashIsStored(name, reply)
		case common.CmdNoteReturn:
			err = vd.validator.CashReturned(name, reply)
		}
	}
	return err
}
