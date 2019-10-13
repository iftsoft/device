package config

type EnumPaperPath  uint16
type EnumShowImage  uint16
type EnumSkipPrefix uint16
type EnumCardAccept uint16
type EnumCardAction uint16
type EnumBillAction uint16
type EnumCalcMethod uint16
type EnumOutputDir  uint16
type EnumUnitUsage  uint16

const (
	PaperPathDefault EnumPaperPath = iota
	PaperPathOne
	PaperPathTwo
)
func (e EnumPaperPath) String() string {
	switch e {
	case PaperPathDefault:		return "Default"
	case PaperPathOne:			return "Path 1"
	case PaperPathTwo:			return "Path 2"
	default:					return "Unknown"
	}
}

const (
	ShowImageNone EnumShowImage = iota
	ShowImageOnes
	ShowImagePage
)
func (e EnumShowImage) String() string {
	switch e {
	case ShowImageNone:			return "Don't show"
	case ShowImageOnes:			return "Only once"
	case ShowImagePage:			return "Every page"
	default:					return "Unknown"
	}
}

const (
	SkipPrefixNone EnumSkipPrefix = iota
	SkipPrefixBite
	SkipPrefixTwo
)
func (e EnumSkipPrefix) String() string {
	switch e {
	case SkipPrefixNone:		return "Don't skip"
	case SkipPrefixBite:		return "Skip 1 byte"
	case SkipPrefixTwo:			return "Skip 2 bytes"
	default:					return "Unknown"
	}
}

const (
	CardAcceptAnyCard EnumCardAccept = iota
	CardAcceptMagnetic
	CardAcceptSmart
)
func (e EnumCardAccept) String() string {
	switch e {
	case CardAcceptAnyCard:		return "Any card"
	case CardAcceptMagnetic:	return "Magnetic"
	case CardAcceptSmart:		return "Smart only"
	default:					return "Unknown"
	}
}

const (
	CardActionDefault EnumCardAction = iota	// Action by default
	CardActionReturn		// Return card to the client
	CardActionHoldOn		// Keep reading, hold the card
	CardActionCapture		// Capture the client card
)
func (e EnumCardAction) String() string {
	switch e {
	case CardActionDefault:		return "Default"
	case CardActionReturn:		return "Return card"
	case CardActionHoldOn:		return "Hod card on"
	case CardActionCapture:		return "Capture card"
	default:					return "Unknown"
	}
}

const (
	BillActionDefault EnumBillAction = iota	// Action by default
	BillActionStore			// Store accepted bill to the cassette
	BillActionReturn		// Return accepted bill to the client
)
func (e EnumBillAction) String() string {
	switch e {
	case BillActionDefault:		return "Default"
	case BillActionStore:		return "Store note"
	case BillActionReturn:		return "Return note"
	default:					return "Unknown"
	}
}

const (
	CalcMethodDefault EnumCalcMethod = iota	// By default
	CalcMethodMaximum		// Maximum nominal output
	CalcMethodUniform		// Uniform distribution notes
	CalcMethodMinimum		// Small nominal output
)
func (e EnumCalcMethod) EnumCalcMethod() string {
	switch e {
	case CalcMethodDefault:		return "Default"
	case CalcMethodMaximum:		return "Maximum nominal"
	case CalcMethodUniform:		return "Uniform distribution"
	case CalcMethodMinimum:		return "Minimum nominal"
	default:					return "Unknown"
	}
}

const (
	UnitIsAbsent EnumUnitUsage = iota
	UnitIsPresent
	UnitIsIgnored
)
func (e EnumUnitUsage) String() string {
	switch e {
	case UnitIsAbsent:			return "Unit is absent"
	case UnitIsPresent:			return "Unit is present"
	case UnitIsIgnored:			return "Unit is ignored"
	default:					return "Unknown"
	}
}



