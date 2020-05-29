package common

import "errors"

type EnumDevError uint16


// Device error codes
const (
	DevErrorSuccess EnumDevError = iota
	DevErrorGeneral
	DevErrorOutOfMemory
	DevErrorNullPointer
	DevErrorBadArgument
	DevErrorNotImplemented
	DevErrorNotInitialized
	DevErrorNotAccepted
	DevErrorNoAccess
	DevErrorCanceled
	DevErrorConfigFault
	DevErrorSystemFault
	DevErrorHardwareFault
	DevErrorSoftwareFault
	DevErrorDatabaseFault
	DevErrorNetworkFault
	DevErrorLinkerFault
	DevErrorLinkerTimeout
	DevErrorProtocolFault
	DevErrorSecurityFault
	DevErrorCommandFault
	DevErrorExecuteFault
	DevErrorWaitTimeout
	DevErrorPaperOut
	DevErrorPaperJam
	DevErrorNoCurrency
	DevErrorBillJammed
	DevErrorStackerFull
	DevErrorStackerEmpty
	DevErrorCantDispense
	DevErrorCassetteMiss
	DevErrorCounterFault
	DevErrorPickFault
	DevErrorBadCardData
	DevErrorBadKeyIndex
	DevErrorBadKeyValue
	DevErrorUnknown
)

// String returns a string explaining of the device error
func (e EnumDevError) String() string {
	switch e {
	case DevErrorSuccess:			return "Success"
	case DevErrorGeneral:			return "General error"
	case DevErrorOutOfMemory:		return "Out of memory"
	case DevErrorNullPointer:		return "Null pointer"
	case DevErrorBadArgument:		return "Bad argument"
	case DevErrorNotImplemented:	return "Not implemented"
	case DevErrorNotInitialized:	return "Not initialized"
	case DevErrorNotAccepted:		return "Not accepted"
	case DevErrorNoAccess:			return "No access"
	case DevErrorCanceled:			return "Canceled"
	case DevErrorConfigFault:		return "Config fault"
	case DevErrorSystemFault:		return "System fault"
	case DevErrorHardwareFault:		return "Hardware fault"
	case DevErrorSoftwareFault:		return "Software fault"
	case DevErrorDatabaseFault:		return "Database fault"
	case DevErrorNetworkFault:		return "Network fault"
	case DevErrorLinkerFault:		return "Linker fault"
	case DevErrorLinkerTimeout:		return "Linker timeout"
	case DevErrorProtocolFault:		return "Protocol fault"
	case DevErrorSecurityFault:		return "Security fault"
	case DevErrorCommandFault:		return "Command fault"
	case DevErrorExecuteFault:		return "Execute fault"
	case DevErrorWaitTimeout:		return "Wait timeout"
	case DevErrorPaperOut:			return "Paper out"
	case DevErrorPaperJam:			return "Paper jam"
	case DevErrorNoCurrency:		return "No currency"
	case DevErrorBillJammed:		return "Bill is jammed"
	case DevErrorStackerFull:		return "Stacker is full"
	case DevErrorStackerEmpty:		return "Stacker is empty"
	case DevErrorCantDispense:		return "Can't dispense"
	case DevErrorCassetteMiss:		return "Cassette mismatch"
	case DevErrorCounterFault:		return "Counter fault"
	case DevErrorPickFault:			return "Pick fault"
	case DevErrorBadCardData:		return "Bad card data"
	case DevErrorBadKeyIndex:		return "Bad key index"
	case DevErrorBadKeyValue:		return "Bad key value"
	case DevErrorUnknown:			return "Unknown error"
	default:						return "Other error"
	}
}

// Error is an implementation of error interface for device reply
type Error struct {
	code   EnumDevError
	reason error
}

func ExtendError(code EnumDevError, err error) error {
	if err == nil { return nil }
	out := &Error{
		code:   code,
		reason: err,
	}
	return out
}

func NewError(code EnumDevError, text string) error {
	out := &Error{
		code:   code,
		reason: nil,
	}
	if text != "" {
		out.reason = errors.New(text)
	}
	return out
}

// Code returns the code of device error
func (e *Error) Code() EnumDevError {
	if e == nil { return DevErrorSuccess }
	return e.code
}

// Error returns the full description of the device error
func (e *Error) Error() string {
	if e == nil { return DevErrorSuccess.String() }
	if e.reason != nil {
		return e.code.String() + ": " + e.reason.Error()
	}
	return e.code.String()
}

// CheckError returns code and description of error interface
func CheckError(err error) (EnumDevError, string) {
	if err == nil {
		return DevErrorSuccess, DevErrorSuccess.String()
	}
	if e, ok := err.(*Error); ok {
		return e.Code(), e.Error()
	}
	return DevErrorGeneral, err.Error()
}
