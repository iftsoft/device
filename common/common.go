package common

import "fmt"

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
func (dev *DeviceQuery) String() string {
	if dev == nil { return "" }
	str := fmt.Sprintf("Timeout = %d",
		dev.Timeout)
	return str
}

type DeviceError struct {
	Action   EnumDevAction	`json:"action"`
	DevState EnumDevState	`json:"devState"`
	ErrCode  EnumDevError	`json:"errCode"`
	ErrText  string			`json:"errText"`
}
func (dev *DeviceError) String() string {
	if dev == nil { return "" }
	str := fmt.Sprintf("Action = %s, DevState = %s, ErrCode = %s, ErrText = %s",
		dev.Action, dev.DevState, dev.ErrCode, dev.ErrText)
	return str
}

type DeviceReply struct {
	Command string
	DeviceError
}
func (dev *DeviceReply) String() string {
	if dev == nil { return "" }
	str := fmt.Sprintf("Command = %s, %s",
		dev.Command, dev.DeviceError.String())
	return str
}

type DeviceState struct {
	Action   EnumDevAction	`json:"action"`
	NewState EnumDevState	`json:"newState"`
	OldState EnumDevState	`json:"oldState"`
}
func (dev *DeviceState) String() string {
	if dev == nil { return "" }
	str := fmt.Sprintf("Action = %s, NewState = %s, OldState = %s",
		dev.Action, dev.NewState, dev.OldState)
	return str
}

type DevicePrompt struct {
	Action EnumDevAction	`json:"action"`
	Prompt EnumDevPrompt	`json:"prompt"`
}
func (dev *DevicePrompt) String() string {
	if dev == nil { return "" }
	str := fmt.Sprintf("Action = %s, Prompt = %s",
		dev.Action, dev.Prompt)
	return str
}

type DeviceInform struct {
	Action EnumDevAction	`json:"action"`
	Inform string	        `json:"inform"`
}
func (dev *DeviceInform) String() string {
	if dev == nil { return "" }
	str := fmt.Sprintf("Action = %s, Inform = %s",
		dev.Action, dev.Inform)
	return str
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
