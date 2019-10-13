package common

import "fmt"

const (
	CmdNoteAccepted   = "NoteAccepted"
	CmdCashIsStored   = "CashIsStored"
	CmdCashReturned   = "CashReturned"
	CmdValidatorStore = "ValidatorStore"
	CmdInitValidator  = "InitValidator"
	CmdDoValidate     = "DoValidate"
	CmdNoteAccept     = "NoteAccept"
	CmdNoteReject     = "NoteReject"
	CmdStopValidate   = "StopValidate"
	CmdCheckValidator = "CheckValidator"
	CmdClearValidator = "ClearValidator"
)

type ValidNoteList []*ValidatorNote

type ValidatorNote struct {
	Currency DevCurrency
	Count    DevCounter
	Nominal  DevAmount
	Amount   DevAmount
}

func (vn *ValidatorNote) String() string {
	if vn == nil {
		return ""
	}
	str := fmt.Sprintf("Note%5.2f *%3d =%8.2f of %3d (%s) - %s",
		vn.Nominal, vn.Count, vn.Amount, vn.Currency, vn.Currency.IsoCode(), vn.Currency.String())
	return str
}

func (vl ValidNoteList) String() string {
	str := "Validator Note List:"
	for i, note := range vl {
		if note != nil {
			str += fmt.Sprintf("\n    Line:%2d - %s", i, note.String())
		}
	}
	return str
}



type ValidatorStore struct {
	DeviceReply
	List ValidNoteList
}
func (dev *ValidatorStore) String() string {
	if dev == nil { return "" }
	str := fmt.Sprintf("%s, %s",
		dev.DeviceReply.String(), dev.List.String())
	return str
}

type ValidatorAccept struct {
	Currency DevCurrency
	Amount   DevAmount
	Count    DevCounter
}
func (dev *ValidatorAccept) String() string {
	if dev == nil { return "" }
	str := fmt.Sprintf("Count: %d, Amount: %5.2f, Currency: %d (%s) %s",
		dev.Count, dev.Amount, dev.Currency, dev.Currency.IsoCode(), dev.Currency.String())
	return str
}

type ValidatorQuery struct {
	Currency DevCurrency
}
func (dev *ValidatorQuery) String() string {
	if dev == nil { return "" }
	str := fmt.Sprintf("Currency = %s",
		dev.Currency)
	return str
}

//type ValidatorReply struct {
//	Currency DevCurrency
//}

type ValidatorCallback interface {
	NoteAccepted(name string, value *ValidatorAccept) error
	CashIsStored(name string, value *ValidatorAccept) error
	CashReturned(name string, value *ValidatorAccept) error
	ValidatorStore(name string, reply *ValidatorStore) error
}

type ValidatorManager interface {
	InitValidator(name string, query *ValidatorQuery) error
	DoValidate(name string, query *ValidatorQuery) error
	NoteAccept(name string, query *ValidatorQuery) error
	NoteReject(name string, query *ValidatorQuery) error
	StopValidate(name string, query *ValidatorQuery) error
	CheckValidator(name string, query *ValidatorQuery) error
	ClearValidator(name string, query *ValidatorQuery) error
}
