package common

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
	CommandReply(reply *SystemReply) error
}

type SystemManager interface {
	Config(query *SystemQuery) error
	Inform(query *SystemQuery) error
	Start(query *SystemQuery) error
	Stop(query *SystemQuery) error
	Restart(query *SystemQuery) error
}
