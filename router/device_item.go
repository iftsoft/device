package router

import (
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/dbase"
	driver2 "github.com/iftsoft/device/driver"
	"sync"
	"time"
)

type DeviceItem struct {
	devName  string
	greeting *common.GreetingInfo
	state    common.EnumSystemState
	error    common.EnumSystemError
	driver   driver2.DeviceDriver
	callback common.ComplexCallback
	manager  common.ManagerSet
	system   common.SystemCallback
	//device    common.DeviceManager
	//printer   common.PrinterManager
	//reader    common.ReaderManager
	//validator common.ValidatorManager
	//pinpad    common.PinPadManager
	config  *config.DeviceConfig
	storage *dbase.DBaseStore
	log     *core.LogAgent
	done    chan struct{}
	wg      sync.WaitGroup
	checkTm time.Time
}

func NewDeviceItem(cfg *config.DeviceConfig) *DeviceItem {
	if cfg == nil {
		return nil
	}
	sd := DeviceItem{
		devName:  cfg.DevName,
		greeting: &common.GreetingInfo{},
		state:    common.SysStateUndefined,
		error:    common.SysErrSuccess,
		driver:   nil,
		system:   nil,
		//device:    nil,
		//printer:   nil,
		//reader:    nil,
		//validator: nil,
		//pinpad:    nil,
		config: cfg,
		//storage:   dbase.GetNewDBaseStore(cfg.Storage),
		log:  core.GetLogAgent(core.LogLevelTrace, cfg.DevName),
		done: make(chan struct{}),
	}
	return &sd
}

func (sd *DeviceItem) InitDevice(driver driver2.DeviceDriver) error {
	if sd.devName == "" {
		return errors.New("device name is not set in client config")
	}

	// Setup Device driver interface
	sd.driver = driver
	sd.system = sd.callback.GetSystemCallback()
	context := driver2.Context{
		DevName:  sd.devName,
		Complex:  sd.callback,
		Storage:  sd.storage,
		Config:   sd.config,
		Greeting: sd.greeting,
	}
	manager := driver.InitDevice(&context)
	sd.greeting.Supported = sd.manager.InitManagers(manager)
	sd.greeting.Supported |= common.ScopeFlagSystem

	return nil
}

func (sd *DeviceItem) StartDeviceLoop() {
	sd.log.Info("Starting system device")
	go sd.deviceLoop(&sd.wg)
}

func (sd *DeviceItem) StopDeviceLoop() {
	sd.log.Info("Stopping system device")
	close(sd.done)
	sd.wg.Wait()
}

