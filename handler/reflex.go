package handler

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
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
	DevType    common.DevTypeMask  // Applied for device types
	Supported  common.DevScopeMask // Callback interfaces that reflex supported
	Required   common.DevScopeMask // Manager interfaces that reflex required
}

func (ri *ReflexInfo) IsMatched(gi *common.GreetingInfo) bool {
	if (ri.DevType & gi.DevType) == ri.DevType {
		return true
	}
	return false
}

type ReflexSet struct {
	reflexMap map[string]ReflexCreator
	log       *core.LogAgent
}

func (rs *ReflexSet) initReflexSet() {
	rs.reflexMap = make(map[string]ReflexCreator)
	rs.log = core.GetLogAgent(core.LogLevelTrace, "Reflex")
}

func (rs *ReflexSet) registerFactory(fact ReflexCreator) {
	if fact == nil {
		return
	}
	info := fact.GetReflexInfo()
	if info == nil {
		return
	}
	rs.log.Debug("ReflexSet.RegisterFactory for reflex:%s", info.ReflexName)
	rs.reflexMap[info.ReflexName] = fact
}

func (rs *ReflexSet) attachReflexes(handler *DeviceHandler, proxy *HandlerProxy, greet *common.GreetingInfo) {
	// Attach mandatory reflexes
	for refName, factory := range rs.reflexMap {
		info := factory.GetReflexInfo()
		if info.Mandatory && info.IsMatched(greet) {
			rs.log.Debug("ReflexSet.AttachReflexes is attaching reflex:%s to device:%s",
				refName, handler.devName)
			err, reflex := factory.CreateReflex(handler.devName, proxy, rs.log)
			if err == nil {
				err = handler.AttachReflex(refName, reflex)
			}
		}
	}
	return
}
