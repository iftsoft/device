package handler

import (
	"github.com/iftsoft/device/core"
)

type PluginManager interface {
	Enabled(on bool)
	Connected(on bool)
	OnTimerTick()
}

type PluginCreator interface {
	GetPluginName() string
	CreatePlugin(devName string, proxy interface{}, log *core.LogAgent) (error, PluginManager)
}
