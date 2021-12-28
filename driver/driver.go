package driver

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/dbase"
)

type Specification struct {
	ModelName   string              `json:"model_name"`  // Unique model name
	Versions    []string            `json:"versions"`    // list of model versions
	Description string              `json:"description"` // Model description
	DeviceType  common.DevTypeMask  `json:"device_type"` // Implemented device types
	Supported   common.DevScopeMask `json:"supported"`   // Manager interfaces that driver supported
	Required    common.DevScopeMask `json:"required"`    // Callback interfaces that driver required
}

type Context struct {
	DevName string
	Complex common.ComplexCallback
	Storage dbase.DBaseLinker
	Config  *config.DeviceConfig
}

type DeviceDriver interface {
	InitDevice(context *Context) common.ComplexManager
	StartDevice(query *common.SystemConfig) error
	StopDevice() error
	DeviceTimer(unix int64) error
	CheckDevice(metrics *common.SystemMetrics) error
}

type DrvFactory func() DeviceDriver
