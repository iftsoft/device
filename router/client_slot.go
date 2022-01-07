package router

import "github.com/iftsoft/device/common"

type ClientSlot struct {
	callbacks common.CallbackSet
	managers  common.ManagerSet
	blocked   bool
}

func (cs *ClientSlot) InitSlot(mng common.ComplexManager, cbk common.ComplexCallback) {
	cs.callbacks.InitCallbacks(cbk)
	cs.managers.InitManagers(mng)
	cs.blocked = false
}

// Implementation of common.ComplexManager

func (cs *ClientSlot) GetSystemManager() common.SystemManager {
	return cs
}
func (cs *ClientSlot) GetDeviceManager() common.DeviceManager {
	return cs
}
func (cs *ClientSlot) GetPrinterManager() common.PrinterManager {
	return cs
}
func (cs *ClientSlot) GetReaderManager() common.ReaderManager {
	return cs
}
func (cs *ClientSlot) GetPinPadManager() common.PinPadManager {
	return cs
}
func (cs *ClientSlot) GetValidatorManager() common.ValidatorManager {
	return cs
}

// Implementation of common.SystemManager

func (cs *ClientSlot) SysInform(name string, query *common.SystemQuery) error {
	if cs.managers.System != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.System.SysInform(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}

func (cs *ClientSlot) SysStart(name string, query *common.SystemQuery) error {
	if cs.managers.System != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.System.SysStart(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}

func (cs *ClientSlot) SysStop(name string, query *common.SystemQuery) error {
	if cs.managers.System != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.System.SysStop(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}

func (cs *ClientSlot) SysRestart(name string, query *common.SystemQuery) error {
	if cs.managers.System != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.System.SysRestart(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}

// Implementation of common.DeviceManager

func (cs *ClientSlot) Cancel(name string, query *common.DeviceQuery) error {
	if cs.managers.Device != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.Device.Cancel(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (cs *ClientSlot) Reset(name string, query *common.DeviceQuery) error {
	if cs.managers.Device != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.Device.Reset(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (cs *ClientSlot) Status(name string, query *common.DeviceQuery) error {
	if cs.managers.Device != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.Device.Status(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (cs *ClientSlot) RunAction(name string, query *common.DeviceQuery) error {
	if cs.managers.Device != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.Device.RunAction(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (cs *ClientSlot) StopAction(name string, query *common.DeviceQuery) error {
	if cs.managers.Device != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.Device.StopAction(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}

// Implementation of common.PrinterManager

func (cs *ClientSlot) InitPrinter(name string, query *common.PrinterSetup) error {
	if cs.managers.Printer != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.Printer.InitPrinter(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (cs *ClientSlot) PrintText(name string, query *common.PrinterQuery) error {
	if cs.managers.Printer != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.Printer.PrintText(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}

// Implementation of common.ReaderManager

func (cs *ClientSlot) EnterCard(name string, query *common.DeviceQuery) error {
	if cs.managers.Reader != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.Reader.EnterCard(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (cs *ClientSlot) EjectCard(name string, query *common.DeviceQuery) error {
	if cs.managers.Reader != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.Reader.EjectCard(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (cs *ClientSlot) CaptureCard(name string, query *common.DeviceQuery) error {
	if cs.managers.Reader != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.Reader.CaptureCard(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (cs *ClientSlot) ReadCard(name string, query *common.DeviceQuery) error {
	if cs.managers.Reader != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.Reader.ReadCard(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (cs *ClientSlot) ChipGetATR(name string, query *common.DeviceQuery) error {
	if cs.managers.Reader != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.Reader.ChipGetATR(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (cs *ClientSlot) ChipPowerOff(name string, query *common.DeviceQuery) error {
	if cs.managers.Reader != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.Reader.ChipPowerOff(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (cs *ClientSlot) ChipCommand(name string, query *common.ReaderChipQuery) error {
	if cs.managers.Reader != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.Reader.ChipCommand(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}

// Implementation of common.ValidatorManager

func (cs *ClientSlot) InitValidator(name string, query *common.ValidatorQuery) error {
	if cs.managers.Validator != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.Validator.InitValidator(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (cs *ClientSlot) DoValidate(name string, query *common.ValidatorQuery) error {
	if cs.managers.Validator != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.Validator.DoValidate(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (cs *ClientSlot) NoteAccept(name string, query *common.ValidatorQuery) error {
	if cs.managers.Validator != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.Validator.NoteAccept(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (cs *ClientSlot) NoteReturn(name string, query *common.ValidatorQuery) error {
	if cs.managers.Validator != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.Validator.NoteReturn(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (cs *ClientSlot) StopValidate(name string, query *common.ValidatorQuery) error {
	if cs.managers.Validator != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.Validator.StopValidate(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (cs *ClientSlot) CheckValidator(name string, query *common.ValidatorQuery) error {
	if cs.managers.Validator != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.Validator.CheckValidator(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (cs *ClientSlot) ClearValidator(name string, query *common.ValidatorQuery) error {
	if cs.managers.Validator != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.Validator.ClearValidator(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}

// Implementation of common.PinPadManager

func (cs *ClientSlot) ReadPIN(name string, query *common.ReaderPinQuery) error {
	if cs.managers.PinPad != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.PinPad.ReadPIN(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (cs *ClientSlot) LoadMasterKey(name string, query *common.ReaderPinQuery) error {
	if cs.managers.PinPad != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.PinPad.LoadMasterKey(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (cs *ClientSlot) LoadWorkKey(name string, query *common.ReaderPinQuery) error {
	if cs.managers.PinPad != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.PinPad.LoadWorkKey(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (cs *ClientSlot) TestMasterKey(name string, query *common.ReaderPinQuery) error {
	if cs.managers.PinPad != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.PinPad.TestMasterKey(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
func (cs *ClientSlot) TestWorkKey(name string, query *common.ReaderPinQuery) error {
	if cs.managers.PinPad != nil {
		if cs.blocked {
			return common.NewError(common.DevErrorNoAccess, "")
		}
		return cs.managers.PinPad.TestWorkKey(name, query)
	}
	return common.NewError(common.DevErrorNotInitialized, "")
}
