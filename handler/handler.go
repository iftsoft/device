package handler

import (
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"sync"
	"time"
)


type DeviceHandler struct {
	devName      string
	pluginMap    map[string]PluginManager
	systemCbk    []common.SystemCallback
	deviceCbk    []common.DeviceCallback
	printerCbk   []common.PrinterCallback
	readerCbk    []common.ReaderCallback
	validatorCbk []common.ValidatorCallback
	pinpadCbk    []common.PinPadCallback
	log          *core.LogAgent
	done         chan struct{}
}

func NewDeviceHandler(name string, log *core.LogAgent) *DeviceHandler {
	oh := DeviceHandler{
		devName:      name,
		pluginMap:    make(map[string]PluginManager),
		systemCbk:    make([]common.SystemCallback, 0),
		deviceCbk:    make([]common.DeviceCallback, 0),
		printerCbk:   make([]common.PrinterCallback, 0),
		readerCbk:    make([]common.ReaderCallback, 0),
		validatorCbk: make([]common.ValidatorCallback, 0),
		pinpadCbk:    make([]common.PinPadCallback, 0),
		log:          log,
		done:         make(chan struct{}),
	}
	return &oh
}

func (oh *DeviceHandler) AttachPlugin(name string, plugin interface{}) error {
	oh.log.Trace("DeviceHandler run AttachPlugin for %s", name)
	_, ok := oh.pluginMap[name]
	if ok { return nil 	}
	if manager, ok := plugin.(PluginManager); ok {
		oh.pluginMap[name] = manager
	} else {
		return errors.New("plugin is corrupted")
	}
	if system, ok := plugin.(common.SystemCallback); ok {
		oh.systemCbk = append(oh.systemCbk, system)
	}
	if device, ok := plugin.(common.DeviceCallback); ok {
		oh.deviceCbk = append(oh.deviceCbk, device)
	}
	if printer, ok := plugin.(common.PrinterCallback); ok {
		oh.printerCbk = append(oh.printerCbk, printer)
	}
	if reader, ok := plugin.(common.ReaderCallback); ok {
		oh.readerCbk = append(oh.readerCbk, reader)
	}
	if valid, ok := plugin.(common.ValidatorCallback); ok {
		oh.validatorCbk = append(oh.validatorCbk, valid)
	}
	if pinpad, ok := plugin.(common.PinPadCallback); ok {
		oh.pinpadCbk = append(oh.pinpadCbk, pinpad)
	}
	return nil
}

func (oh *DeviceHandler) panicRecover() {
	if r := recover(); r != nil {
		if oh.log != nil {
			oh.log.Panic("Panic is recovered: %+v", r)
		}
	}
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
	oh.log.Debug("Device handler loop for dev:%s is started", oh.devName)
	defer oh.log.Debug("Device handler loop for dev:%s is stopped", oh.devName)

	tick := time.NewTicker(1000 * time.Millisecond)
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
	oh.log.Trace("Device handler %s onTimerTick %s", oh.devName, tm.Format(time.StampMilli))
	for _, pl := range oh.pluginMap {
		go func() {
			defer oh.panicRecover()
			pl.OnTimerTick()
		}()
	}
}

// Implementation of duplex.ClientManager
func (oh *DeviceHandler) OnClientStarted(name string) {
	oh.log.Debug("DeviceHandler.OnClientStarted dev:%s", name)
	for _, pl := range oh.pluginMap {
		go func() {
			defer oh.panicRecover()
			pl.Connected(true)
		}()
	}
}

func (oh *DeviceHandler) OnClientStopped(name string) {
	oh.log.Debug("DeviceHandler.OnClientStopped dev:%s", name)
	for _, pl := range oh.pluginMap {
		go func() {
			defer oh.panicRecover()
			pl.Connected(false)
		}()
	}
}



// Implementation of common.SystemCallback
func (oh *DeviceHandler) SystemReply(name string, reply *common.SystemReply) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.SystemReply dev:%s get cmd:%s",
			name, reply.Command)
	}
	for _, callback := range oh.systemCbk {
		go func() {
			defer oh.panicRecover()
			_ = callback.SystemReply(name, reply)
		}()
	}
	return nil
}

func (oh *DeviceHandler) SystemHealth(name string, reply *common.SystemHealth) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.SystemHealth dev:%s for moment:%d",
			name, reply.Moment)
	}
	for _, callback := range oh.systemCbk {
		go func() {
			defer oh.panicRecover()
			_ = callback.SystemHealth(name, reply)
		}()
	}
	return nil
}

