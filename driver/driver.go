package driver

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/dbase"
)

type Context struct {
	DevName  string
	Complex  common.ComplexCallback
	Storage  dbase.DBaseLinker
	Config   *config.DeviceConfig
	Greeting *common.GreetingInfo
}

type DeviceDriver interface {
	InitDevice(context *Context) common.ComplexManager
	StartDevice(query *common.SystemConfig) error
	StopDevice() error
	DeviceTimer(unix int64) error
	CheckDevice(metrics *common.SystemMetrics) error
}
