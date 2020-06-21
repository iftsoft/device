package common

import "fmt"

const (
	CmdPinPadReply   = "PinPadReply"
	CmdReadPIN       = "ReadPIN"
	CmdLoadMasterKey = "LoadMasterKey"
	CmdLoadWorkKey   = "LoadWorkKey"
	CmdTestMasterKey = "TestMasterKey"
	CmdTestWorkKey   = "TestWorkKey"
)

type EnumPinKeyType uint16

const (
	PinPadKeyPIN EnumPinKeyType = iota
	PinPadKeyMAC
	PinPadKeyData
)
func (e EnumPinKeyType) String() string {
	switch e {
	case PinPadKeyPIN:			return "Key for PIN"
	case PinPadKeyMAC:			return "Key for MAC"
	case PinPadKeyData:			return "Key for data"
	default:					return "Unknown"
	}
}

type ReaderPinQuery struct {
//	UseMode  int16
	KeyType  EnumPinKeyType	`json:"keyType"`
	KeyIndex uint16	        `json:"keyIndex"`
	KeyValue []byte	        `json:"keyValue"`
	CardPan  string	        `json:"cardPan"`
}

type ReaderPinReply struct {
	DeviceReply
	PinLength uint16	`json:"pinLenght"`
	PinBlock  []byte	`json:"pinBlock"`
}
func (dev *ReaderPinReply) String() string {
	if dev == nil { return "" }
	str := fmt.Sprintf("%s, PinLength = %d",
		dev.DeviceReply.String(), dev.PinLength)
	return str
}


type PinPadCallback interface {
	PinPadReply(name string, reply *ReaderPinReply) error
}

type PinPadManager interface {
	ReadPIN(name string, query *ReaderPinQuery) error
	LoadMasterKey(name string, query *ReaderPinQuery) error
	LoadWorkKey(name string, query *ReaderPinQuery) error
	TestMasterKey(name string, query *ReaderPinQuery) error
	TestWorkKey(name string, query *ReaderPinQuery) error
}
