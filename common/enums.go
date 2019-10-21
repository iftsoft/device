package common

import (
	"strings"
)

type EnumDevState uint16
type EnumDevAction uint16
type EnumDevPrompt uint16
type DevScopeMask uint64
type DevTypeMask uint64

// Scope Flags
const (
	ScopeFlagSystem  DevScopeMask = 1 << iota
	ScopeFlagDevice
	ScopeFlagPrinter
	ScopeFlagReader
	ScopeFlagPinPad
	ScopeFlagValidator
	ScopeFlagDispenser
	ScopeFlagVending
	ScopeFlagCustom
	ScopeFlagUnknown = 0
)

// Device types
const (
	DevTypePrinter DevTypeMask = 1 << iota
	DevTypeCardReader
	DevTypeBarScanner
	DevTypeCashValidator
	DevTypeCoinValidator
	DevTypeCashDispenser
	DevTypeCoinDispenser
	DevTypeVending
	DevTypePINEntry
	DevTypeCustom
	DevTypeUndefined = 0
	DevTypeCommon    = 0x0FF
)

func (e DevTypeMask) ToString() string {
	if e == DevTypeUndefined {	return "Undefined"	}
	list := make([]string, 0)
	if (e & DevTypePrinter) == DevTypePrinter				{	list = append(list, "Printer")	}
	if (e & DevTypeCardReader) == DevTypeCardReader			{	list = append(list, "CardReader")	}
	if (e & DevTypeBarScanner) == DevTypeBarScanner			{	list = append(list, "BarScanner")	}
	if (e & DevTypeCashValidator) == DevTypeCashValidator	{	list = append(list, "CashValidator")	}
	if (e & DevTypeCoinValidator) == DevTypeCoinValidator	{	list = append(list, "CoinValidator")	}
	if (e & DevTypeCashDispenser) == DevTypeCashDispenser	{	list = append(list, "CashDispenser")	}
	if (e & DevTypeCoinDispenser) == DevTypeCoinDispenser	{	list = append(list, "CoinDispenser")	}
	if (e & DevTypeVending) == DevTypeVending				{	list = append(list, "Vending")	}
	if (e & DevTypePINEntry) == DevTypePINEntry				{	list = append(list, "PINEntry")	}
	if (e & DevTypeCustom) == DevTypeCustom					{	list = append(list, "Custom")	}
	return strings.Join(list,",")
}

// Device status codes
const (
	DevStateUndefined EnumDevState = iota
	DevStateReady
	DevStateWorking
	DevStateWaiting
	DevStateStandby
	DevStateOffLine
	DevStateFailure
	DevStateHardError
	DevStateSoftError
	DevStatePrnTonerOut
	DevStatePrnPaperOut
	DevStatePrnPaperJam
	DevStatePrnCoverOpen
	DevStatePrnOutputBin
	DevStateCardInFront
	DevStateCardInside
	DevStateCardInTrack
	DevStateCardPowered
	DevStateCashAccepting
	DevStateCashEscrowed
	DevStateCashStacking
	DevStateCashStacked
	DevStateCashReturning
	DevStateCashReturned
	DevStateCashRejecting
	DevStateCashBillJammed
	DevStateCashStackerFull
	DevStateDispCapturing
	DevStateDispDispensing
	DevStateDispDispensed
	DevStateDispEmptyStack
	DevStateDoorBroken
)

// String returns a string explaining the device status
func (e EnumDevState) String() string {
	switch e {
	case DevStateUndefined:			return "Undefined"
	case DevStateReady:				return "Ready"
	case DevStateWorking:			return "Working"
	case DevStateWaiting:			return "Waiting"
	case DevStateStandby:			return "Standby"
	case DevStateOffLine:			return "Off line"
	case DevStateFailure:			return "Failure"
	case DevStateHardError:			return "Hardware error"
	case DevStateSoftError:			return "Software error"
	case DevStatePrnTonerOut:		return "Toner out"
	case DevStatePrnPaperOut:		return "Paper out"
	case DevStatePrnPaperJam:		return "Paper jam"
	case DevStatePrnCoverOpen:		return "Cover open"
	case DevStatePrnOutputBin:		return "Output bin"
	case DevStateCardInFront:		return "Card in front"
	case DevStateCardInside:		return "Card inside"
	case DevStateCardInTrack:		return "Card in track"
	case DevStateCardPowered:		return "Card powered"
	case DevStateCashAccepting:		return "Accepting note"
	case DevStateCashEscrowed:		return "Note is escrowed"
	case DevStateCashStacking:		return "Stacking note"
	case DevStateCashStacked:		return "Note is stacked"
	case DevStateCashReturning:		return "Returning note"
	case DevStateCashReturned:		return "Note is returned"
	case DevStateCashRejecting:		return "Rejecting note"
	case DevStateCashBillJammed:	return "Note is jammed"
	case DevStateCashStackerFull:	return "Stacker is full"
	case DevStateDispCapturing:		return "Capturing"
	case DevStateDispDispensing:	return "Dispensing"
	case DevStateDispDispensed:		return "Dispensed"
	case DevStateDispEmptyStack:	return "Stack is empty"
	case DevStateDoorBroken:		return "Door is broken"
	default:						return "Unknown"
	}
}

