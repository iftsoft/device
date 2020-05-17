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
	vd.DevName = context.DevName
	vd.begTime = time.Now().Unix()
	vd.Log.Debug("ValidatorDriver run cmd:%s", "InitDevice")

	mask := common.ScopeFlagSystem
	if device, ok := context.Manager.(common.DeviceCallback); ok {
		vd.CbDevice = device
		mask |= common.ScopeFlagDevice
	}
	if validator, ok := context.Manager.(common.ValidatorCallback); ok {
		vd.CbValidator = validator
		mask |= common.ScopeFlagValidator
	}
	if context.Storage != nil {
		vd.storage = context.Storage
		vd.booker  = dbvalid.NewDBaseValidator(vd.storage, vd.DevName)
	}
	if context.Greeting != nil {
		context.Greeting.DevType = common.DevTypeCashValidator
		context.Greeting.Required = mask
	}
	return nil
}

func (vd *ValidatorDriver) StartDevice() error {
	vd.Log.Debug("ValidatorDriver run cmd:%s", "StartDevice")
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
	vd.Log.Trace("ValidatorDriver run cmd:%s", "DeviceTimer")
	vd.NextMimicStage()
	return nil
}
func (vd *ValidatorDriver) StopDevice() error {
	vd.Log.Debug("ValidatorDriver run cmd:%s", "StopDevice")
	err := vd.DevCleanup()
	if vd.storage != nil {
		_ = vd.storage.Close()
	}
	return err
}
func (vd *ValidatorDriver) CheckDevice(metrics *common.SystemMetrics) error {
	vd.Log.Debug("ValidatorDriver run cmd:%s", "CheckDevice")
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
	return vd.RunDeviceReply(common.CmdDeviceCancel)
}
func (vd *ValidatorDriver) Reset(name string, query *common.DeviceQuery) error {
	err := vd.DevReset()
	if err == nil {
		err = vd.DevIdent()
	}
	vd.DevError, vd.DevReply = common.CheckError(err)
	return vd.RunDeviceReply(common.CmdDeviceReset)
}
func (vd *ValidatorDriver) Status(name string, query *common.DeviceQuery) error {
	err := vd.DevStatus()
	vd.DevError, vd.DevReply = common.CheckError(err)
	return vd.RunDeviceReply(common.CmdDeviceStatus)
}
func (vd *ValidatorDriver) RunAction(name string, query *common.DeviceQuery) error {
	err := vd.DevEnableBills(common.CurrencyUAH)
	if err == nil {
		err = vd.DevStatus()
	}
	vd.DevError, vd.DevReply = common.CheckError(err)
	return vd.RunDeviceReply(common.CmdRunAction)
}
func (vd *ValidatorDriver) StopAction(name string, query *common.DeviceQuery) error {
	err := vd.DevDisableBills()
	if err == nil {
		err = vd.DevStatus()
	}
	vd.DevError, vd.DevReply = common.CheckError(err)
	return vd.RunDeviceReply(common.CmdStopAction)
}


// Implementation of common.ValidatorManager
//
func (vd *ValidatorDriver) InitValidator(name string, query *common.ValidatorQuery) error {
	err := vd.DevReset()
	if err == nil {
		err = vd.DevInitBillList()
	}
	vd.DevError, vd.DevReply = common.CheckError(err)
	return vd.RunValidatorStore(common.CmdInitValidator)
}
func (vd *ValidatorDriver) DoValidate(name string, query *common.ValidatorQuery) error {
	err := vd.DevEnableBills(query.Currency)
	if err == nil {
		err = vd.DevStatus()
	}
	vd.DevError, vd.DevReply = common.CheckError(err)
	return vd.RunValidatorStore(common.CmdDoValidate)
}
func (vd *ValidatorDriver) NoteAccept(name string, query *common.ValidatorQuery) error {
	err := vd.DevNoteAccept()
	vd.DevError, vd.DevReply = common.CheckError(err)
	return err
}
func (vd *ValidatorDriver) NoteReturn(name string, query *common.ValidatorQuery) error {
	err := vd.DevNoteReturn()
	vd.DevError, vd.DevReply = common.CheckError(err)
	return err
}
func (vd *ValidatorDriver) StopValidate(name string, query *common.ValidatorQuery) error {
	err := vd.DevDisableBills()
	if err == nil {
		err = vd.DevStatus()
	}
	vd.DevError, vd.DevReply = common.CheckError(err)
	return vd.RunValidatorStore(common.CmdStopValidate)
}
func (vd *ValidatorDriver) CheckValidator(name string, query *common.ValidatorQuery) error {
	err := vd.DevCheckBatch()
	vd.DevError, vd.DevReply = common.CheckError(err)
	return vd.RunValidatorStore(common.CmdCheckValidator)
}
func (vd *ValidatorDriver) ClearValidator(name string, query *common.ValidatorQuery) error {
	err := vd.DevClearBatch()
	vd.DevError, vd.DevReply = common.CheckError(err)
	err = vd.RunValidatorStore(common.CmdClearValidator)
	err = vd.DevCheckBatch()
	return err
}


