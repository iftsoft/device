package driver

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
)

type DeviceDriver interface {
	InitDevice(manager interface{}) common.DevScopeMask
	StartDevice(cfg *config.DeviceConfig) error
	DeviceTimer(unix int64) error
	CheckDevice(metrics *common.SystemMetrics) error
	StopDevice() error
}
