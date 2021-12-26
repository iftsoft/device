package common

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

func (set *CallbackSet) InitCallbacks(complex ComplexCallback) DevScopeMask {
	var mask DevScopeMask
	if complex == nil {
		return mask
	}
	// Setup System scope interface
	set.System = complex.GetSystemCallback()
	if set.System != nil {
		mask |= ScopeFlagSystem
	}
	// Setup Device scope interface
	set.Device = complex.GetDeviceCallback()
	if set.Device != nil {
		mask |= ScopeFlagDevice
	}
	// Setup Printer scope interface
	set.Printer = complex.GetPrinterCallback()
	if set.Printer != nil {
		mask |= ScopeFlagPrinter
	}
	// Setup Reader scope interface
	set.Reader = complex.GetReaderCallback()
	if set.Reader != nil {
		mask |= ScopeFlagReader
	}
	// Setup Validator scope interface
	set.Validator = complex.GetValidatorCallback()
	if set.Validator != nil {
		mask |= ScopeFlagValidator
	}
	// Setup PinPad scope interface
	set.PinPad = complex.GetPinPadCallback()
	if set.PinPad != nil {
		mask |= ScopeFlagPinPad
	}
	return mask
}

func (set *ManagerSet) InitManagers(complex ComplexManager) DevScopeMask {
	var mask DevScopeMask
	if complex == nil {
		return mask
	}
	// Setup System scope interface
	set.System = complex.GetSystemManager()
	if set.System != nil {
		mask |= ScopeFlagSystem
	}
	// Setup Device scope interface
	set.Device = complex.GetDeviceManager()
	if set.Device != nil {
		mask |= ScopeFlagDevice
	}
	// Setup Printer scope interface
	set.Printer = complex.GetPrinterManager()
	if set.Printer != nil {
		mask |= ScopeFlagPrinter
	}
	// Setup Reader scope interface
	set.Reader = complex.GetReaderManager()
	if set.Reader != nil {
		mask |= ScopeFlagReader
	}
	// Setup Validator scope interface
	set.Validator = complex.GetValidatorManager()
	if set.Validator != nil {
		mask |= ScopeFlagValidator
	}
	// Setup PinPad scope interface
	set.PinPad = complex.GetPinPadManager()
	if set.PinPad != nil {
		mask |= ScopeFlagPinPad
	}
	return mask
}
