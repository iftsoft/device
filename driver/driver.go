package driver

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
)

type DeviceDriver interface {
	InitDevice(manager interface{}, cfg *config.DeviceConfig) common.DevScopeMask
	StartDevice() error
	StopDevice() error
	DeviceTimer(unix int64) error
	CheckDevice(metrics *common.SystemMetrics) error
}
