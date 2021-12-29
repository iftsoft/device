package router

import (
	"github.com/iftsoft/device/common"
)

type CallbackChannel struct {
	reply chan<- common.Packet
}

func (ch *CallbackChannel) initChannel(reply chan<- common.Packet) {
	ch.reply = reply
}

func (ch *CallbackChannel) sendReply(command, devName string, query interface{}) error {
	packet := common.Packet{
		Command: command,
		DevName: devName,
		Content: query,
	}
	ch.reply <- packet
	return nil
}

// Implementation of common.ComplexCallback

func (ch *CallbackChannel) GetSystemCallback() common.SystemCallback {
	return ch
}
func (ch *CallbackChannel) GetDeviceCallback() common.DeviceCallback {
	return ch
}
func (ch *CallbackChannel) GetPrinterCallback() common.PrinterCallback {
	return ch
}
func (ch *CallbackChannel) GetReaderCallback() common.ReaderCallback {
	return ch
}
func (ch *CallbackChannel) GetPinPadCallback() common.PinPadCallback {
	return ch
}
func (ch *CallbackChannel) GetValidatorCallback() common.ValidatorCallback {
	return ch
}

// Implementation of common.SystemCallback

func (ch *CallbackChannel) SystemReply(name string, reply *common.SystemReply) error {
	return ch.sendReply(common.CmdSystemReply, name, reply)
}
func (ch *CallbackChannel) SystemHealth(name string, reply *common.SystemHealth) error {
	return ch.sendReply(common.CmdSystemHealth, name, reply)
}

// Implementation of common.DeviceCallback

func (ch *CallbackChannel) DeviceReply(name string, reply *common.DeviceReply) error {
	return ch.sendReply(common.CmdDeviceReply, name, reply)
}
func (ch *CallbackChannel) ExecuteError(name string, reply *common.DeviceError) error {
	return ch.sendReply(common.CmdExecuteError, name, reply)
}
func (ch *CallbackChannel) StateChanged(name string, reply *common.DeviceState) error {
	return ch.sendReply(common.CmdStateChanged, name, reply)
}
func (ch *CallbackChannel) ActionPrompt(name string, reply *common.DevicePrompt) error {
	return ch.sendReply(common.CmdActionPrompt, name, reply)
}
func (ch *CallbackChannel) ReaderReturn(name string, reply *common.DeviceInform) error {
	return ch.sendReply(common.CmdReaderReturn, name, reply)
}

// Implementation of common.PrinterCallback

func (ch *CallbackChannel) PrinterProgress(name string, reply *common.PrinterProgress) error {
	return ch.sendReply(common.CmdPrinterProgress, name, reply)
}

// Implementation of common.ReaderCallback

func (ch *CallbackChannel) CardPosition(name string, reply *common.ReaderCardPos) error {
	return ch.sendReply(common.CmdCardPosition, name, reply)
}
func (ch *CallbackChannel) CardDescription(name string, reply *common.ReaderCardInfo) error {
	return ch.sendReply(common.CmdCardDescription, name, reply)
}
func (ch *CallbackChannel) ChipResponse(name string, reply *common.ReaderChipReply) error {
	return ch.sendReply(common.CmdChipResponse, name, reply)
}

// Implementation of common.ValidatorCallback

func (ch *CallbackChannel) NoteAccepted(name string, reply *common.ValidatorAccept) error {
	return ch.sendReply(common.CmdNoteAccepted, name, reply)
}
func (ch *CallbackChannel) CashIsStored(name string, reply *common.ValidatorAccept) error {
	return ch.sendReply(common.CmdCashIsStored, name, reply)
}
func (ch *CallbackChannel) CashReturned(name string, reply *common.ValidatorAccept) error {
	return ch.sendReply(common.CmdCashReturned, name, reply)
}
func (ch *CallbackChannel) ValidatorStore(name string, reply *common.ValidatorStore) error {
	return ch.sendReply(common.CmdValidatorStore, name, reply)
}

// Implementation of common.ReaderCallback

func (ch *CallbackChannel) PinPadReply(name string, reply *common.ReaderPinReply) error {
	return ch.sendReply(common.CmdPinPadReply, name, reply)
}
