package driver

import "github.com/iftsoft/device/config"

type DeviceDriver interface {
	InitDevice(manager interface{}) error
	StartDevice(cfg *config.DeviceConfig) error
	DeviceTimer(unix int64) error
	StopDevice() error
	ClearDevice() error
}
