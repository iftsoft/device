package handler

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"sync"
	"time"
)

type TestFunc func(*ObjectHandler) error

type ObjectHandler struct {
	devName string
	system  common.SystemManager
	device  common.DeviceManager
	log     *core.LogAgent
	done    chan struct{}
	tests   []TestFunc
}

func NewObjectState(name string, log *core.LogAgent) *ObjectHandler {
	oh := ObjectHandler{
		devName: name,
		system:  nil,
		device:  nil,
		log:     log,
		done:    make(chan struct{}),
		tests:   make([]TestFunc, 0),
	}
	return &oh
}

func (oh *ObjectHandler) InitObject(proxy interface{}) error {
	oh.log.Debug("ObjectHandler run cmd:%s", "InitObject")
	if system, ok := proxy.(common.SystemManager); ok {
		oh.system = system
	}
	if device, ok := proxy.(common.DeviceManager); ok {
		oh.device = device
	}
	return nil
}

func (oh *ObjectHandler) StartObject(wg *sync.WaitGroup) {
	oh.log.Info("Starting object handle")
	go oh.objectHandlerLoop(wg)
}

func (oh *ObjectHandler) StopObject() {
	oh.log.Info("Stopping object handle")
	close(oh.done)
}

func (oh *ObjectHandler) objectHandlerLoop(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	oh.log.Debug("Object handler loop for dev:%s is started", oh.devName)
	defer oh.log.Debug("Object handler loop for dev:%s is stopped", oh.devName)

	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()

	for {
		select {
		case <-oh.done:
			return
		case tm := <-tick.C:
			oh.onTimerTick(tm)
		}
	}
}

func (oh *ObjectHandler) onTimerTick(tm time.Time) {
	oh.log.Trace("Object handler %s onTimerTick %s", oh.devName, tm.Format(time.StampMilli))
	if len(oh.tests) > 0 {
		tstFunc := oh.tests[0]
		oh.tests = oh.tests[1:]
		if tstFunc != nil {
			_ = tstFunc(oh)
		}
	}
}

// Implementation of duplex.ClientManager
func (oh *ObjectHandler) OnClientStarted(name string) {
	oh.log.Debug("ObjectHandler.OnClientStarted dev:%s", name)
	oh.fillTestList()
}

func (oh *ObjectHandler) OnClientStopped(name string) {
	oh.log.Debug("ObjectHandler.OnClientStopped dev:%s", name)
	oh.tests = nil
}

func (oh *ObjectHandler) fillTestList() {
	oh.tests = append(oh.tests, func(hnd *ObjectHandler) error {
		return hnd.system.Config(hnd.devName, &common.SystemQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *ObjectHandler) error {
		return hnd.system.Start(hnd.devName, &common.SystemQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *ObjectHandler) error {
		return hnd.system.Inform(hnd.devName, &common.SystemQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *ObjectHandler) error {
		return hnd.device.Reset(hnd.devName, &common.DeviceQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *ObjectHandler) error {
		return hnd.device.Status(hnd.devName, &common.DeviceQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *ObjectHandler) error {
		return hnd.device.Cancel(hnd.devName, &common.DeviceQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *ObjectHandler) error {
		return hnd.system.Stop(hnd.devName, &common.SystemQuery{})
	})
}

// Implementation of common.SystemCallback
func (oh *ObjectHandler) SystemReply(name string, reply *common.SystemReply) error {
	if oh.log != nil {
		oh.log.Debug("ObjectHandler.SystemReply dev:%s get cmd:%s",
			name, reply.Command)
	}
	return nil
}

func (oh *ObjectHandler) SystemHealth(name string, reply *common.SystemHealth) error {
	if oh.log != nil {
		oh.log.Debug("ObjectHandler.SystemHealth dev:%s for moment:%d",
			name, reply.Moment)
	}
	return nil
}

// Implementation of common.DeviceCallback
func (oh *ObjectHandler) DeviceReply(name string, reply *common.DeviceReply) error {
	if oh.log != nil {
		oh.log.Debug("ObjectHandler.DeviceReply dev:%s, cmd:%s, state:%d",
			name, reply.Command, reply.DevState)
	}
	return nil
}

func (oh *ObjectHandler) ExecuteError(name string, reply *common.DeviceError) error {
	if oh.log != nil {
		oh.log.Debug("ObjectHandler.ExecuteError dev:%s, action:%d, error:%d - %s",
			name, reply.Action, reply.ErrCode, reply.ErrText)
	}
	return nil
}

func (oh *ObjectHandler) StateChanged(name string, reply *common.DeviceState) error {
	if oh.log != nil {
		oh.log.Debug("ObjectHandler.StateChanged dev:%s, old state:%d, new state:%d",
			name, reply.OldState, reply.NewState)
	}
	return nil
}

func (oh *ObjectHandler) ActionPrompt(name string, reply *common.DevicePrompt) error {
	if oh.log != nil {
		oh.log.Debug("ObjectHandler.ActionPrompt dev:%s, action:%d, prompt:%d",
			name, reply.Action, reply.Prompt)
	}
	return nil
}