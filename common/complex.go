package common

type Packet struct {
	Command string
	DevName string
	Content interface{}
}

type ComplexCallback interface {
	GetSystemCallback() SystemCallback
	GetDeviceCallback() DeviceCallback
	GetPrinterCallback() PrinterCallback
	GetReaderCallback() ReaderCallback
	GetPinPadCallback() PinPadCallback
	GetValidatorCallback() ValidatorCallback
}

type ComplexManager interface {
	GetSystemManager() SystemManager
	GetDeviceManager() DeviceManager
	GetPrinterManager() PrinterManager
	GetReaderManager() ReaderManager
	GetPinPadManager() PinPadManager
	GetValidatorManager() ValidatorManager
}

type CallbackSet struct {
	System    SystemCallback
	Device    DeviceCallback
	Printer   PrinterCallback
	Reader    ReaderCallback
	PinPad    PinPadCallback
	Validator ValidatorCallback
}

type ManagerSet struct {
	System    SystemManager
	Device    DeviceManager
	Printer   PrinterManager
	Reader    ReaderManager
	PinPad    PinPadManager
	Validator ValidatorManager
}

func (set *CallbackSet) InitCallbacks(complex ComplexCallback) {
	if complex == nil {
		return
	}
	// Setup System scope interface
	set.System = complex.GetSystemCallback()
	// Setup Device scope interface
	set.Device = complex.GetDeviceCallback()
	// Setup Printer scope interface
	set.Printer = complex.GetPrinterCallback()
	// Setup Reader scope interface
	set.Reader = complex.GetReaderCallback()
	// Setup Validator scope interface
	set.Validator = complex.GetValidatorCallback()
	// Setup PinPad scope interface
	set.PinPad = complex.GetPinPadCallback()
}

func (set *ManagerSet) InitManagers(complex ComplexManager) {
	if complex == nil {
		return
	}
	// Setup System scope interface
	set.System = complex.GetSystemManager()
	// Setup Device scope interface
	set.Device = complex.GetDeviceManager()
	// Setup Printer scope interface
	set.Printer = complex.GetPrinterManager()
	// Setup Reader scope interface
	set.Reader = complex.GetReaderManager()
	// Setup Validator scope interface
	set.Validator = complex.GetValidatorManager()
	// Setup PinPad scope interface
	set.PinPad = complex.GetPinPadManager()
}

