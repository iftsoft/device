package driver

/*
import (
	"encoding/json"
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/dbase"
	"github.com/iftsoft/device/duplex"
	"github.com/iftsoft/device/proxy"
	"sync"
	"time"
)

type SystemDevice struct {
	devName   string
	greeting  *common.GreetingInfo
	state     common.EnumSystemState
	error     common.EnumSystemError
	driver    DeviceDriver
	duplex    *duplex.DuplexClient
	config    *config.DeviceConfig
	system    *proxy.SystemClient
	device    *proxy.DeviceClient
	printer   *proxy.PrinterClient
	reader    *proxy.ReaderClient
	validator *proxy.ValidatorClient
	pinpad    *proxy.PinPadClient
	storage   *dbase.DBaseStore
	log       *core.LogAgent
	done      chan struct{}
	wg        sync.WaitGroup
	checkTm   time.Time
}

func NewSystemDevice(cfg *config.AppConfig) *SystemDevice {
	if cfg == nil {
		return nil
	}
	sd := SystemDevice{
		devName:   cfg.Duplex.DevName,
		greeting:  &common.GreetingInfo{},
		state:     common.SysStateUndefined,
		error:     common.SysErrSuccess,
		driver:    nil,
		duplex:    duplex.NewDuplexClient(cfg.Duplex),
		config:    cfg.Devices[0],
		system:    proxy.NewSystemClient(),
		device:    proxy.NewDeviceClient(),
		printer:   proxy.NewPrinterClient(),
		reader:    proxy.NewReaderClient(),
		validator: proxy.NewValidatorClient(),
		pinpad:    proxy.NewPinPadClient(),
		storage:   dbase.GetNewDBaseStore(cfg.Storage),
		log:       core.GetLogAgent(core.LogLevelTrace, "Device"),
		done:      make(chan struct{}),
	}
	return &sd
}

func (sd *SystemDevice) InitDevice(worker interface{}) error {
	if sd.devName == "" {
		return errors.New("device name is not set in client config")
	}
	// Setup System scope interface
	sd.system.Init(sd, sd.log)
	sd.duplex.AddDispatcher(duplex.ScopeSystem, sd.system.GetDispatcher())
	sd.greeting.Supported = common.ScopeFlagSystem

	// Setup Device scope interface
	if device, ok := worker.(common.DeviceManager); ok {
		sd.device.Init(device, sd.log)
		sd.duplex.AddDispatcher(duplex.ScopeDevice, sd.device.GetDispatcher())
		sd.greeting.Supported |= common.ScopeFlagDevice
	}
	// Setup Printer scope interface
	if printer, ok := worker.(common.PrinterManager); ok {
		sd.printer.Init(printer, sd.log)
		sd.duplex.AddDispatcher(duplex.ScopePrinter, sd.printer.GetDispatcher())
		sd.greeting.Supported |= common.ScopeFlagPrinter
	}
	// Setup Reader scope interface
	if reader, ok := worker.(common.ReaderManager); ok {
		sd.reader.Init(reader, sd.log)
		sd.duplex.AddDispatcher(duplex.ScopeReader, sd.reader.GetDispatcher())
		sd.greeting.Supported |= common.ScopeFlagReader
	}
	// Setup Validator scope interface
	if valid, ok := worker.(common.ValidatorManager); ok {
		sd.validator.Init(valid, sd.log)
		sd.duplex.AddDispatcher(duplex.ScopeValidator, sd.validator.GetDispatcher())
		sd.greeting.Supported |= common.ScopeFlagValidator
	}
	// Setup PinPad scope interface
	if pinpad, ok := worker.(common.PinPadManager); ok {
		sd.pinpad.Init(pinpad, sd.log)
		sd.duplex.AddDispatcher(duplex.ScopePinPad, sd.pinpad.GetDispatcher())
		sd.greeting.Supported |= common.ScopeFlagPinPad
	}
	// Setup Device driver interface
	if drv, ok := worker.(DeviceDriver); ok {
		sd.driver = drv
		context := Context{
			DevName: sd.devName,
			//Manager:  sd,
			Storage: sd.storage,
			Config:  sd.config,
			//Greeting: sd.greeting,
		}
		drv.InitDevice(&context)
		return nil
	}
	return errors.New("device driver is not implemented")
}

func (sd *SystemDevice) StartDeviceLoop() {
	sd.log.Info("Starting system device")
	sd.duplex.StartClient(&sd.wg, sd.greeting)
	go sd.deviceLoop(&sd.wg)
}

func (sd *SystemDevice) StopDeviceLoop() {
	sd.log.Info("Stopping system device")
	close(sd.done)
	sd.duplex.StopClient(&sd.wg)
	sd.wg.Wait()
}

func (sd *SystemDevice) deviceLoop(wg *sync.WaitGroup) {
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
	err = sd.SystemHealth(sd.devName, health)
}

// Implementation of common.SystemManager
//
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
	core.SendQuitSignal(100)
	return sd.SystemReply(name, reply)
}

func (sd *SystemDevice) SysInform(name string, query *common.SystemQuery) error {
	reply := &common.SystemReply{}
	reply.Command = common.CmdSystemInform
	reply.State = sd.state
	return sd.SystemReply(name, reply)
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
	return sd.SystemReply(name, reply)
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
	return sd.SystemReply(name, reply)
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
	return sd.SystemReply(name, reply)
}

// Implementation of common.SystemCallback
func (sd *SystemDevice) SystemReply(name string, reply *common.SystemReply) error {
	return sd.encodeReply(duplex.ScopeSystem, common.CmdSystemReply, reply)
}
func (sd *SystemDevice) SystemHealth(name string, reply *common.SystemHealth) error {
	return sd.encodeReply(duplex.ScopeSystem, common.CmdSystemHealth, reply)
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
	pack := duplex.NewPacket(scope, sd.devName, cmd, dump)
	if sd.duplex != nil {
		err = sd.duplex.SendPacket(pack)
	}
	return err
}
*/
