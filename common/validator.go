package common

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
	CmdCheckStore     = "CheckStore"
	CmdClearStore     = "ClearStore"
)

type ValidatorBill struct {
	Currency DevCurrency
	Count    DevCounter
	Nominal  DevAmount
	Amount   DevAmount
}

type ValidatorStore struct {
	Note []*ValidatorBill
}

type ValidatorAccept struct {
	Currency DevCurrency
	Amount   DevAmount
	Count    DevCounter
}
type ValidatorQuery struct {
	Currency DevCurrency
}
type ValidatorReply struct {
	Currency DevCurrency
}

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
	CheckStore(name string, query *ValidatorQuery) error
	ClearStore(name string, query *ValidatorQuery) error
}
