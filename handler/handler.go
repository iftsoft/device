package handler

import (
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"sync"
	"time"
)


type DeviceHandler struct {
	devName      string
	config       *config.HandlerConfig
	proxy        *HandlerProxy
	reflexMap    map[string]ReflexManager
	systemCbk    []common.SystemCallback
	deviceCbk    []common.DeviceCallback
	printerCbk   []common.PrinterCallback
	readerCbk    []common.ReaderCallback
	validatorCbk []common.ValidatorCallback
	pinpadCbk    []common.PinPadCallback
	isRunning    bool
	log          *core.LogAgent
	done         chan struct{}
}

func NewDeviceHandler(name string, cfg *config.HandlerConfig, log *core.LogAgent) *DeviceHandler {
	dh := DeviceHandler{
		devName:      name,
		config:       cfg,
		reflexMap:    make(map[string]ReflexManager),
		systemCbk:    make([]common.SystemCallback, 0),
		deviceCbk:    make([]common.DeviceCallback, 0),
		printerCbk:   make([]common.PrinterCallback, 0),
		readerCbk:    make([]common.ReaderCallback, 0),
		validatorCbk: make([]common.ValidatorCallback, 0),
		pinpadCbk:    make([]common.PinPadCallback, 0),
		isRunning:    false,
		log:          log,
		done:         make(chan struct{}),
	}
	return &dh
}

func (dh *DeviceHandler) AttachProxy(proxy *HandlerProxy) {
	dh.proxy = proxy
}

func (dh *DeviceHandler) AttachReflex(name string, reflex interface{}) error {
	dh.log.Trace("DeviceHandler run AttachReflex for reflex:%s to device:%s", name, dh.devName)
	_, ok := dh.reflexMap[name]
	if ok { return nil 	}
	if manager, ok := reflex.(ReflexManager); ok {
		dh.reflexMap[name] = manager
	} else {
		return errors.New("reflex is corrupted")
	}
	if system, ok := reflex.(common.SystemCallback); ok {
		dh.systemCbk = append(dh.systemCbk, system)
	}
	if device, ok := reflex.(common.DeviceCallback); ok {
		dh.deviceCbk = append(dh.deviceCbk, device)
	}
	if printer, ok := reflex.(common.PrinterCallback); ok {
		dh.printerCbk = append(dh.printerCbk, printer)
	}
	if reader, ok := reflex.(common.ReaderCallback); ok {
		dh.readerCbk = append(dh.readerCbk, reader)
	}
	if valid, ok := reflex.(common.ValidatorCallback); ok {
		dh.validatorCbk = append(dh.validatorCbk, valid)
	}
	if pinpad, ok := reflex.(common.PinPadCallback); ok {
		dh.pinpadCbk = append(dh.pinpadCbk, pinpad)
	}
	return nil
}

func (dh *DeviceHandler) panicRecover() {
	if r := recover(); r != nil {
		if dh.log != nil {
			dh.log.Panic("Panic is recovered: %+v", r)
		}
	}
}

func (dh *DeviceHandler) StartObject(wg *sync.WaitGroup) {
	dh.log.Info("Starting device handle")
	go dh.objectHandlerLoop(wg)
}

func (dh *DeviceHandler) StopObject() {
	dh.log.Info("Stopping device handle")
	close(dh.done)
}

func (dh *DeviceHandler) objectHandlerLoop(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	dh.log.Debug("Device handler loop for dev:%s is started", dh.devName)
	defer dh.log.Debug("Device handler loop for dev:%s is stopped", dh.devName)

	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()

	for {
		select {
		case <-dh.done:
			return
		case tm := <-tick.C:
			dh.onTimerTick(tm)
		}
	}
}

func (dh *DeviceHandler) onTimerTick(tm time.Time) {
	if dh.isRunning {
		dh.log.Trace("Device handler %s onTimerTick %s", dh.devName, tm.Format(time.StampMilli))
		for _, item := range dh.reflexMap {
			go func(reflex ReflexManager) {
				defer dh.panicRecover()
				reflex.OnTimerTick()
			}(item)
		}
	}
}

// Implementation of duplex.ClientManager
func (dh *DeviceHandler) OnClientStarted(name string) {
	dh.log.Debug("DeviceHandler.OnClientStarted dev:%s", name)
	for _, item := range dh.reflexMap {
		go func(reflex ReflexManager) {
			defer dh.panicRecover()
			reflex.Connected(true)
		}(item)
	}
	if dh.proxy != nil {
		query := &common.SystemConfig{}
		if dh.config != nil {
			query = dh.config.Config.SystemConfig()
		}
		_ = dh.proxy.SysStart(name, query)
	}
	dh.isRunning = true
}

