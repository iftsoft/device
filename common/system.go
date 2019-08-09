package common

const (
	CmdSystemCommandReply = "CommandReply"
	CmdSystemConfig       = "Config"
	CmdSystemInform       = "Inform"
	CmdSystemStart        = "Start"
	CmdSystemStop         = "Stop"
	CmdSystemRestart      = "Restart"
)

type SystemQuery struct {
	DevName string
}

type SystemReply struct {
	DevName string
	Command string
	DevType EnumDevType
	//	Reply   interface{}
}

type SystemCallback interface {
	CommandReply(name string, reply *SystemReply) error
}

type SystemManager interface {
	Config(name string, query *SystemQuery) error
	Inform(name string, query *SystemQuery) error
	Start(name string, query *SystemQuery) error
	Stop(name string, query *SystemQuery) error
	Restart(name string, query *SystemQuery) error
}
