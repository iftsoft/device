package handler

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
)

type DeviceTesterFactory struct {
}

func (dtf DeviceTesterFactory) GetPluginName() string {
	return "DeviceTesterBehavior"
}

func (dtf DeviceTesterFactory) CreatePlugin(devName string, proxy interface{}, log *core.LogAgent) (error, PluginManager) {
	dt := &DeviceTester{
		devName:      devName,
		enabled:      false,
		connected:    false,
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
	err := dt.InitPlugin(proxy)
	return err, dt
}

type TestFunc func(*DeviceTester) error

type DeviceTester struct {
	devName      string
	enabled      bool
	connected    bool
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

func (dt *DeviceTester) InitPlugin(proxy interface{}) error {
	dt.log.Debug("DeviceTester run InitPlugin for:%s", dt.devName)
	if system, ok := proxy.(common.SystemManager); ok {
		dt.systemMng = system
	}
	if device, ok := proxy.(common.DeviceManager); ok {
		dt.deviceMng = device
	}
	if printer, ok := proxy.(common.PrinterManager); ok {
		dt.printerMng = printer
	}
	if reader, ok := proxy.(common.ReaderManager); ok {
		dt.readerMng = reader
	}
	if valid, ok := proxy.(common.ValidatorManager); ok {
		dt.validatorMng = valid
	}
	if pinpad, ok := proxy.(common.PinPadManager); ok {
		dt.pinpadMng = pinpad
	}
	return nil
}


func (dt *DeviceTester) Enabled(on bool) {
	dt.enabled = on
}

func (dt *DeviceTester) Connected(on bool){
	dt.connected = on
}

func (dt *DeviceTester) OnTimerTick(){
	dt.log.Trace("Object handler %s onTimerTick", dt.devName)
	if len(dt.tests) > 0 {
		tstFunc := dt.tests[0]
		dt.tests = dt.tests[1:]
		if tstFunc != nil {
			_ = tstFunc(dt)
		}
	}
}



// Implementation of duplex.ClientManager
func (oh *DeviceTester) OnClientStarted(name string) {
	oh.log.Debug("DeviceTester.OnClientStarted dev:%s", name)
	oh.fillTestList()
}

func (oh *DeviceTester) OnClientStopped(name string) {
	oh.log.Debug("DeviceTester.OnClientStopped dev:%s", name)
	oh.tests = nil
}

func (oh *DeviceTester) fillTestList() {
	oh.tests = append(oh.tests, func(hnd *DeviceTester) error {
		return hnd.systemMng.SysInform(hnd.devName, &common.SystemQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *DeviceTester) error {
		return hnd.systemMng.SysStart(hnd.devName, &common.SystemQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *DeviceTester) error {
		return hnd.deviceMng.Reset(hnd.devName, &common.DeviceQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *DeviceTester) error {
		return hnd.deviceMng.Status(hnd.devName, &common.DeviceQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *DeviceTester) error {
		return hnd.deviceMng.Status(hnd.devName, &common.DeviceQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *DeviceTester) error {
		return hnd.deviceMng.RunAction(hnd.devName, &common.DeviceQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *DeviceTester) error {
		return hnd.deviceMng.Status(hnd.devName, &common.DeviceQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *DeviceTester) error {
		return hnd.deviceMng.Cancel(hnd.devName, &common.DeviceQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *DeviceTester) error {
		return hnd.deviceMng.Status(hnd.devName, &common.DeviceQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *DeviceTester) error {
		return hnd.deviceMng.StopAction(hnd.devName, &common.DeviceQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *DeviceTester) error {
		return hnd.deviceMng.Status(hnd.devName, &common.DeviceQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *DeviceTester) error {
		return hnd.systemMng.SysStop(hnd.devName, &common.SystemQuery{})
	})
	oh.tests = append(oh.tests, func(hnd *DeviceTester) error {
		return hnd.systemMng.Terminate(hnd.devName, &common.SystemQuery{})
	})
}

// Implementation of common.SystemCallback
func (oh *DeviceTester) SystemReply(name string, reply *common.SystemReply) error {
	if oh.log != nil {
		oh.log.Debug("DeviceTester.SystemReply dev:%s get cmd:%s",
			name, reply.Command)
	}
	return nil
}

func (oh *DeviceTester) SystemHealth(name string, reply *common.SystemHealth) error {
	if oh.log != nil {
		oh.log.Debug("DeviceTester.SystemHealth dev:%s for moment:%d",
			name, reply.Moment)
	}
	return nil
}

// Implementation of common.DeviceCallback
func (oh *DeviceTester) DeviceReply(name string, reply *common.DeviceReply) error {
	if oh.log != nil {
		oh.log.Debug("DeviceTester.DeviceReply dev:%s, cmd:%s, state:%s, error:%d - %s",
			name, reply.Command, reply.DevState, reply.ErrCode, reply.ErrText)
	}
	return nil
}
func (oh *DeviceTester) ExecuteError(name string, reply *common.DeviceError) error {
	if oh.log != nil {
		oh.log.Debug("DeviceTester.ExecuteError dev:%s, action:%s, error:%d - %s",
			name, reply.DevState, reply.ErrCode, reply.ErrText)
	}
	return nil
}
func (oh *DeviceTester) StateChanged(name string, reply *common.DeviceState) error {
	if oh.log != nil {
		oh.log.Debug("DeviceTester.StateChanged dev:%s, old state:%s, new state:%s",
			name, reply.OldState, reply.NewState)
	}
	return nil
}
func (oh *DeviceTester) ActionPrompt(name string, reply *common.DevicePrompt) error {
	if oh.log != nil {
		oh.log.Debug("DeviceTester.ActionPrompt dev:%s, action:%s, prompt:%s",
			name, reply.Action, reply.Prompt)
	}
	return nil
}
func (oh *DeviceTester) ReaderReturn(name string, reply *common.DeviceInform) error {
	if oh.log != nil {
		oh.log.Debug("DeviceTester.ReaderReturn dev:%s, action:%s, info:%s",
			name, reply.Action, reply.Inform)
	}
	return nil
}

// Implementation of common.PrinterCallback
func (oh *DeviceTester) PrinterProgress(name string, reply *common.PrinterProgress) error {
	if oh.log != nil {
		oh.log.Debug("DeviceTester.PrinterProgress dev:%s, Done:%d, From:%d",
			name, reply.PageDone, reply.PagesAll)
	}
	return nil
}

// Implementation of common.ReaderCallback
func (oh *DeviceTester) CardPosition(name string, reply *common.ReaderCardPos) error {
	if oh.log != nil {
		oh.log.Debug("DeviceTester.CardPosition dev:%s, Position:%d",
			name, reply.Position)
	}
	return nil
}
func (oh *DeviceTester) CardDescription(name string, reply *common.ReaderCardInfo) error {
	if oh.log != nil {
		oh.log.Debug("DeviceTester.CardDescription dev:%s, CardPAN:%s, ExpDate:%s",
			name, reply.CardPan, reply.ExpDate)
	}
	return nil
}
func (oh *DeviceTester) ChipResponse(name string, reply *common.ReaderChipReply) error {
	if oh.log != nil {
		oh.log.Debug("DeviceTester.ChipResponse dev:%s, Protocol:%d",
			name, reply.Protocol)
	}
	return nil
}

// Implementation of common.ValidatorCallback
func (oh *DeviceTester) NoteAccepted(name string, reply *common.ValidatorAccept) error {
	if oh.log != nil {
		oh.log.Debug("DeviceTester.NoteAccepted dev:%s, Reply: %s",
			name, reply.String())
	}
	return nil
}
func (oh *DeviceTester) CashIsStored(name string, reply *common.ValidatorAccept) error {
	if oh.log != nil {
		oh.log.Debug("DeviceTester.CashIsStored dev:%s, Reply: %s",
			name, reply.String())
	}
	return nil
}
func (oh *DeviceTester) CashReturned(name string, reply *common.ValidatorAccept) error {
	if oh.log != nil {
		oh.log.Debug("DeviceTester.CashReturned dev:%s, Reply: %s",
			name, reply.String())
	}
	return nil
}
func (oh *DeviceTester) ValidatorStore(name string, reply *common.ValidatorStore) error {
	if oh.log != nil {
		oh.log.Debug("DeviceTester.ValidatorStore dev:%s, Reply: %s",
			name, reply.String())
	}
	return nil
}

// Implementation of common.PinPadCallback
func (oh *DeviceTester) PinPadReply(name string, reply *common.ReaderPinReply) error {
	if oh.log != nil {
		oh.log.Debug("DeviceTester.PinPadReply dev:%s, PinLen:%d",
			name, reply.PinLength)
	}
	return nil
}

