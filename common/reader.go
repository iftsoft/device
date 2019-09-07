package common

const (
	CmdCardDescription = "CardDescription"
	CmdChipResponse    = "ChipResponse"
	CmdPinPadReply     = "PinPadReply"
	CmdChipGetATR      = "ChipGetATR"
	CmdChipPowerOff    = "ChipPowerOff"
	CmdChipCommand     = "ChipCommand"
	CmdReadPIN         = "ReadPIN"
	CmdLoadMasterKey   = "LoadMasterKey"
	CmdLoadWorkKey     = "LoadWorkKey"
	CmdTestMasterKey   = "TestMasterKey"
	CmdTestWorkKey     = "TestWorkKey"
)

type ReaderCardInfo struct {
	Track1  string
	Track2  string
	Track3  string
	RawData string
	CardPan string
	ExpDate string
	Holder  string
}

type ReaderChipQuery struct {
	Protocol int16
	Query    []byte
}

type ReaderChipReply struct {
	Protocol int16
	Reply    []byte
	ErrCode  EnumDevError
	ErrText  string
}

type ReaderPinQuery struct {
	UseMode  int16
	KeyType  int16
	KeyIndex int16
	KeyValue []byte
	CardPan  string
}

type ReaderPinReply struct {
	PinLength int16
	PinBlock  []byte
	ErrCode   EnumDevError
	ErrText   string
}

type ReaderCallback interface {
	CardDescription(name string, query *ReaderCardInfo) error
	ChipResponse(name string, query *ReaderChipReply) error
	PinPadReply(name string, query *ReaderPinReply) error
}

type ReaderManager interface {
	ChipGetAT(name string, query *DeviceQuery) error
	ChipPowerOff(name string, query *DeviceQuery) error
	ChipCommand(name string, query *ReaderChipQuery) error
	ReadPIN(name string, query *ReaderPinQuery) error
	LoadMasterKey(name string, query *ReaderPinQuery) error
	LoadWorkKey(name string, query *ReaderPinQuery) error
	TestMasterKey(name string, query *ReaderPinQuery) error
	TestWorkKey(name string, query *ReaderPinQuery) error
}
