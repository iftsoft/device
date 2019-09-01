package common

type EnumSystemError int16
type EnumSystemState int16

// System state codes
const (
	SysErrSuccess EnumSystemError = iota
	SysErrSystemFail
	SysErrDeviceFail
)

// System state codes
const (
	SysStateUndefined EnumSystemState = iota
	SysStateRunning
	SysStateStopped
)

const (
	CmdSystemReply   = "SystemReply"
	CmdSystemHealth  = "SystemHealth"
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
	Message string
	Error   EnumSystemError
	State   EnumSystemState
}

type SystemMetrics struct {
	Uptime uint32
	Counts map[string]uint32
	Totals map[string]uint32
	Topics map[string]string
}

type SystemHealth struct {
	Moment  int64
	Error   EnumSystemError
	State   EnumSystemState
	Metrics SystemMetrics
}

type SystemCallback interface {
	SystemReply(name string, reply *SystemReply) error
	SystemHealth(name string, reply *SystemHealth) error
}

type SystemManager interface {
	Config(name string, query *SystemQuery) error
	Inform(name string, query *SystemQuery) error
	Start(name string, query *SystemQuery) error
	Stop(name string, query *SystemQuery) error
	Restart(name string, query *SystemQuery) error
}