func (dh *DeviceHandler) OnClientStopped(name string) {
	dh.isRunning = false
	dh.log.Debug("DeviceHandler.OnClientStopped dev:%s", name)
	for _, item := range dh.reflexMap {
		go func(reflex ReflexManager) {
			defer dh.panicRecover()
			reflex.Connected(false)
		}(item)
	}
}



// Implementation of common.SystemCallback
func (dh *DeviceHandler) SystemReply(name string, reply *common.SystemReply) error {
	if dh.log != nil {
		dh.log.Debug("DeviceHandler.SystemReply dev:%s get cmd:%s",
			name, reply.Command)
	}
	for _, cb := range dh.systemCbk {
		go func(callback common.SystemCallback) {
			defer dh.panicRecover()
			_ = callback.SystemReply(name, reply)
		}(cb)
	}
	return nil
}

func (dh *DeviceHandler) SystemHealth(name string, reply *common.SystemHealth) error {
	if dh.log != nil {
		dh.log.Debug("DeviceHandler.SystemHealth dev:%s for moment:%d",
			name, reply.Moment)
	}
	for _, cb := range dh.systemCbk {
		go func(callback common.SystemCallback) {
			defer dh.panicRecover()
			_ = callback.SystemHealth(name, reply)
		}(cb)
	}
	return nil
}

// Implementation of common.DeviceCallback
func (dh *DeviceHandler) DeviceReply(name string, reply *common.DeviceReply) error {
	if dh.log != nil {
		dh.log.Debug("DeviceHandler.DeviceReply dev:%s, cmd:%s, state:%s, error:%d - %s",
			name, reply.Command, reply.DevState, reply.ErrCode, reply.ErrText)
	}
	for _, cb := range dh.deviceCbk {
		go func(callback common.DeviceCallback) {
			defer dh.panicRecover()
			_ = callback.DeviceReply(name, reply)
		}(cb)
	}
	return nil
}
func (dh *DeviceHandler) ExecuteError(name string, reply *common.DeviceError) error {
	if dh.log != nil {
		dh.log.Debug("DeviceHandler.ExecuteError dev:%s, action:%s, error:%d - %s",
			name, reply.DevState, reply.ErrCode, reply.ErrText)
	}
	for _, cb := range dh.deviceCbk {
		go func(callback common.DeviceCallback) {
			defer dh.panicRecover()
			_ = callback.ExecuteError(name, reply)
		}(cb)
	}
	return nil
}
func (dh *DeviceHandler) StateChanged(name string, reply *common.DeviceState) error {
	if dh.log != nil {
		dh.log.Debug("DeviceHandler.StateChanged dev:%s, old state:%s, new state:%s",
			name, reply.OldState, reply.NewState)
	}
	for _, cb := range dh.deviceCbk {
		go func(callback common.DeviceCallback) {
			defer dh.panicRecover()
			_ = callback.StateChanged(name, reply)
		}(cb)
	}
	return nil
}
func (dh *DeviceHandler) ActionPrompt(name string, reply *common.DevicePrompt) error {
	if dh.log != nil {
		dh.log.Debug("DeviceHandler.ActionPrompt dev:%s, action:%s, prompt:%s",
			name, reply.Action, reply.Prompt)
	}
	for _, cb := range dh.deviceCbk {
		go func(callback common.DeviceCallback) {
			defer dh.panicRecover()
			_ = callback.ActionPrompt(name, reply)
		}(cb)
	}
	return nil
}
func (dh *DeviceHandler) ReaderReturn(name string, reply *common.DeviceInform) error {
	if dh.log != nil {
		dh.log.Debug("DeviceHandler.ReaderReturn dev:%s, action:%s, info:%s",
			name, reply.Action, reply.Inform)
	}
	for _, cb := range dh.deviceCbk {
		go func(callback common.DeviceCallback) {
			defer dh.panicRecover()
			_ = callback.ReaderReturn(name, reply)
		}(cb)
	}
	return nil
}

