package system

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/driver"
	"github.com/iftsoft/device/duplex"
	"github.com/iftsoft/device/proxy"
)

type SystemDevice struct {
	driver driver.DeviceDriver
	duplex *duplex.DuplexClient
	//	client   *duplex.DuplexClientConfig
	config *config.DeviceConfig
	//	system   common.SystemCallback
	system *proxy.SystemClient
	device *proxy.DeviceClient
	done   chan struct{}
	log    *core.LogAgent
}

func NewBaseDevice(cfg *config.AppConfig) *SystemDevice {
	sd := SystemDevice{
		driver: nil,
		duplex: duplex.NewDuplexClient(&cfg.Client),
		config: &cfg.Device,
		system: proxy.NewSystemClient(),
		device: proxy.NewDeviceClient(),
		log:    core.GetLogAgent(core.LogLevelTrace, "Device"),
	}
	return &sd
}

func (sd *SystemDevice) InitDevice(worker interface{}) {
	driver, okDrv := worker.(driver.DeviceDriver)
	if okDrv {
		sd.driver = driver
	}
	sd.system.Init(sd.duplex, sd, sd.log)
	sd.duplex.AddScopeItem(sd.system.GetScopeItem())

	device, okDev := worker.(common.DeviceManager)
	if okDev {
		sd.device.Init(sd.duplex, device, sd.log)
		sd.duplex.AddScopeItem(sd.device.GetScopeItem())
	}
}

// Implementation of common.SystemManager
//
func (sd *SystemDevice) Config(name string, query *common.SystemQuery) error {
	return nil
}

func (sd *SystemDevice) Inform(name string, query *common.SystemQuery) error {
	return nil
}

func (sd *SystemDevice) Start(name string, query *common.SystemQuery) error {
	return nil
}

func (sd *SystemDevice) Stop(name string, query *common.SystemQuery) error {
	return nil
}

func (sd *SystemDevice) Restart(name string, query *common.SystemQuery) error {
	return nil
}
