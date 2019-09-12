package common

const (
	CmdPinPadReply   = "PinPadReply"
	CmdReadPIN       = "ReadPIN"
	CmdLoadMasterKey = "LoadMasterKey"
	CmdLoadWorkKey   = "LoadWorkKey"
	CmdTestMasterKey = "TestMasterKey"
	CmdTestWorkKey   = "TestWorkKey"
)

type ReaderPinQuery struct {
	UseMode  int16
	KeyType  int16
	KeyIndex int16
	KeyValue []byte
	CardPan  string
}

type ReaderPinReply struct {
	DeviceReply
	PinLength int16
	PinBlock  []byte
}

type PinPadCallback interface {
	PinPadReply(name string, query *ReaderPinReply) error
}

type PinPadManager interface {
	ReadPIN(name string, query *ReaderPinQuery) error
	LoadMasterKey(name string, query *ReaderPinQuery) error
	LoadWorkKey(name string, query *ReaderPinQuery) error
	TestMasterKey(name string, query *ReaderPinQuery) error
	TestWorkKey(name string, query *ReaderPinQuery) error
}
