package common

const (
	CmdDeviceReply  = "DeviceReply"
	CmdExecuteError = "ExecuteError"
	CmdStateChanged = "StateChanged"
	CmdActionPrompt = "ActionPrompt"
	CmdDeviceCancel = "Cancel"
	CmdDeviceReset  = "Reset"
	CmdDeviceStatus = "Status"
)

type DeviceQuery struct {
	Timeout int32
}

type DeviceReply struct {
	Command  string
	DevType  EnumDevType
	DevState EnumDevState
	ErrCode  EnumDevError
	ErrText  string
}

type DeviceError struct {
	Action  EnumDevAction
	ErrCode EnumDevError
	ErrText string
}
type DeviceState struct {
	NewState EnumDevState
	OldState EnumDevState
}
type DevicePrompt struct {
	Prompt EnumDevPrompt
	Action EnumDevAction
}

type DeviceCallback interface {
	CommandReply(name string, reply *DeviceReply) error
	ExecuteError(name string, value *DeviceError) error
	StateChanged(name string, value *DeviceState) error
	ActionPrompt(name string, value *DevicePrompt) error
}

type DeviceManager interface {
	Cancel(name string, query *DeviceQuery) error
	Reset(name string, query *DeviceQuery) error
	Status(name string, query *DeviceQuery) error
}