// Implementation of common.PrinterCallback
func (dh *DeviceHandler) PrinterProgress(name string, reply *common.PrinterProgress) error {
	if dh.log != nil {
		dh.log.Debug("DeviceHandler.PrinterProgress dev:%s, done:%d, From:%d",
			name, reply.PageDone, reply.PagesAll)
	}
	for _, cb := range dh.printerCbk {
		go func(callback common.PrinterCallback) {
			defer dh.panicRecover()
			_ = callback.PrinterProgress(name, reply)
		}(cb)
	}
	return nil
}

// Implementation of common.ReaderCallback
func (dh *DeviceHandler) CardPosition(name string, reply *common.ReaderCardPos) error {
	if dh.log != nil {
		dh.log.Debug("DeviceHandler.CardPosition dev:%s, Position:%d",
			name, reply.Position)
	}
	for _, cb := range dh.readerCbk {
		go func(callback common.ReaderCallback) {
			defer dh.panicRecover()
			_ = callback.CardPosition(name, reply)
		}(cb)
	}
	return nil
}
func (dh *DeviceHandler) CardDescription(name string, reply *common.ReaderCardInfo) error {
	if dh.log != nil {
		dh.log.Debug("DeviceHandler.CardDescription dev:%s, CardPAN:%s, ExpDate:%s",
			name, reply.CardPan, reply.ExpDate)
	}
	for _, cb := range dh.readerCbk {
		go func(callback common.ReaderCallback) {
			defer dh.panicRecover()
			_ = callback.CardDescription(name, reply)
		}(cb)
	}
	return nil
}
func (dh *DeviceHandler) ChipResponse(name string, reply *common.ReaderChipReply) error {
	if dh.log != nil {
		dh.log.Debug("DeviceHandler.ChipResponse dev:%s, Protocol:%d",
			name, reply.Protocol)
	}
	for _, cb := range dh.readerCbk {
		go func(callback common.ReaderCallback) {
			defer dh.panicRecover()
			_ = callback.ChipResponse(name, reply)
		}(cb)
	}
	return nil
}

// Implementation of common.ValidatorCallback
func (dh *DeviceHandler) NoteAccepted(name string, reply *common.ValidatorAccept) error {
	if dh.log != nil {
		dh.log.Debug("DeviceHandler.NoteAccepted dev:%s, Reply: %s",
			name, reply.String())
	}
	for _, cb := range dh.validatorCbk {
		go func(callback common.ValidatorCallback) {
			defer dh.panicRecover()
			_ = callback.NoteAccepted(name, reply)
		}(cb)
	}
	return nil
}
func (dh *DeviceHandler) CashIsStored(name string, reply *common.ValidatorAccept) error {
	if dh.log != nil {
		dh.log.Debug("DeviceHandler.CashIsStored dev:%s, Reply: %s",
			name, reply.String())
	}
	for _, cb := range dh.validatorCbk {
		go func(callback common.ValidatorCallback) {
			defer dh.panicRecover()
			_ = callback.CashIsStored(name, reply)
		}(cb)
	}
	return nil
}
func (dh *DeviceHandler) CashReturned(name string, reply *common.ValidatorAccept) error {
	if dh.log != nil {
		dh.log.Debug("DeviceHandler.CashReturned dev:%s, Reply: %s",
			name, reply.String())
	}
	for _, cb := range dh.validatorCbk {
		go func(callback common.ValidatorCallback) {
			defer dh.panicRecover()
			_ = callback.CashReturned(name, reply)
		}(cb)
	}
	return nil
}
func (dh *DeviceHandler) ValidatorStore(name string, reply *common.ValidatorStore) error {
	if dh.log != nil {
		dh.log.Debug("DeviceHandler.ValidatorStore dev:%s, Reply: %s",
			name, reply.String())
	}
	for _, cb := range dh.validatorCbk {
		go func(callback common.ValidatorCallback) {
			defer dh.panicRecover()
			_ = callback.ValidatorStore(name, reply)
		}(cb)
	}
	return nil
}

// Implementation of common.PinPadCallback
func (dh *DeviceHandler) PinPadReply(name string, reply *common.ReaderPinReply) error {
	if dh.log != nil {
		dh.log.Debug("DeviceHandler.PinPadReply dev:%s, PinLen:%d",
			name, reply.PinLength)
	}
	for _, cb := range dh.pinpadCbk {
		go func(callback common.PinPadCallback) {
			defer dh.panicRecover()
			_ = callback.PinPadReply(name, reply)
		}(cb)
	}
	return nil
}

