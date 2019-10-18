package driver

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/duplex"
)

type DeviceDriver interface {
	InitDevice(manager interface{}, cfg *config.DeviceConfig, info *duplex.GreetingInfo) error
	StartDevice() error
	StopDevice() error
	DeviceTimer(unix int64) error
	CheckDevice(metrics *common.SystemMetrics) error
}
