package router

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/dbase"
	"github.com/iftsoft/device/driver"
	"sync"
	"time"
)

type SystemDevice struct {
	devName string
	state   common.EnumSystemState
	error   common.EnumSystemError
	driver  driver.DeviceDriver
	storage dbase.DBaseLinker
	manager common.ManagerSet
	system  common.SystemCallback
	config  *config.DeviceConfig
	log     *core.LogAgent
	done    chan struct{}
	wg      sync.WaitGroup
	checkTm time.Time
}

func NewSystemDevice(name string) *SystemDevice {
	sd := SystemDevice{
		devName: name,
		state:   common.SysStateUndefined,
		error:   common.SysErrSuccess,
		driver:  nil,
		system:  nil,
		log:     core.GetLogAgent(core.LogLevelTrace, name),
		done:    make(chan struct{}),
	}
	return &sd
}

func (sd *SystemDevice) InitDevice(drv driver.DeviceDriver, ctx *driver.Context) {
	sd.driver = drv
	sd.config = ctx.Config
	sd.storage = ctx.Storage
	sd.system = ctx.Complex.GetSystemCallback()
	manager := drv.InitDevice(ctx)
	sd.manager.InitManagers(manager)
}

func (sd *SystemDevice) StartDeviceLoop() {
	sd.log.Info("Starting system device %s", sd.devName)
	go sd.deviceLoop(&sd.wg)
}

func (sd *SystemDevice) StopDeviceLoop() {
	sd.log.Info("Stopping system device %s", sd.devName)
	close(sd.done)
	sd.wg.Wait()
}

func (sd *SystemDevice) deviceLoop(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	sd.log.Debug("System device %s loop is started", sd.devName)
	defer sd.log.Debug("System device %s loop is stopped", sd.devName)

	sd.state = common.SysStateUndefined
	if sd.config.Common.AutoLoad {
		err := sd.driver.StartDevice(nil)
		if err == nil {
			sd.state = common.SysStateRunning
		}
	}
	defer func() {
		err := sd.driver.StopDevice()
		if err == nil {
			sd.state = common.SysStateStopped
		}
	}()

	sd.checkTm = time.Now()
	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()

	for {
		select {
		case <-sd.done:
			return
		case tm := <-tick.C:
			sd.onTimerTick(tm)
		}
	}
}

func (sd *SystemDevice) onTimerTick(tm time.Time) {
	sd.log.Trace("System device %s onTimerTick %s", sd.devName, tm.Format(time.StampMilli))
	if sd.driver != nil && sd.state == common.SysStateRunning {
		_ = sd.driver.DeviceTimer(tm.Unix())
	}
	delta := int(tm.Sub(sd.checkTm).Seconds())
	if delta == 60 {
		sd.sendDeviceMetrics(tm)
	}
}

func (sd *SystemDevice) sendDeviceMetrics(tm time.Time) {
	sd.log.Debug("System device %s check device metric", sd.devName)
	sd.checkTm = tm
	health := common.NewSystemHealth()
	health.Moment = tm.Unix()
	health.State = sd.state
	health.Error = sd.error
	err := sd.driver.CheckDevice(&health.Metrics)
	if err != nil {
		health.Error = common.SysErrDeviceFail
	}
	if sd.system != nil {
		err = sd.system.SystemHealth(sd.devName, health)
	}
}

//// Implementation of common.ComplexManager
//
//func (sd *SystemDevice) GetSystemManager() common.SystemManager {
//	return sd
//}
//func (sd *SystemDevice) GetDeviceManager() common.DeviceManager {
//	return sd
//}
//func (sd *SystemDevice) GetPrinterManager() common.PrinterManager {
//	return sd
//}
//func (sd *SystemDevice) GetReaderManager() common.ReaderManager {
//	return sd
//}
//func (sd *SystemDevice) GetPinPadManager() common.PinPadManager {
//	return sd
//}
//func (sd *SystemDevice) GetValidatorManager() common.ValidatorManager {
//	return sd
//}

// Implementation of common.SystemManager

