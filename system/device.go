package system

import (
	"encoding/json"
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/driver"
	"github.com/iftsoft/device/duplex"
	"github.com/iftsoft/device/proxy"
	"sync"
	"time"
)

type SystemDevice struct {
	devName string
	state   common.EnumSystemState
	driver  driver.DeviceDriver
	duplex  *duplex.DuplexClient
	config  *config.DeviceConfig
	system  *proxy.SystemClient
	device  *proxy.DeviceClient
	log     *core.LogAgent
	done    chan struct{}
	wg      sync.WaitGroup
	checkTm time.Time
}

func NewSystemDevice(cfg *config.AppConfig) *SystemDevice {
	if cfg == nil {
		return nil
	}
	sd := SystemDevice{
		devName: cfg.Client.DevName,
		state:   common.SysStateUndefined,
		driver:  nil,
		duplex:  duplex.NewDuplexClient(&cfg.Client),
		config:  &cfg.Device,
		system:  proxy.NewSystemClient(),
		device:  proxy.NewDeviceClient(),
		log:     core.GetLogAgent(core.LogLevelTrace, "Device"),
		done:    make(chan struct{}),
	}
	return &sd
}

func (sd *SystemDevice) InitDevice(worker interface{}) error {
	if sd.devName == "" {
		return errors.New("device name is not set in client config")
	}
	// Setup System scope interface
	sd.system.Init(sd, sd.log)
	sd.duplex.AddScopeItem(sd.system.GetScopeItem())

	// Setup Device scope interface
	device, okDev := worker.(common.DeviceManager)
	if okDev {
		sd.device.Init(device, sd.log)
		sd.duplex.AddScopeItem(sd.device.GetScopeItem())
	}
	// Setup Device driver interface
	drv, okDrv := worker.(driver.DeviceDriver)
	if okDrv {
		sd.driver = drv
		return drv.InitDevice(sd)
	}
	return errors.New("device driver is not implemented")
}

func (sd *SystemDevice) StartDevice() {
	sd.log.Info("Starting system device")
	sd.duplex.StartClient(&sd.wg)
	go sd.deviceLoop(&sd.wg)
}

func (sd *SystemDevice) StopDevice() {
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

	if sd.config.Common.AutoLoad {
		err := sd.driver.StartDevice(sd.config)
		if err == nil {
			sd.state = common.SysStateRunning
		}
	}
	defer sd.driver.StopDevice()

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
	err := sd.driver.CheckDevice(&health.Metrics)
	if err != nil {
		health.Error = common.SysErrDeviceFail
	}
	err = sd.SystemHealth(sd.devName, health)
}

// Implementation of common.SystemManager
//
func (sd *SystemDevice) Config(name string, query *common.SystemQuery) error {
	reply := &common.SystemReply{}
	reply.Command = common.CmdSystemConfig
	reply.State = sd.state
	return sd.SystemReply(name, reply)
}

func (sd *SystemDevice) Inform(name string, query *common.SystemQuery) error {
	reply := &common.SystemReply{}
	reply.Command = common.CmdSystemInform
	reply.State = sd.state
	return sd.SystemReply(name, reply)
}

func (sd *SystemDevice) Start(name string, query *common.SystemQuery) error {
	sd.state = common.SysStateUndefined
	err := sd.driver.StartDevice(sd.config)
	if err == nil {
		sd.state = common.SysStateRunning
	}
	reply := &common.SystemReply{}
	reply.Command = common.CmdSystemStart
	reply.State = sd.state
	if err != nil {
		reply.Message = err.Error()
	}
	return sd.SystemReply(name, reply)
}

func (sd *SystemDevice) Stop(name string, query *common.SystemQuery) error {
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

func (sd *SystemDevice) Restart(name string, query *common.SystemQuery) error {
	sd.state = common.SysStateUndefined
	var err error
	if sd.driver != nil {
		err = sd.driver.StopDevice()
		if err == nil {
			sd.state = common.SysStateStopped
		}
		err = sd.driver.StartDevice(sd.config)
		if err == nil {
			sd.state = common.SysStateRunning
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
