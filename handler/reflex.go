package handler

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type ReflexManager interface {
	Enabled(on bool)
	Connected(on bool)
	OnTimerTick()
}

type ReflexCreator interface {
	GetReflexInfo() *ReflexInfo
	CreateReflex(devName string, proxy interface{}, log *core.LogAgent) (error, ReflexManager)
}

type ReflexInfo struct {
	ReflexName string
	Mandatory  bool
	Supported  common.DevScopeMask	// Callback interfaces that reflex supported
	Required   common.DevScopeMask	// Manager interfaces that reflex required
}

func (ri *ReflexInfo) IsMatched(gi *duplex.GreetingInfo) bool {
	if  (ri.Supported & gi.Required) == gi.Required &&
		(gi.Supported & ri.Required) == ri.Required {
		return true
	}
	return false
}