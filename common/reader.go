package common

const (
	CmdCardPosition    = "CardPosition"
	CmdCardDescription = "CardDescription"
	CmdChipResponse    = "ChipResponse"
	CmdEnterCard       = "EnterCard"
	CmdEjectCard       = "EjectCard"
	CmdCaptureCard     = "CaptureCard"
	CmdReadCard        = "ReadCard"
	CmdChipGetATR      = "ChipGetATR"
	CmdChipPowerOff    = "ChipPowerOff"
	CmdChipCommand     = "ChipCommand"
)

type ReaderCardPos struct {
	Position int16
}

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
	DeviceReply
	Protocol int16
	Reply    []byte
}

type ReaderCallback interface {
	CardPosition(name string, query *ReaderCardPos) error
	CardDescription(name string, query *ReaderCardInfo) error
	ChipResponse(name string, query *ReaderChipReply) error
}

type ReaderManager interface {
	EnterCard(name string, query *DeviceQuery) error
	EjectCard(name string, query *DeviceQuery) error
	CaptureCard(name string, query *DeviceQuery) error
	ReadCard(name string, query *DeviceQuery) error
	ChipGetATR(name string, query *DeviceQuery) error
	ChipPowerOff(name string, query *DeviceQuery) error
	ChipCommand(name string, query *ReaderChipQuery) error
}