func (sd *SystemDevice) Terminate(name string, query *common.SystemQuery) error {
	sd.state = common.SysStateUndefined
	var err error
	if sd.driver != nil {
		err = sd.driver.StopDevice()
	}
	reply := &common.SystemReply{}
	reply.Command = common.CmdSystemTerminate
	reply.State = sd.state
	if err != nil {
		reply.Message = err.Error()
	}
	//core.SendQuitSignal(100)
	if sd.system != nil {
		_ = sd.system.SystemReply(sd.devName, reply)
	}
	return err
}

func (sd *SystemDevice) SysInform(name string, query *common.SystemQuery) error {
	reply := &common.SystemReply{}
	reply.Command = common.CmdSystemInform
	reply.State = sd.state
	if sd.system != nil {
		_ = sd.system.SystemReply(sd.devName, reply)
	}
	return nil
}

func (sd *SystemDevice) SysStart(name string, query *common.SystemConfig) error {
	sd.state = common.SysStateUndefined
	var err error
	if sd.driver != nil {
		err = sd.driver.StartDevice(query)
		if err == nil {
			sd.state = common.SysStateRunning
		} else {
			sd.state = common.SysStateFailed
		}
	}
	reply := &common.SystemReply{}
	reply.Command = common.CmdSystemStart
	reply.State = sd.state
	if err != nil {
		reply.Message = err.Error()
	}
	if sd.system != nil {
		_ = sd.system.SystemReply(sd.devName, reply)
	}
	return err
}

func (sd *SystemDevice) SysStop(name string, query *common.SystemQuery) error {
	sd.state = common.SysStateUndefined
	var err error
	if sd.driver != nil {
		err = sd.driver.StopDevice()
		if err == nil {
			sd.state = common.SysStateStopped
		}
	}
	reply := &common.SystemReply{}
	reply.Command = common.CmdSystemStop
	reply.State = sd.state
	if err != nil {
		reply.Message = err.Error()
	}
	if sd.system != nil {
		_ = sd.system.SystemReply(sd.devName, reply)
	}
	return err
}

func (sd *SystemDevice) SysRestart(name string, query *common.SystemConfig) error {
	sd.state = common.SysStateUndefined
	var err error
	if sd.driver != nil {
		err = sd.driver.StopDevice()
		if err == nil {
			sd.state = common.SysStateStopped
		}
		err = sd.driver.StartDevice(query)
		if err == nil {
			sd.state = common.SysStateRunning
		} else {
			sd.state = common.SysStateFailed
		}
	}
	reply := &common.SystemReply{}
	reply.Command = common.CmdSystemRestart
	reply.State = sd.state
	if err != nil {
		reply.Message = err.Error()
	}
	if sd.system != nil {
		_ = sd.system.SystemReply(sd.devName, reply)
	}
	return err
}

// Implementation of common.DeviceManager