func (set *CallbackSet) RunCommand(packet Packet) error {
	switch packet.Command {
	// Run System scope interface
	case CmdSystemReply:
		if set.System != nil {
			if query, ok := packet.Content.(*SystemReply); ok {
				return set.System.SystemReply(packet.DevName, query)
			}
		}
	case CmdSystemHealth:
		if set.System != nil {
			if query, ok := packet.Content.(*SystemHealth); ok {
				return set.System.SystemHealth(packet.DevName, query)
			}
		}
	// Run Device scope interface
	case CmdDeviceReply:
		if set.Device != nil {
			if query, ok := packet.Content.(*DeviceReply); ok {
				return set.Device.DeviceReply(packet.DevName, query)
			}
		}
	case CmdExecuteError:
		if set.Device != nil {
			if query, ok := packet.Content.(*DeviceError); ok {
				return set.Device.ExecuteError(packet.DevName, query)
			}
		}
	case CmdStateChanged:
		if set.Device != nil {
			if query, ok := packet.Content.(*DeviceState); ok {
				return set.Device.StateChanged(packet.DevName, query)
			}
		}
	case CmdActionPrompt:
		if set.Device != nil {
			if query, ok := packet.Content.(*DevicePrompt); ok {
				return set.Device.ActionPrompt(packet.DevName, query)
			}
		}
	case CmdReaderReturn:
		if set.Device != nil {
			if query, ok := packet.Content.(*DeviceInform); ok {
				return set.Device.ReaderReturn(packet.DevName, query)
			}
		}
	// Run Printer scope interface
	case CmdPrinterProgress:
		if set.Printer != nil {
			if query, ok := packet.Content.(*PrinterProgress); ok {
				return set.Printer.PrinterProgress(packet.DevName, query)
			}
		}
	// Run Reader scope interface
	case CmdCardPosition:
		if set.Reader != nil {
			if query, ok := packet.Content.(*ReaderCardPos); ok {
				return set.Reader.CardPosition(packet.DevName, query)
			}
		}
	case CmdCardDescription:
		if set.Reader != nil {
			if query, ok := packet.Content.(*ReaderCardInfo); ok {
				return set.Reader.CardDescription(packet.DevName, query)
			}
		}
	case CmdChipResponse:
		if set.Reader != nil {
			if query, ok := packet.Content.(*ReaderChipReply); ok {
				return set.Reader.ChipResponse(packet.DevName, query)
			}
		}
	// Run Validator scope interface
	case CmdNoteAccepted:
		if set.Validator != nil {
			if query, ok := packet.Content.(*ValidatorAccept); ok {
				return set.Validator.NoteAccepted(packet.DevName, query)
			}
		}
	case CmdCashIsStored:
		if set.Validator != nil {
			if query, ok := packet.Content.(*ValidatorAccept); ok {
				return set.Validator.CashIsStored(packet.DevName, query)
			}
		}
	case CmdCashReturned:
		if set.Validator != nil {
			if query, ok := packet.Content.(*ValidatorAccept); ok {
				return set.Validator.CashReturned(packet.DevName, query)
			}
		}
	case CmdValidatorStore:
		if set.Validator != nil {
			if query, ok := packet.Content.(*ValidatorStore); ok {
				return set.Validator.ValidatorStore(packet.DevName, query)
			}
		}
	// Run PinPad scope interface
	case CmdPinPadReply:
		if set.PinPad != nil {
			if query, ok := packet.Content.(*ReaderPinReply); ok {
				return set.PinPad.PinPadReply(packet.DevName, query)
			}
		}

	}
	return nil
}

