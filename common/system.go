package common

type EnumSystemError int16
type EnumSystemState int16

const (
	CmdSystemReply  = "SystemReply"
	CmdSystemHealth = "SystemHealth"
	//	CmdSystemTerminate = "Terminate"
	CmdSystemInform  = "SysInform"
	CmdSystemStart   = "SysStart"
	CmdSystemStop    = "SysStop"
	CmdSystemRestart = "SysRestart"
)

// System state codes
const (
	SysErrSuccess EnumSystemError = iota
	SysErrSystemFail
	SysErrDeviceFail
)

func (e EnumSystemError) String() string {
	switch e {
	case SysErrSuccess:
		return "Success"
	case SysErrSystemFail:
		return "System fail"
	case SysErrDeviceFail:
		return "Device fail"
	default:
		return "Undefined"
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
	case SysStateUndefined:
		return "Undefined"
	case SysStateRunning:
		return "Running"
	case SysStateStopped:
		return "Stopped"
	case SysStateFailed:
		return "Failed"
	default:
		return "Unknown"
	}
}

type SystemQuery struct {
}

//type SystemConfig struct {
//	LinkType  uint16 `json:"link_type"`  // 0-none, 1-COM, 2-USB
//	PortName  string `json:"port_name"`  // Serial port name
//	VendorID  uint16 `json:"vendor_id"`  // Device Vendor ID
//	ProductID uint16 `json:"product_id"` // Device Product ID
//}

type SystemReply struct {
	Command string          `json:"command"`
	Message string          `json:"message"`
	Error   EnumSystemError `json:"error"`
	State   EnumSystemState `json:"state"`
}

type SystemMetrics struct {
	Uptime   int64              `json:"uptime"`
	DevError EnumDevError       `json:"dev_error"`
	DevState EnumDevState       `json:"dev_state"`
	Counts   map[string]uint32  `json:"counts"`
	Totals   map[string]float32 `json:"totals"`
	Topics   map[string]string  `json:"topics"`
}

type SystemHealth struct {
	Moment  int64           `json:"moment"`
	Error   EnumSystemError `json:"error"`
	State   EnumSystemState `json:"state"`
	Metrics SystemMetrics   `json:"metrics"`
}

func NewSystemHealth() *SystemHealth {
	sh := &SystemHealth{
		Moment: 0,
		Error:  0,
		State:  0,
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

type GreetingInfo struct {
	DevType   DevTypeMask  `json:"devType"`   // Implemented device types
	Supported DevScopeMask `json:"supported"` // Manager interfaces that driver supported
	Required  DevScopeMask `json:"required"`  // Callback interfaces that driver required
}

type SystemCallback interface {
	SystemReply(name string, reply *SystemReply) error
	SystemHealth(name string, reply *SystemHealth) error
}

type SystemManager interface {
	SysInform(name string, query *SystemQuery) error
	SysStart(name string, query *SystemQuery) error
	SysStop(name string, query *SystemQuery) error
	SysRestart(name string, query *SystemQuery) error
}