func (sd *SystemDevice) Cancel(name string, query *common.DeviceQuery) error {
	if sd.manager.Device != nil {
		return sd.manager.Device.Cancel(name, query)
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}
func (sd *SystemDevice) Reset(name string, query *common.DeviceQuery) error {
	if sd.manager.Device != nil {
		return sd.manager.Device.Reset(name, query)
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}
func (sd *SystemDevice) Status(name string, query *common.DeviceQuery) error {
	if sd.manager.Device != nil {
		return sd.manager.Device.Status(name, query)
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}
func (sd *SystemDevice) RunAction(name string, query *common.DeviceQuery) error {
	if sd.manager.Device != nil {
		return sd.manager.Device.RunAction(name, query)
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}
func (sd *SystemDevice) StopAction(name string, query *common.DeviceQuery) error {
	if sd.manager.Device != nil {
		return sd.manager.Device.StopAction(name, query)
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}

// Implementation of common.PrinterManager

func (sd *SystemDevice) InitPrinter(name string, query *common.PrinterSetup) error {
	if sd.manager.Printer != nil {
		return sd.manager.Printer.InitPrinter(name, query)
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}
func (sd *SystemDevice) PrintText(name string, query *common.PrinterQuery) error {
	if sd.manager.Printer != nil {
		return sd.manager.Printer.PrintText(name, query)
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}

// Implementation of common.ReaderManager

func (sd *SystemDevice) EnterCard(name string, query *common.DeviceQuery) error {
	if sd.manager.Reader != nil {
		return sd.manager.Reader.EnterCard(name, query)
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}
func (sd *SystemDevice) EjectCard(name string, query *common.DeviceQuery) error {
	if sd.manager.Reader != nil {
		return sd.manager.Reader.EjectCard(name, query)
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}
func (sd *SystemDevice) CaptureCard(name string, query *common.DeviceQuery) error {
	if sd.manager.Reader != nil {
		return sd.manager.Reader.CaptureCard(name, query)
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}
func (sd *SystemDevice) ReadCard(name string, query *common.DeviceQuery) error {
	if sd.manager.Reader != nil {
		return sd.manager.Reader.ReadCard(name, query)
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}
func (sd *SystemDevice) ChipGetATR(name string, query *common.DeviceQuery) error {
	if sd.manager.Reader != nil {
		return sd.manager.Reader.ChipGetATR(name, query)
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}
func (sd *SystemDevice) ChipPowerOff(name string, query *common.DeviceQuery) error {
	if sd.manager.Reader != nil {
		return sd.manager.Reader.ChipPowerOff(name, query)
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}
func (sd *SystemDevice) ChipCommand(name string, query *common.ReaderChipQuery) error {
	if sd.manager.Reader != nil {
		return sd.manager.Reader.ChipCommand(name, query)
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}

// Implementation of common.ValidatorManager

func (sd *SystemDevice) InitValidator(name string, query *common.ValidatorQuery) error {
	if sd.manager.Validator != nil {
		return sd.manager.Validator.InitValidator(name, query)
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}
func (sd *SystemDevice) DoValidate(name string, query *common.ValidatorQuery) error {
	if sd.manager.Validator != nil {
		return sd.manager.Validator.DoValidate(name, query)
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}
func (sd *SystemDevice) NoteAccept(name string, query *common.ValidatorQuery) error {
	if sd.manager.Validator != nil {
		return sd.manager.Validator.NoteAccept(name, query)
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}
func (sd *SystemDevice) NoteReturn(name string, query *common.ValidatorQuery) error {
	if sd.manager.Validator != nil {
		return sd.manager.Validator.NoteReturn(name, query)
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}
func (sd *SystemDevice) StopValidate(name string, query *common.ValidatorQuery) error {
	if sd.manager.Validator != nil {
		return sd.manager.Validator.StopValidate(name, query)
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}
func (sd *SystemDevice) CheckValidator(name string, query *common.ValidatorQuery) error {
	if sd.manager.Validator != nil {
		return sd.manager.Validator.CheckValidator(name, query)
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}
func (sd *SystemDevice) ClearValidator(name string, query *common.ValidatorQuery) error {
	if sd.manager.Validator != nil {
		return sd.manager.Validator.ClearValidator(name, query)
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}

// Implementation of common.PinPadManager

func (sd *SystemDevice) ReadPIN(name string, query *common.ReaderPinQuery) error {
	if sd.manager.PinPad != nil {
		return sd.manager.PinPad.ReadPIN(name, query)
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}
func (sd *SystemDevice) LoadMasterKey(name string, query *common.ReaderPinQuery) error {
	if sd.manager.PinPad != nil {
		return sd.manager.PinPad.LoadMasterKey(name, query)
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}
func (sd *SystemDevice) LoadWorkKey(name string, query *common.ReaderPinQuery) error {
	if sd.manager.PinPad != nil {
		return sd.manager.PinPad.LoadWorkKey(name, query)
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}
func (sd *SystemDevice) TestMasterKey(name string, query *common.ReaderPinQuery) error {
	if sd.manager.PinPad != nil {
		return sd.manager.PinPad.TestMasterKey(name, query)
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}
func (sd *SystemDevice) TestWorkKey(name string, query *common.ReaderPinQuery) error {
	if sd.manager.PinPad != nil {
		return sd.manager.PinPad.TestWorkKey(name, query)
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}

/*
// Implementation of common.SystemCallback

func (sd *SystemDevice) SystemReply(name string, reply *common.SystemReply) error {
	if sd.callback != nil {
		sysCb := sd.callback.GetSystemCallback()
		if sysCb != nil {
			return sysCb.SystemReply(sd.devName, reply)
		}
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}
func (sd *SystemDevice) SystemHealth(name string, reply *common.SystemHealth) error {
	if sd.callback != nil {
		sysCb := sd.callback.GetSystemCallback()
		if sysCb != nil {
			return sysCb.SystemHealth(sd.devName, reply)
		}
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}

// Implementation of common.DeviceCallback
func (sd *SystemDevice) DeviceReply(name string, reply *common.DeviceReply) error {
	return sd.encodeReply(duplex.ScopeDevice, common.CmdDeviceReply, reply)
}
func (sd *SystemDevice) ExecuteError(name string, reply *common.DeviceError) error {
	return sd.encodeReply(duplex.ScopeDevice, common.CmdExecuteError, reply)
}
func (sd *SystemDevice) StateChanged(name string, reply *common.DeviceState) error {
	return sd.encodeReply(duplex.ScopeDevice, common.CmdStateChanged, reply)
}
func (sd *SystemDevice) ActionPrompt(name string, reply *common.DevicePrompt) error {
	return sd.encodeReply(duplex.ScopeDevice, common.CmdActionPrompt, reply)
}
func (sd *SystemDevice) ReaderReturn(name string, reply *common.DeviceInform) error {
	return sd.encodeReply(duplex.ScopeDevice, common.CmdReaderReturn, reply)
}

// Implementation of common.PrinterCallback
func (sd *SystemDevice) PrinterProgress(name string, reply *common.PrinterProgress) error {
	return sd.encodeReply(duplex.ScopePrinter, common.CmdPrinterProgress, reply)
}

// Implementation of common.ReaderCallback
func (sd *SystemDevice) CardPosition(name string, reply *common.ReaderCardPos) error {
	return sd.encodeReply(duplex.ScopeReader, common.CmdCardPosition, reply)
}
func (sd *SystemDevice) CardDescription(name string, reply *common.ReaderCardInfo) error {
	return sd.encodeReply(duplex.ScopeReader, common.CmdCardDescription, reply)
}
func (sd *SystemDevice) ChipResponse(name string, reply *common.ReaderChipReply) error {
	return sd.encodeReply(duplex.ScopeReader, common.CmdChipResponse, reply)
}

// Implementation of common.ValidatorCallback
func (sd *SystemDevice) NoteAccepted(name string, reply *common.ValidatorAccept) error {
	return sd.encodeReply(duplex.ScopeValidator, common.CmdNoteAccepted, reply)
}
func (sd *SystemDevice) CashIsStored(name string, reply *common.ValidatorAccept) error {
	return sd.encodeReply(duplex.ScopeValidator, common.CmdCashIsStored, reply)
}
func (sd *SystemDevice) CashReturned(name string, reply *common.ValidatorAccept) error {
	return sd.encodeReply(duplex.ScopeValidator, common.CmdCashReturned, reply)
}
func (sd *SystemDevice) ValidatorStore(name string, reply *common.ValidatorStore) error {
	return sd.encodeReply(duplex.ScopeValidator, common.CmdValidatorStore, reply)
}

// Implementation of common.ReaderCallback
func (sd *SystemDevice) PinPadReply(name string, reply *common.ReaderPinReply) error {
	return sd.encodeReply(duplex.ScopePinPad, common.CmdPinPadReply, reply)
}

// Common function for reply encoding
func (sd *SystemDevice) encodeReply(scope duplex.PacketScope, cmd string, reply interface{}) error {
	dump, err := json.Marshal(reply)
	if err != nil {
		return err
	}
	if sd.log != nil {
		sd.log.Dump("SystemDevice dev:%s send scope:%s, cmd:%s pack:%s",
			sd.devName, duplex.GetScopeName(scope), cmd, string(dump))
	}
	return err
}
*/
