package router

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
)

type Multiplexer struct {
	router  DeviceRouter
	clients map[string]common.CallbackSet
	config  *config.AppConfig
	log     *core.LogAgent
}

func NewMultiplexer(config *config.AppConfig) Multiplexer {
	mx := Multiplexer{}
	mx.initMultiplexer(config)
	return mx
}

func (mx *Multiplexer) initMultiplexer(config *config.AppConfig) {
	mx.config = config
	mx.clients = make(map[string]common.CallbackSet)
	mx.log = core.GetLogAgent(core.LogLevelTrace, "Router")
	mx.router.initRouter(mx.log, config, mx)
}

func (mx *Multiplexer) Startup() {
	mx.router.startupRouter()
}

func (mx *Multiplexer) Cleanup() {
	for name, _ := range mx.clients {
		delete(mx.clients, name)
	}
	mx.router.cleanupRouter()
}

func (mx *Multiplexer) AttachClient(name string, cbk common.ComplexCallback) common.ComplexManager {
	set := common.CallbackSet{}
	set.InitCallbacks(cbk)
	mx.clients[name] = set
	return &mx.router
}

func (mx *Multiplexer) DetachClient(name string) {
	delete(mx.clients, name)
}

// Implementation of common.ComplexCallback

func (mx *Multiplexer) GetSystemCallback() common.SystemCallback {
	return mx
}
func (mx *Multiplexer) GetDeviceCallback() common.DeviceCallback {
	return mx
}
func (mx *Multiplexer) GetPrinterCallback() common.PrinterCallback {
	return mx
}
func (mx *Multiplexer) GetReaderCallback() common.ReaderCallback {
	return mx
}
func (mx *Multiplexer) GetPinPadCallback() common.PinPadCallback {
	return mx
}
func (mx *Multiplexer) GetValidatorCallback() common.ValidatorCallback {
	return mx
}

// Implementation of common.SystemCallback

func (mx *Multiplexer) SystemReply(name string, reply *common.SystemReply) error {
	for _, set := range mx.clients {
		if set.System != nil {
			go func(cb common.SystemCallback) {
				defer mx.panicRecover()
				_ = cb.SystemReply(name, reply)
			}(set.System)
		}
	}
	return nil
}
func (mx *Multiplexer) SystemHealth(name string, reply *common.SystemHealth) error {
	for _, set := range mx.clients {
		if set.System != nil {
			go func(cb common.SystemCallback) {
				defer mx.panicRecover()
				_ = cb.SystemHealth(name, reply)
			}(set.System)
		}
	}
	return nil
}

// Implementation of common.DeviceCallback

func (mx *Multiplexer) DeviceReply(name string, reply *common.DeviceReply) error {
	for _, set := range mx.clients {
		if set.Device != nil {
			go func(cb common.DeviceCallback) {
				defer mx.panicRecover()
				_ = cb.DeviceReply(name, reply)
			}(set.Device)
		}
	}
	return nil
}
func (mx *Multiplexer) ExecuteError(name string, reply *common.DeviceError) error {
	for _, set := range mx.clients {
		if set.Device != nil {
			go func(cb common.DeviceCallback) {
				defer mx.panicRecover()
				_ = cb.ExecuteError(name, reply)
			}(set.Device)
		}
	}
	return nil
}
func (mx *Multiplexer) StateChanged(name string, reply *common.DeviceState) error {
	for _, set := range mx.clients {
		if set.Device != nil {
			go func(cb common.DeviceCallback) {
				defer mx.panicRecover()
				_ = cb.StateChanged(name, reply)
			}(set.Device)
		}
	}
	return nil
}
func (mx *Multiplexer) ActionPrompt(name string, reply *common.DevicePrompt) error {
	for _, set := range mx.clients {
		if set.Device != nil {
			go func(cb common.DeviceCallback) {
				defer mx.panicRecover()
				_ = cb.ActionPrompt(name, reply)
			}(set.Device)
		}
	}
	return nil
}
func (mx *Multiplexer) ReaderReturn(name string, reply *common.DeviceInform) error {
	for _, set := range mx.clients {
		if set.Device != nil {
			go func(cb common.DeviceCallback) {
				defer mx.panicRecover()
				_ = cb.ReaderReturn(name, reply)
			}(set.Device)
		}
	}
	return nil
}