func (set *ManagerSet) RunCommand(packet Packet) error {
	switch packet.Command {
	// Run System scope interface
	case CmdSystemInform:
		if set.System != nil {
			if query, ok := packet.Content.(*SystemQuery); ok {
				return set.System.SysInform(packet.DevName, query)
			}
		}
	case CmdSystemStart:
		if set.System != nil {
			if query, ok := packet.Content.(*SystemConfig); ok {
				return set.System.SysStart(packet.DevName, query)
			}
		}
	case CmdSystemStop:
		if set.System != nil {
			if query, ok := packet.Content.(*SystemQuery); ok {
				return set.System.SysStop(packet.DevName, query)
			}
		}
	case CmdSystemRestart:
		if set.System != nil {
			if query, ok := packet.Content.(*SystemConfig); ok {
				return set.System.SysRestart(packet.DevName, query)
			}
		}
	// Run Device scope interface
	case CmdDeviceCancel:
		if set.Device != nil {
			if query, ok := packet.Content.(*DeviceQuery); ok {
				return set.Device.Cancel(packet.DevName, query)
			}
		}
	case CmdDeviceReset:
		if set.Device != nil {
			if query, ok := packet.Content.(*DeviceQuery); ok {
				return set.Device.Reset(packet.DevName, query)
			}
		}
	case CmdDeviceStatus:
		if set.Device != nil {
			if query, ok := packet.Content.(*DeviceQuery); ok {
				return set.Device.Status(packet.DevName, query)
			}
		}
	case CmdRunAction:
		if set.Device != nil {
			if query, ok := packet.Content.(*DeviceQuery); ok {
				return set.Device.RunAction(packet.DevName, query)
			}
		}
	case CmdStopAction:
		if set.Device != nil {
			if query, ok := packet.Content.(*DeviceQuery); ok {
				return set.Device.StopAction(packet.DevName, query)
			}
		}
	// Run Printer scope interface
	case CmdInitPrinter:
		if set.Printer != nil {
			if query, ok := packet.Content.(*PrinterSetup); ok {
				return set.Printer.InitPrinter(packet.DevName, query)
			}
		}
	case CmdPrintText:
		if set.Printer != nil {
			if query, ok := packet.Content.(*PrinterQuery); ok {
				return set.Printer.PrintText(packet.DevName, query)
			}
		}
	// Run Reader scope interface
	case CmdEnterCard:
		if set.Reader != nil {
			if query, ok := packet.Content.(*DeviceQuery); ok {
				return set.Reader.EnterCard(packet.DevName, query)
			}
		}
	case CmdEjectCard:
		if set.Reader != nil {
			if query, ok := packet.Content.(*DeviceQuery); ok {
				return set.Reader.EjectCard(packet.DevName, query)
			}
		}
	case CmdCaptureCard:
		if set.Reader != nil {
			if query, ok := packet.Content.(*DeviceQuery); ok {
				return set.Reader.CaptureCard(packet.DevName, query)
			}
		}
	case CmdReadCard:
		if set.Reader != nil {
			if query, ok := packet.Content.(*DeviceQuery); ok {
				return set.Reader.ReadCard(packet.DevName, query)
			}
		}
	case CmdChipGetATR:
		if set.Reader != nil {
			if query, ok := packet.Content.(*DeviceQuery); ok {
				return set.Reader.ChipGetATR(packet.DevName, query)
			}
		}
	case CmdChipPowerOff:
		if set.Reader != nil {
			if query, ok := packet.Content.(*DeviceQuery); ok {
				return set.Reader.ChipPowerOff(packet.DevName, query)
			}
		}
	case CmdChipCommand:
		if set.Reader != nil {
			if query, ok := packet.Content.(*ReaderChipQuery); ok {
				return set.Reader.ChipCommand(packet.DevName, query)
			}
		}
	// Run Validator scope interface
	case CmdInitValidator:
		if set.Validator != nil {
			if query, ok := packet.Content.(*ValidatorQuery); ok {
				return set.Validator.InitValidator(packet.DevName, query)
			}
		}
	case CmdDoValidate:
		if set.Validator != nil {
			if query, ok := packet.Content.(*ValidatorQuery); ok {
				return set.Validator.DoValidate(packet.DevName, query)
			}
		}
	case CmdNoteAccept:
		if set.Validator != nil {
			if query, ok := packet.Content.(*ValidatorQuery); ok {
				return set.Validator.NoteAccept(packet.DevName, query)
			}
		}
	case CmdNoteReturn:
		if set.Validator != nil {
			if query, ok := packet.Content.(*ValidatorQuery); ok {
				return set.Validator.NoteReturn(packet.DevName, query)
			}
		}
	case CmdStopValidate:
		if set.Validator != nil {
			if query, ok := packet.Content.(*ValidatorQuery); ok {
				return set.Validator.StopValidate(packet.DevName, query)
			}
		}
	case CmdCheckValidator:
		if set.Validator != nil {
			if query, ok := packet.Content.(*ValidatorQuery); ok {
				return set.Validator.CheckValidator(packet.DevName, query)
			}
		}
	case CmdClearValidator:
		if set.Validator != nil {
			if query, ok := packet.Content.(*ValidatorQuery); ok {
				return set.Validator.ClearValidator(packet.DevName, query)
			}
		}
	// Run PinPad scope interface
	case CmdReadPIN:
		if set.PinPad != nil {
			if query, ok := packet.Content.(*ReaderPinQuery); ok {
				return set.PinPad.ReadPIN(packet.DevName, query)
			}
		}
	case CmdLoadMasterKey:
		if set.PinPad != nil {
			if query, ok := packet.Content.(*ReaderPinQuery); ok {
				return set.PinPad.LoadMasterKey(packet.DevName, query)
			}
		}
	case CmdLoadWorkKey:
		if set.PinPad != nil {
			if query, ok := packet.Content.(*ReaderPinQuery); ok {
				return set.PinPad.LoadWorkKey(packet.DevName, query)
			}
		}
	case CmdTestMasterKey:
		if set.PinPad != nil {
			if query, ok := packet.Content.(*ReaderPinQuery); ok {
				return set.PinPad.TestMasterKey(packet.DevName, query)
			}
		}
	case CmdTestWorkKey:
		if set.PinPad != nil {
			if query, ok := packet.Content.(*ReaderPinQuery); ok {
				return set.PinPad.TestWorkKey(packet.DevName, query)
			}
		}

	}
	return nil
}