func (sd *DeviceItem) deviceLoop(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	sd.log.Debug("System device loop is started")
	defer sd.log.Debug("System device loop is stopped")

	sd.state = common.SysStateUndefined
	//if sd.config.Common.AutoLoad {
	//	err := sd.driver.StartDevice(nil)
	//	if err == nil {
	//		sd.state = common.SysStateRunning
	//	}
	//}
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

func (sd *DeviceItem) onTimerTick(tm time.Time) {
	sd.log.Trace("System device %s onTimerTick %s", sd.devName, tm.Format(time.StampMilli))
	if sd.driver != nil && sd.state == common.SysStateRunning {
		_ = sd.driver.DeviceTimer(tm.Unix())
	}
	delta := int(tm.Sub(sd.checkTm).Seconds())
	if delta == 60 {
		sd.sendDeviceMetrics(tm)
	}
}

func (sd *DeviceItem) sendDeviceMetrics(tm time.Time) {
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

// Implementation of common.SystemManager

func (sd *DeviceItem) Terminate(name string, query *common.SystemQuery) error {
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

func (sd *DeviceItem) SysInform(name string, query *common.SystemQuery) error {
	reply := &common.SystemReply{}
	reply.Command = common.CmdSystemInform
	reply.State = sd.state
	if sd.system != nil {
		_ = sd.system.SystemReply(sd.devName, reply)
	}
	return nil
}

func (sd *DeviceItem) SysStart(name string, query *common.SystemConfig) error {
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

func (sd *DeviceItem) SysStop(name string, query *common.SystemQuery) error {
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

func (sd *DeviceItem) SysRestart(name string, query *common.SystemConfig) error {
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

/*
// Implementation of common.SystemCallback

func (sd *DeviceItem) SystemReply(name string, reply *common.SystemReply) error {
	if sd.callback != nil {
		sysCb := sd.callback.GetSystemCallback()
		if sysCb != nil {
			return sysCb.SystemReply(sd.devName, reply)
		}
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}
func (sd *DeviceItem) SystemHealth(name string, reply *common.SystemHealth) error {
	if sd.callback != nil {
		sysCb := sd.callback.GetSystemCallback()
		if sysCb != nil {
			return sysCb.SystemHealth(sd.devName, reply)
		}
	}
	return common.NewError(common.DevErrorNotImplemented, "")
}

// Implementation of common.DeviceCallback
func (sd *DeviceItem) DeviceReply(name string, reply *common.DeviceReply) error {
	return sd.encodeReply(duplex.ScopeDevice, common.CmdDeviceReply, reply)
}
func (sd *DeviceItem) ExecuteError(name string, reply *common.DeviceError) error {
	return sd.encodeReply(duplex.ScopeDevice, common.CmdExecuteError, reply)
}
func (sd *DeviceItem) StateChanged(name string, reply *common.DeviceState) error {
	return sd.encodeReply(duplex.ScopeDevice, common.CmdStateChanged, reply)
}
func (sd *DeviceItem) ActionPrompt(name string, reply *common.DevicePrompt) error {
	return sd.encodeReply(duplex.ScopeDevice, common.CmdActionPrompt, reply)
}
func (sd *DeviceItem) ReaderReturn(name string, reply *common.DeviceInform) error {
	return sd.encodeReply(duplex.ScopeDevice, common.CmdReaderReturn, reply)
}

// Implementation of common.PrinterCallback
func (sd *DeviceItem) PrinterProgress(name string, reply *common.PrinterProgress) error {
	return sd.encodeReply(duplex.ScopePrinter, common.CmdPrinterProgress, reply)
}

// Implementation of common.ReaderCallback
func (sd *DeviceItem) CardPosition(name string, reply *common.ReaderCardPos) error {
	return sd.encodeReply(duplex.ScopeReader, common.CmdCardPosition, reply)
}
func (sd *DeviceItem) CardDescription(name string, reply *common.ReaderCardInfo) error {
	return sd.encodeReply(duplex.ScopeReader, common.CmdCardDescription, reply)
}
func (sd *DeviceItem) ChipResponse(name string, reply *common.ReaderChipReply) error {
	return sd.encodeReply(duplex.ScopeReader, common.CmdChipResponse, reply)
}

// Implementation of common.ValidatorCallback
func (sd *DeviceItem) NoteAccepted(name string, reply *common.ValidatorAccept) error {
	return sd.encodeReply(duplex.ScopeValidator, common.CmdNoteAccepted, reply)
}
func (sd *DeviceItem) CashIsStored(name string, reply *common.ValidatorAccept) error {
	return sd.encodeReply(duplex.ScopeValidator, common.CmdCashIsStored, reply)
}
func (sd *DeviceItem) CashReturned(name string, reply *common.ValidatorAccept) error {
	return sd.encodeReply(duplex.ScopeValidator, common.CmdCashReturned, reply)
}
func (sd *DeviceItem) ValidatorStore(name string, reply *common.ValidatorStore) error {
	return sd.encodeReply(duplex.ScopeValidator, common.CmdValidatorStore, reply)
}

// Implementation of common.ReaderCallback
func (sd *DeviceItem) PinPadReply(name string, reply *common.ReaderPinReply) error {
	return sd.encodeReply(duplex.ScopePinPad, common.CmdPinPadReply, reply)
}

// Common function for reply encoding
func (sd *DeviceItem) encodeReply(scope duplex.PacketScope, cmd string, reply interface{}) error {
	dump, err := json.Marshal(reply)
	if err != nil {
		return err
	}
	if sd.log != nil {
		sd.log.Dump("DeviceItem dev:%s send scope:%s, cmd:%s pack:%s",
			sd.devName, duplex.GetScopeName(scope), cmd, string(dump))
	}
	return err
}
*/
