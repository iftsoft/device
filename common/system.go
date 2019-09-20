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
	CmdSystemReply     = "SystemReply"
	CmdSystemHealth    = "SystemHealth"
	CmdSystemTerminate = "Terminate"
	CmdSystemInform    = "Inform"
	CmdSystemStart     = "Start"
	CmdSystemStop      = "Stop"
	CmdSystemRestart   = "Restart"
)

type SystemQuery struct {
}

type SystemReply struct {
	Command string
	Message string
	Error   EnumSystemError
	State   EnumSystemState
}

type SystemMetrics struct {
	Uptime int64
	Counts map[string]uint32
	Totals map[string]float32
	Topics map[string]string
}

type SystemHealth struct {
	Moment  int64
	Error   EnumSystemError
	State   EnumSystemState
	Metrics SystemMetrics
}

func NewSystemHealth() *SystemHealth {
	sh := &SystemHealth{
		Error: 0,
		State: 0,
		Metrics: SystemMetrics{
			Uptime: 0,
			Counts: make(map[string]uint32),
			Totals: make(map[string]float32),
			Topics: make(map[string]string),
		},
	}
	return sh
}

type SystemCallback interface {
	SystemReply(name string, reply *SystemReply) error
	SystemHealth(name string, reply *SystemHealth) error
}

type SystemManager interface {
	Terminate(name string, query *SystemQuery) error
	Inform(name string, query *SystemQuery) error
	Start(name string, query *SystemQuery) error
	Stop(name string, query *SystemQuery) error
	Restart(name string, query *SystemQuery) error
}
