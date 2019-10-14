package handler

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"sync"
	"time"
)

type TestFunc func(*DeviceHandler) error

type DeviceHandler struct {
	devName      string
	systemMng    common.SystemManager
	deviceMng    common.DeviceManager
	printerMng   common.PrinterManager
	readerMng    common.ReaderManager
	validatorMng common.ValidatorManager
	pinpadMng    common.PinPadManager
	log          *core.LogAgent
	done         chan struct{}
	tests        []TestFunc
}

func NewDeviceHandler(name string, log *core.LogAgent) *DeviceHandler {
	oh := DeviceHandler{
		devName:      name,
		systemMng:    nil,
		deviceMng:    nil,
		printerMng:   nil,
		readerMng:    nil,
		validatorMng: nil,
		pinpadMng:    nil,
		log:          log,
		done:         make(chan struct{}),
		tests:        make([]TestFunc, 0),
	}
	return &oh
}

func (oh *DeviceHandler) InitObject(proxy interface{}) error {
	oh.log.Debug("DeviceHandler run cmd:%s", "InitObject")
	if system, ok := proxy.(common.SystemManager); ok {
		oh.systemMng = system
	}
	if device, ok := proxy.(common.DeviceManager); ok {
		oh.deviceMng = device
	}
	if printer, ok := proxy.(common.PrinterManager); ok {
		oh.printerMng = printer
	}
	if reader, ok := proxy.(common.ReaderManager); ok {
		oh.readerMng = reader
	}
	if valid, ok := proxy.(common.ValidatorManager); ok {
		oh.validatorMng = valid
	}
	if pinpad, ok := proxy.(common.PinPadManager); ok {
		oh.pinpadMng = pinpad
	}
	return nil
}

func (oh *DeviceHandler) StartObject(wg *sync.WaitGroup) {
	oh.log.Info("Starting object handle")
	go oh.objectHandlerLoop(wg)
}

func (oh *DeviceHandler) StopObject() {
	oh.log.Info("Stopping object handle")
	close(oh.done)
}

func (oh *DeviceHandler) objectHandlerLoop(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	oh.log.Debug("Object handler loop for dev:%s is started", oh.devName)
	defer oh.log.Debug("Object handler loop for dev:%s is stopped", oh.devName)

	tick := time.NewTicker(5 * time.Second)
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

func (oh *DeviceHandler) onTimerTick(tm time.Time) {
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
func (oh *DeviceHandler) OnClientStarted(name string) {
	oh.log.Debug("DeviceHandler.OnClientStarted dev:%s", name)
	oh.fillTestList()
}

func (oh *DeviceHandler) OnClientStopped(name string) {
	oh.log.Debug("DeviceHandler.OnClientStopped dev:%s", name)
	oh.tests = nil
}

func (oh *DeviceHandler) fillTestList() {
	oh.tests = append(oh.tests, func(hnd *DeviceHandler) error {
		return hnd.systemMng.SysInform(hnd.devName, &common.SystemQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *DeviceHandler) error {
		return hnd.systemMng.SysStart(hnd.devName, &common.SystemQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *DeviceHandler) error {
		return hnd.deviceMng.Reset(hnd.devName, &common.DeviceQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *DeviceHandler) error {
		return hnd.deviceMng.Status(hnd.devName, &common.DeviceQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *DeviceHandler) error {
		return hnd.deviceMng.Status(hnd.devName, &common.DeviceQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *DeviceHandler) error {
		return hnd.deviceMng.RunAction(hnd.devName, &common.DeviceQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *DeviceHandler) error {
		return hnd.deviceMng.Status(hnd.devName, &common.DeviceQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *DeviceHandler) error {
		return hnd.deviceMng.Cancel(hnd.devName, &common.DeviceQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *DeviceHandler) error {
		return hnd.deviceMng.Status(hnd.devName, &common.DeviceQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *DeviceHandler) error {
		return hnd.deviceMng.StopAction(hnd.devName, &common.DeviceQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *DeviceHandler) error {
		return hnd.deviceMng.Status(hnd.devName, &common.DeviceQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *DeviceHandler) error {
		return hnd.systemMng.SysStop(hnd.devName, &common.SystemQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *DeviceHandler) error {
		return hnd.systemMng.Terminate(hnd.devName, &common.SystemQuery{})
	})
}

// Implementation of common.SystemCallback
func (oh *DeviceHandler) SystemReply(name string, reply *common.SystemReply) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.SystemReply dev:%s get cmd:%s",
			name, reply.Command)
	}
	return nil
}

func (oh *DeviceHandler) SystemHealth(name string, reply *common.SystemHealth) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.SystemHealth dev:%s for moment:%d",
			name, reply.Moment)
	}
	return nil
}

// Implementation of common.DeviceCallback
func (oh *DeviceHandler) DeviceReply(name string, reply *common.DeviceReply) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.DeviceReply dev:%s, cmd:%s, state:%s, error:%d - %s",
			name, reply.Command, reply.DevState, reply.ErrCode, reply.ErrText)
	}
	return nil
}
func (oh *DeviceHandler) ExecuteError(name string, reply *common.DeviceError) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.ExecuteError dev:%s, action:%s, error:%d - %s",
			name, reply.DevState, reply.ErrCode, reply.ErrText)
	}
	return nil
}
func (oh *DeviceHandler) StateChanged(name string, reply *common.DeviceState) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.StateChanged dev:%s, old state:%s, new state:%s",
			name, reply.OldState, reply.NewState)
	}
	return nil
}
func (oh *DeviceHandler) ActionPrompt(name string, reply *common.DevicePrompt) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.ActionPrompt dev:%s, action:%s, prompt:%s",
			name, reply.Action, reply.Prompt)
	}
	return nil
}
func (oh *DeviceHandler) ReaderReturn(name string, reply *common.DeviceInform) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.ReaderReturn dev:%s, action:%s, info:%s",
			name, reply.Action, reply.Inform)
	}
	return nil
}

