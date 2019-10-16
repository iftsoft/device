package handler

import (
	"github.com/iftsoft/device/core"
)

type ReflexManager interface {
	Enabled(on bool)
	Connected(on bool)
	OnTimerTick()
}

type ReflexInfo struct {
	ReflexName string
	Mandatory  bool
}

type ReflexCreator interface {
	GetReflexInfo() *ReflexInfo
	CreateReflex(devName string, proxy interface{}, log *core.LogAgent) (error, ReflexManager)
}
