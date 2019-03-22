package common

type SystemReply struct {
	DevName string
	Command string
	DevType EnumDevType
	Reply   interface{}
}

type SystemCallback interface {
	CommandReply(reply *SystemReply)
}

type SystemManager interface {
	Config(scb SystemCallback) error
	Inform(scb SystemCallback) error
	Start(scb SystemCallback) error
	Stop(scb SystemCallback) error
	Restart(scb SystemCallback) error
}
