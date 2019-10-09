package common

const (
	CmdDeviceReply  = "DeviceReply"
	CmdExecuteError = "ExecuteError"
	CmdStateChanged = "StateChanged"
	CmdActionPrompt = "ActionPrompt"
	CmdReaderReturn = "ReaderReturn"
	CmdDeviceCancel = "Cancel"
	CmdDeviceReset  = "Reset"
	CmdDeviceStatus = "Status"
	CmdRunAction    = "RunAction"
	CmdStopAction   = "StopAction"
)

type DeviceQuery struct {
	Timeout int32
}

type DeviceReply struct {
	Command string
	DeviceError
}

type DeviceError struct {
	DevState EnumDevState
	ErrCode  EnumDevError
	ErrText  string
}
type DeviceState struct {
	NewState EnumDevState
	OldState EnumDevState
}
type DevicePrompt struct {
	Prompt EnumDevPrompt
	Action EnumDevAction
}
type DeviceInform struct {
	Inform string
	Action EnumDevAction
}

type DeviceCallback interface {
	DeviceReply(name string, reply *DeviceReply) error
	ExecuteError(name string, value *DeviceError) error
	StateChanged(name string, value *DeviceState) error
	ActionPrompt(name string, value *DevicePrompt) error
	ReaderReturn(name string, value *DeviceInform) error
}

type DeviceManager interface {
	Cancel(name string, query *DeviceQuery) error
	Reset(name string, query *DeviceQuery) error
	Status(name string, query *DeviceQuery) error
	RunAction(name string, query *DeviceQuery) error
	StopAction(name string, query *DeviceQuery) error
}