// Implementation of common.PrinterCallback

func (mx *Multiplexer) PrinterProgress(name string, reply *common.PrinterProgress) error {
	for _, set := range mx.clients {
		if set.Printer != nil {
			go func(cb common.PrinterCallback) {
				defer mx.panicRecover()
				_ = cb.PrinterProgress(name, reply)
			}(set.Printer)
		}
	}
	return nil
}

// Implementation of common.ReaderCallback

func (mx *Multiplexer) CardPosition(name string, reply *common.ReaderCardPos) error {
	for _, set := range mx.clients {
		if set.Reader != nil {
			go func(cb common.ReaderCallback) {
				defer mx.panicRecover()
				_ = cb.CardPosition(name, reply)
			}(set.Reader)
		}
	}
	return nil
}
func (mx *Multiplexer) CardDescription(name string, reply *common.ReaderCardInfo) error {
	for _, set := range mx.clients {
		if set.Reader != nil {
			go func(cb common.ReaderCallback) {
				defer mx.panicRecover()
				_ = cb.CardDescription(name, reply)
			}(set.Reader)
		}
	}
	return nil
}
func (mx *Multiplexer) ChipResponse(name string, reply *common.ReaderChipReply) error {
	for _, set := range mx.clients {
		if set.Reader != nil {
			go func(cb common.ReaderCallback) {
				defer mx.panicRecover()
				_ = cb.ChipResponse(name, reply)
			}(set.Reader)
		}
	}
	return nil
}

// Implementation of common.ValidatorCallback

func (mx *Multiplexer) NoteAccepted(name string, reply *common.ValidatorAccept) error {
	for _, set := range mx.clients {
		if set.Validator != nil {
			go func(cb common.ValidatorCallback) {
				defer mx.panicRecover()
				_ = cb.NoteAccepted(name, reply)
			}(set.Validator)
		}
	}
	return nil
}
func (mx *Multiplexer) CashIsStored(name string, reply *common.ValidatorAccept) error {
	for _, set := range mx.clients {
		if set.Validator != nil {
			go func(cb common.ValidatorCallback) {
				defer mx.panicRecover()
				_ = cb.CashIsStored(name, reply)
			}(set.Validator)
		}
	}
	return nil
}
func (mx *Multiplexer) CashReturned(name string, reply *common.ValidatorAccept) error {
	for _, set := range mx.clients {
		if set.Validator != nil {
			go func(cb common.ValidatorCallback) {
				defer mx.panicRecover()
				_ = cb.CashReturned(name, reply)
			}(set.Validator)
		}
	}
	return nil
}
func (mx *Multiplexer) ValidatorStore(name string, reply *common.ValidatorStore) error {
	for _, set := range mx.clients {
		if set.Validator != nil {
			go func(cb common.ValidatorCallback) {
				defer mx.panicRecover()
				_ = cb.ValidatorStore(name, reply)
			}(set.Validator)
		}
	}
	return nil
}

// Implementation of common.ReaderCallback

func (mx *Multiplexer) PinPadReply(name string, reply *common.ReaderPinReply) error {
	for _, set := range mx.clients {
		if set.PinPad != nil {
			go func(cb common.PinPadCallback) {
				defer mx.panicRecover()
				_ = cb.PinPadReply(name, reply)
			}(set.PinPad)
		}
	}
	return nil
}

func (mx *Multiplexer) panicRecover() {
	if r := recover(); r != nil {
		if mx.log != nil {
			mx.log.Panic("Panic is recovered: %+v", r)
		}
	}
}
