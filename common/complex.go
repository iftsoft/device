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
