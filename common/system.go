package common

type EnumSystemError int16
type EnumSystemState int16

const (
	CmdSystemReply     = "SystemReply"
	CmdSystemHealth    = "SystemHealth"
	CmdSystemTerminate = "Terminate"
	CmdSystemInform    = "SysInform"
	CmdSystemStart     = "SysStart"
	CmdSystemStop      = "SysStop"
	CmdSystemRestart   = "SysRestart"
)

// System state codes
const (
	SysErrSuccess EnumSystemError = iota
	SysErrSystemFail
	SysErrDeviceFail
)
func (e EnumSystemError) String() string {
	switch e {
	case SysErrSuccess:			return "Success"
	case SysErrSystemFail:		return "System fail"
	case SysErrDeviceFail:		return "Device fail"
	default:					return "Undefined"
	}
}


// System state codes
const (
	SysStateUndefined EnumSystemState = iota
	SysStateRunning
	SysStateStopped
	SysStateFailed
)
func (e EnumSystemState) String() string {
	switch e {
	case SysStateUndefined:		return "Undefined"
	case SysStateRunning:		return "Running"
	case SysStateStopped:		return "Stopped"
	case SysStateFailed:		return "Failed"
	default:					return "Unknown"
	}
}


type SystemQuery struct {
}

type SystemReply struct {
	Command string
	Message string
	Error   EnumSystemError
	State   EnumSystemState
}

type SystemMetrics struct {
	Uptime   int64
	DevError EnumDevError
	DevState EnumDevState
	Counts   map[string]uint32
	Totals   map[string]float32
	Topics   map[string]string
}

type SystemHealth struct {
	Moment  int64
	Error   EnumSystemError
	State   EnumSystemState
	Metrics SystemMetrics
}

func NewSystemHealth() *SystemHealth {
	sh := &SystemHealth{
		Moment: 0,
		Error: 0,
		State: 0,
		Metrics: SystemMetrics{
			Uptime:   0,
			DevError: 0,
			DevState: 0,
			Counts:   make(map[string]uint32),
			Totals:   make(map[string]float32),
			Topics:   make(map[string]string),
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
	SysInform(name string, query *SystemQuery) error
	SysStart(name string, query *SystemQuery) error
	SysStop(name string, query *SystemQuery) error
	SysRestart(name string, query *SystemQuery) error
}