// Device prompt codes
const (
	DevPromptNone EnumDevPrompt = iota
	DevPromptUnitWait
	DevPromptUnitDone
	DevPromptUnitError
	DevPromptCardSwipe
	DevPromptCardInsert
	DevPromptCardRemove
	DevPromptCardCapture
	DevPromptCardFailure
	DevPromptScanBarcode
	DevPromptPrintText
	DevPromptPpadEntryData
	DevPromptPpadEntryPIN
	DevPromptCashInsertBill
	DevPromptCashAccepting
	DevPromptCashEscrowed
	DevPromptCashStacking
	DevPromptCashReturning
	DevPromptCashFailure
	DevPromptCashBillJammed
	DevPromptCashStackerFull
	DevPromptDispTakeItem
	DevPromptDispTakeCard
	DevPromptDispTakeBill
	DevPromptDispTakeCoin
	DevPromptDispCapture
	DevPromptDispFailure
)

// String returns a string explaining the device prompt
func (e EnumDevPrompt) String() string {
	switch e {
	case DevPromptNone:				return "Thank you."
	case DevPromptUnitWait:			return "Please, wait while device is working."
	case DevPromptUnitDone:			return "All done. Thank you."
	case DevPromptUnitError:		return "Device error has occurred."
	case DevPromptCardSwipe:		return "Swipe your card through card reader"
	case DevPromptCardInsert:		return "Insert your card into card reader"
	case DevPromptCardRemove:		return "Remove your card from card reader"
	case DevPromptCardCapture:		return "Your card is captured. Call to support."
	case DevPromptCardFailure:		return "Device can't read your card."
	case DevPromptScanBarcode:		return "Scan barcode by the scanner."
	case DevPromptPrintText:		return "Your receipt is printing."
	case DevPromptPpadEntryData:	return "Enter your number."
	case DevPromptPpadEntryPIN:		return "Enter your PIN code."
	case DevPromptCashInsertBill:	return "Insert your banknote."
	case DevPromptCashAccepting:	return "Wait while your note is being accepted."
	case DevPromptCashEscrowed:		return "Wait while your note is being processed."
	case DevPromptCashStacking:		return "Wait while your note is being stacked."
	case DevPromptCashReturning:	return "Wait while your note is being returned."
	case DevPromptCashFailure:		return "Validator error has occurred."
	case DevPromptCashBillJammed:	return "Your note is jammed. Call to support."
	case DevPromptCashStackerFull:	return "Out of service. Stacker is full."
	case DevPromptDispTakeItem:		return "Please, take your item"
	case DevPromptDispTakeCard:		return "Please, take your card"
	case DevPromptDispTakeBill:		return "Please, take your note"
	case DevPromptDispTakeCoin:		return "Please, take your coin"
	case DevPromptDispCapture:		return "Your item has been captured."
	case DevPromptDispFailure:		return "Dispenser error has occurred."
	default:						return "Unknown"
	}
}

// Device action codes
const (
	DevActionDoNothing EnumDevAction = iota
	DevActionInitialization
	DevActionReconciliation
	DevActionDeviceStarting
	DevActionDeviceStopping
	DevActionDeviceResetting
	DevActionCardEntering
	DevActionCardReading
	DevActionCardEjecting
	DevActionCardCapturing
	DevActionCardProcessing
	DevActionBarScanning
	DevActionTextPrinting
	DevActionDataEntering
	DevActionKeyEntering
	DevActionPinEntering
	DevActionNoteWaiting
	DevActionNoteAccepting
	DevActionNoteStacking
	DevActionNoteReturning
	DevActionNoteRejecting
	DevActionNotePicking
	DevActionNoteDispensing
	DevActionNoteDiverting
	DevActionLightSwitching
	DevActionRelaySwitching
	DevActionSensorChecking
	DevActionItemVending
	DevAction
)

// String returns a string explaining the device action
func (e EnumDevAction) String() string {
	switch e {
	case DevActionDoNothing:		return "Do nothing"
	case DevActionInitialization:	return "Initialization"
	case DevActionReconciliation:	return "reconciliation"
	case DevActionDeviceStarting:	return "Device starting"
	case DevActionDeviceStopping:	return "Device stopping"
	case DevActionDeviceResetting:	return "Device resetting"
	case DevActionCardEntering:		return "Card entering"
	case DevActionCardReading:		return "Card reading"
	case DevActionCardEjecting:		return "Card ejecting"
	case DevActionCardCapturing:	return "Card capturing"
	case DevActionCardProcessing:	return "Card processing"
	case DevActionBarScanning:		return "Barcode scanning"
	case DevActionTextPrinting:		return "Text printing"
	case DevActionDataEntering:		return "Data entering"
	case DevActionKeyEntering:		return "Key entering"
	case DevActionPinEntering:		return "PIN entering"
	case DevActionNoteWaiting:		return "Waiting for note"
	case DevActionNoteAccepting:	return "Note accepting"
	case DevActionNoteStacking:		return "Note stacking"
	case DevActionNoteReturning:	return "Note returning"
	case DevActionNoteRejecting:	return "Note rejecting"
	case DevActionNotePicking:		return "Note picking"
	case DevActionNoteDispensing:	return "Note dispensing"
	case DevActionNoteDiverting:	return "Note diverting"
	case DevActionLightSwitching:	return "Light switching"
	case DevActionRelaySwitching:	return "Relay switching"
	case DevActionSensorChecking:	return "Sensor checking"
	case DevActionItemVending:		return "Item vending"
	case DevAction:					return "Action"
	default:						return "Unknown"
	}
}

