package driver

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/duplex"
	"github.com/iftsoft/device/dbase"
)

type Context struct {
	Manager  interface{}
	Storage  dbase.DBaseLinker
	Config   *config.DeviceConfig
	Greeting *duplex.GreetingInfo
}

type DeviceDriver interface {
	InitDevice(context *Context) error
	StartDevice() error
	StopDevice() error
	DeviceTimer(unix int64) error
	CheckDevice(metrics *common.SystemMetrics) error
}
