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
	channel CallbackChannel
	config  *config.DeviceConfig
	log     *core.LogAgent
	query   chan common.Packet
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
		log:     core.GetLogAgent(core.LogLevelTrace, name),
		query:   make(chan common.Packet, 8),
		done:    make(chan struct{}),
	}
	return &sd
}

func (sd *SystemDevice) InitDevice(drv driver.DeviceDriver, ctx *driver.Context) {
	sd.driver = drv
	sd.config = ctx.Config
	sd.storage = ctx.Storage
	manager := drv.InitDevice(ctx)
	sd.manager.InitManagers(manager)
	sd.manager.System = sd
}

func (sd *SystemDevice) RunCommand(command, devName string, query interface{}) error {
	packet := common.Packet{
		Command: command,
		DevName: devName,
		Content: query,
	}
	sd.query <- packet
	return nil
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
		case pack := <-sd.query:
			_ = sd.manager.RunCommand(pack)
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
	err = sd.channel.SystemHealth(sd.devName, health)
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
	_ = sd.channel.SystemReply(sd.devName, reply)
	return err
}

func (sd *SystemDevice) SysInform(name string, query *common.SystemQuery) error {
	reply := &common.SystemReply{}
	reply.Command = common.CmdSystemInform
	reply.State = sd.state
	_ = sd.channel.SystemReply(sd.devName, reply)
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
	_ = sd.channel.SystemReply(sd.devName, reply)
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
	_ = sd.channel.SystemReply(sd.devName, reply)
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
	_ = sd.channel.SystemReply(sd.devName, reply)
	return err
}

/*
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
*/
