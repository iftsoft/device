package handler

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"sync"
)

type ObjectRouter struct {
	proxy  interface{}
	objMap map[string]*ObjectHandler
	log    *core.LogAgent
	wg     sync.WaitGroup
}

func NewObjectRouter() *ObjectRouter {
	or := ObjectRouter{
		proxy:  nil,
		objMap: make(map[string]*ObjectHandler),
		log:    nil,
	}
	return &or
}

func (or *ObjectRouter) InitRouter(log *core.LogAgent, proxy interface{}) {
	or.log = log
	or.proxy = proxy
}

func (or *ObjectRouter) Cleanup() {
	for name, obj := range or.objMap {
		obj.StopObject()
		delete(or.objMap, name)
	}
	or.wg.Wait()
}

func (or *ObjectRouter) GetObjectHandler(name string) *ObjectHandler {
	if name == "" {
		return nil
	}
	obj, ok := or.objMap[name]
	if ok {
		return obj
	}
	obj = NewObjectHandler(name, or.log)
	err := obj.InitObject(or.proxy)
	if err != nil {
		return nil
	}
	obj.StartObject(&or.wg)
	or.objMap[name] = obj
	return obj
}

func (or *ObjectRouter) DelObjectHandler(name string) {
	if name == "" {
		return
	}
	obj, ok := or.objMap[name]
	if ok {
		obj.StopObject()
	}
	delete(or.objMap, name)
}

// Implementation of duplex.ClientManager
func (or *ObjectRouter) OnClientStarted(name string) {
	if name == "" {
		return
	}
	or.log.Trace("ObjectProxy.OnClientStarted dev:%s", name)
	obj := or.GetObjectHandler(name)
	obj.OnClientStarted(name)
}

func (or *ObjectRouter) OnClientStopped(name string) {
	if name == "" {
		return
	}
	or.log.Trace("ObjectProxy.OnClientStopped dev:%s", name)
	obj := or.GetObjectHandler(name)
	obj.OnClientStopped(name)
}

// Implementation of common.SystemCallback
func (or *ObjectRouter) SystemReply(name string, reply *common.SystemReply) error {
	object := or.GetObjectHandler(name)
	if object != nil {
		return object.SystemReply(name, reply)
	}
	return nil
}
func (or *ObjectRouter) SystemHealth(name string, reply *common.SystemHealth) error {
	object := or.GetObjectHandler(name)
	if object != nil {
		return object.SystemHealth(name, reply)
	}
	return nil
}

// Implementation of common.DeviceCallback
func (or *ObjectRouter) DeviceReply(name string, reply *common.DeviceReply) error {
	object := or.GetObjectHandler(name)
	if object != nil {
		return object.DeviceReply(name, reply)
	}
	return nil
}
func (or *ObjectRouter) ExecuteError(name string, reply *common.DeviceError) error {
	object := or.GetObjectHandler(name)
	if object != nil {
		return object.ExecuteError(name, reply)
	}
	return nil
}
func (or *ObjectRouter) StateChanged(name string, reply *common.DeviceState) error {
	object := or.GetObjectHandler(name)
	if object != nil {
		return object.StateChanged(name, reply)
	}
	return nil
}
func (or *ObjectRouter) ActionPrompt(name string, reply *common.DevicePrompt) error {
	object := or.GetObjectHandler(name)
	if object != nil {
		return object.ActionPrompt(name, reply)
	}
	return nil
}
func (or *ObjectRouter) ReaderReturn(name string, reply *common.DeviceInform) error {
	object := or.GetObjectHandler(name)
	if object != nil {
		return object.ReaderReturn(name, reply)
	}
	return nil
}

// Implementation of common.ReaderCallback
func (or *ObjectRouter) CardDescription(name string, reply *common.ReaderCardInfo) error {
	object := or.GetObjectHandler(name)
	if object != nil {
		return object.CardDescription(name, reply)
	}
	return nil
}
func (or *ObjectRouter) ChipResponse(name string, reply *common.ReaderChipReply) error {
	object := or.GetObjectHandler(name)
	if object != nil {
		return object.ChipResponse(name, reply)
	}
	return nil
}
func (or *ObjectRouter) PinPadReply(name string, reply *common.ReaderPinReply) error {
	object := or.GetObjectHandler(name)
	if object != nil {
		return object.PinPadReply(name, reply)
	}
	return nil
}

// Implementation of common.ValidatorCallback
func (or *ObjectRouter) NoteAccepted(name string, reply *common.ValidatorAccept) error {
	object := or.GetObjectHandler(name)
	if object != nil {
		return object.NoteAccepted(name, reply)
	}
	return nil
}
func (or *ObjectRouter) CashIsStored(name string, reply *common.ValidatorAccept) error {
	object := or.GetObjectHandler(name)
	if object != nil {
		return object.CashIsStored(name, reply)
	}
	return nil
}
func (or *ObjectRouter) CashReturned(name string, reply *common.ValidatorAccept) error {
	object := or.GetObjectHandler(name)
	if object != nil {
		return object.CashReturned(name, reply)
	}
	return nil
}
func (or *ObjectRouter) ValidatorStore(name string, reply *common.ValidatorStore) error {
	object := or.GetObjectHandler(name)
	if object != nil {
		return object.ValidatorStore(name, reply)
	}
	return nil
}
