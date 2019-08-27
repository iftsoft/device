package common

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
	NoteAccepted(name string, value *ValidatorAccept)
	CashIsStored(name string, value *ValidatorAccept)
	ValidatorReturn(name string, reply *ValidatorReply)
}

type ValidatorManager interface {
	Setup(name string, query *ValidatorQuery) error
	Validate(name string, query *ValidatorQuery) error
	Accept(name string, query *ValidatorQuery) error
	Reject(name string, query *ValidatorQuery) error
	Finish(name string, query *ValidatorQuery) error
}

type DeviceReader struct {
	DevName string
	Inform  string
	Action  EnumDevAction
}

type DeviceNotifier interface {
	ReaderReturn(name string, value *DeviceReader)
}
