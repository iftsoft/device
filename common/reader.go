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
	Position int16	`json:"position"`
}

type ReaderCardInfo struct {
	Track1  string	`json:"track1"`
	Track2  string	`json:"track2"`
	Track3  string	`json:"track3"`
	RawData string	`json:"rawData"`
	CardPan string	`json:"cardPan"`
	ExpDate string	`json:"expDate"`
	Holder  string	`json:"holder"`
}

type ReaderChipQuery struct {
	Protocol int16	`json:"protocol"`
	Query    []byte	`json:"query"`
}

type ReaderChipReply struct {
	DeviceReply
	Protocol int16	`json:"protocol"`
	Reply    []byte	`json:"reply"`
}

type ReaderCallback interface {
	CardPosition(name string, value *ReaderCardPos) error
	CardDescription(name string, value *ReaderCardInfo) error
	ChipResponse(name string, reply *ReaderChipReply) error
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
