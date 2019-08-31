package common

type EnumSystemState int16

// System state codes
const (
	SysStateUndefined EnumSystemState = iota
	SysStateRunning
	SysStateStopped
)

const (
	CmdSystemReply   = "SystemReply"
	CmdSystemConfig  = "Config"
	CmdSystemInform  = "Inform"
	CmdSystemStart   = "Start"
	CmdSystemStop    = "Stop"
	CmdSystemRestart = "Restart"
)

type SystemQuery struct {
	//	DevName string
}

type SystemReply struct {
	//	DevName string
	Command string
	Error   string
	State   EnumSystemState
}

type SystemCallback interface {
	SystemReply(name string, reply *SystemReply) error
}

type SystemManager interface {
	Config(name string, query *SystemQuery) error
	Inform(name string, query *SystemQuery) error
	Start(name string, query *SystemQuery) error
	Stop(name string, query *SystemQuery) error
	Restart(name string, query *SystemQuery) error
}
