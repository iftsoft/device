package common

type DeviceReply struct {
	DevName string
	Command string
	DevType EnumDevType
	Reply   interface{}
}
type DeviceError struct {
	DevName string
	ErrCode EnumDevError
	Action  EnumDevAction
}
type DeviceState struct {
	DevName  string
	NewState EnumDevState
	OldState EnumDevState
}
type DevicePrompt struct {
	DevName string
	Prompt  EnumDevPrompt
	Action  EnumDevAction
}
type DeviceAccept struct {
	DevName  string
	Amount   DevAmount
	Currency DevCurrency
	Count    DevCounter
}
type DeviceReader struct {
	DevName string
	Inform  string
	Action  EnumDevAction
}

type DeviceCallback interface {
	CommandReply(reply *DeviceReply)
}

type DeviceNotifier interface {
	DeviceCallback
	ExecuteError(value *DeviceError)
	StateChanged(value *DeviceState)
	ActionPrompt(value *DevicePrompt)
	NoteAccepted(value *DeviceAccept)
	CashIsStored(value *DeviceAccept)
	ReaderReturn(value *DeviceReader)
}

type DeviceManager interface {
	Cancel(dcb DeviceCallback) error
	Initialization(dcb DeviceCallback) error
	Reset(dcb DeviceCallback) error
	Status(dcb DeviceCallback) error
}