// Implementation of common.DeviceCallback
func (oh *DeviceHandler) DeviceReply(name string, reply *common.DeviceReply) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.DeviceReply dev:%s, cmd:%s, state:%s, error:%d - %s",
			name, reply.Command, reply.DevState, reply.ErrCode, reply.ErrText)
	}
	for _, callback := range oh.deviceCbk {
		go func() {
			defer oh.panicRecover()
			_ = callback.DeviceReply(name, reply)
		}()
	}
	return nil
}
func (oh *DeviceHandler) ExecuteError(name string, reply *common.DeviceError) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.ExecuteError dev:%s, action:%s, error:%d - %s",
			name, reply.DevState, reply.ErrCode, reply.ErrText)
	}
	for _, callback := range oh.deviceCbk {
		go func() {
			defer oh.panicRecover()
			_ = callback.ExecuteError(name, reply)
		}()
	}
	return nil
}
func (oh *DeviceHandler) StateChanged(name string, reply *common.DeviceState) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.StateChanged dev:%s, old state:%s, new state:%s",
			name, reply.OldState, reply.NewState)
	}
	for _, callback := range oh.deviceCbk {
		go func() {
			defer oh.panicRecover()
			_ = callback.StateChanged(name, reply)
		}()
	}
	return nil
}
func (oh *DeviceHandler) ActionPrompt(name string, reply *common.DevicePrompt) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.ActionPrompt dev:%s, action:%s, prompt:%s",
			name, reply.Action, reply.Prompt)
	}
	for _, callback := range oh.deviceCbk {
		go func() {
			defer oh.panicRecover()
			_ = callback.ActionPrompt(name, reply)
		}()
	}
	return nil
}
func (oh *DeviceHandler) ReaderReturn(name string, reply *common.DeviceInform) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.ReaderReturn dev:%s, action:%s, info:%s",
			name, reply.Action, reply.Inform)
	}
	for _, callback := range oh.deviceCbk {
		go func() {
			defer oh.panicRecover()
			_ = callback.ReaderReturn(name, reply)
		}()
	}
	return nil
}

// Implementation of common.PrinterCallback
func (oh *DeviceHandler) PrinterProgress(name string, reply *common.PrinterProgress) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.PrinterProgress dev:%s, Done:%d, From:%d",
			name, reply.PageDone, reply.PagesAll)
	}
	for _, callback := range oh.printerCbk {
		go func() {
			defer oh.panicRecover()
			_ = callback.PrinterProgress(name, reply)
		}()
	}
	return nil
}

// Implementation of common.ReaderCallback
func (oh *DeviceHandler) CardPosition(name string, reply *common.ReaderCardPos) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.CardPosition dev:%s, Position:%d",
			name, reply.Position)
	}
	for _, callback := range oh.readerCbk {
		go func() {
			defer oh.panicRecover()
			_ = callback.CardPosition(name, reply)
		}()
	}
	return nil
}
func (oh *DeviceHandler) CardDescription(name string, reply *common.ReaderCardInfo) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.CardDescription dev:%s, CardPAN:%s, ExpDate:%s",
			name, reply.CardPan, reply.ExpDate)
	}
	for _, callback := range oh.readerCbk {
		go func() {
			defer oh.panicRecover()
			_ = callback.CardDescription(name, reply)
		}()
	}
	return nil
}
func (oh *DeviceHandler) ChipResponse(name string, reply *common.ReaderChipReply) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.ChipResponse dev:%s, Protocol:%d",
			name, reply.Protocol)
	}
	for _, callback := range oh.readerCbk {
		go func() {
			defer oh.panicRecover()
			_ = callback.ChipResponse(name, reply)
		}()
	}
	return nil
}

// Implementation of common.ValidatorCallback
func (oh *DeviceHandler) NoteAccepted(name string, reply *common.ValidatorAccept) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.NoteAccepted dev:%s, Reply: %s",
			name, reply.String())
	}
	for _, callback := range oh.validatorCbk {
		go func() {
			defer oh.panicRecover()
			_ = callback.NoteAccepted(name, reply)
		}()
	}
	return nil
}
func (oh *DeviceHandler) CashIsStored(name string, reply *common.ValidatorAccept) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.CashIsStored dev:%s, Reply: %s",
			name, reply.String())
	}
	for _, callback := range oh.validatorCbk {
		go func() {
			defer oh.panicRecover()
			_ = callback.CashIsStored(name, reply)
		}()
	}
	return nil
}
func (oh *DeviceHandler) CashReturned(name string, reply *common.ValidatorAccept) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.CashReturned dev:%s, Reply: %s",
			name, reply.String())
	}
	for _, callback := range oh.validatorCbk {
		go func() {
			defer oh.panicRecover()
			_ = callback.CashReturned(name, reply)
		}()
	}
	return nil
}
func (oh *DeviceHandler) ValidatorStore(name string, reply *common.ValidatorStore) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.ValidatorStore dev:%s, Reply: %s",
			name, reply.String())
	}
	for _, callback := range oh.validatorCbk {
		go func() {
			defer oh.panicRecover()
			_ = callback.ValidatorStore(name, reply)
		}()
	}
	return nil
}

// Implementation of common.PinPadCallback
func (oh *DeviceHandler) PinPadReply(name string, reply *common.ReaderPinReply) error {
	if oh.log != nil {
		oh.log.Debug("DeviceHandler.PinPadReply dev:%s, PinLen:%d",
			name, reply.PinLength)
	}
	for _, callback := range oh.pinpadCbk {
		go func() {
			defer oh.panicRecover()
			_ = callback.PinPadReply(name, reply)
		}()
	}
	return nil
}