// Implementation of common.PrinterCallback
func (oh *DeviceHandler) PrinterProgress(name string, reply *common.PrinterProgress) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.PrinterProgress dev:%s, Done:%d, From:%d",
			name, reply.PageDone, reply.PagesAll)
	}
	return nil
}

// Implementation of common.ReaderCallback
func (oh *DeviceHandler) CardPosition(name string, reply *common.ReaderCardPos) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.CardPosition dev:%s, Position:%d",
			name, reply.Position)
	}
	return nil
}
func (oh *DeviceHandler) CardDescription(name string, reply *common.ReaderCardInfo) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.CardDescription dev:%s, CardPAN:%s, ExpDate:%s",
			name, reply.CardPan, reply.ExpDate)
	}
	return nil
}
func (oh *DeviceHandler) ChipResponse(name string, reply *common.ReaderChipReply) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.ChipResponse dev:%s, Protocol:%d",
			name, reply.Protocol)
	}
	return nil
}

// Implementation of common.ValidatorCallback
func (oh *DeviceHandler) NoteAccepted(name string, reply *common.ValidatorAccept) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.NoteAccepted dev:%s, Reply: %s",
			name, reply.String())
	}
	return nil
}
func (oh *DeviceHandler) CashIsStored(name string, reply *common.ValidatorAccept) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.CashIsStored dev:%s, Reply: %s",
			name, reply.String())
	}
	return nil
}
func (oh *DeviceHandler) CashReturned(name string, reply *common.ValidatorAccept) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.CashReturned dev:%s, Reply: %s",
			name, reply.String())
	}
	return nil
}
func (oh *DeviceHandler) ValidatorStore(name string, reply *common.ValidatorStore) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.ValidatorStore dev:%s, Reply: %s",
			name, reply.String())
	}
	return nil
}

// Implementation of common.PinPadCallback
func (oh *DeviceHandler) PinPadReply(name string, reply *common.ReaderPinReply) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.PinPadReply dev:%s, PinLen:%d",
			name, reply.PinLength)
	}
	return nil
}
