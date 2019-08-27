package common

type DevAmount float32
type DevCounter int32
type DevCurrency int16
type EnumDevType int16
type EnumDevError int16
type EnumDevState int16
type EnumDevAction int16
type EnumDevPrompt int16
type EnumDevWarning int16
type EnumSystemState int16

// System state codes
const (
	SysStateUndefined EnumSystemState = iota
	SysStateRunning
	SysStateStopped
)

// Device types
const (
	DevTypeDefault EnumDevType = iota
	DevTypeCardReader
	DevTypeBarScanner
	DevTypePINEntry
	DevTypePrinter
	DevTypeCashValidator
	DevTypeCoinValidator
	DevTypeCashDispenser
	DevTypeCoinDispenser
	DevTypeVending
	DevTypeCustom
)

// Device error codes
const (
	DevErrorSuccess EnumDevError = iota
	DevErrorBadArgument
)

// Device status codes
const (
	DevStateUndefined EnumDevState = iota
	DevStateReady
	DevStateWorking
	DevStateWaiting
	DevStateStandby
	DevStateBadState
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
	DevState
)

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
	DevPrompt
)

// Device action codes
const (
	DevActionNothing EnumDevAction = iota
	DevAction
)

/*
// Device error codes
enum HS_DeviceErrors {
	HS_OBJ_ERR_SUCCESS = 0,
	HS_OBJ_ERR_OUT_OF_MEMORY,
	HS_OBJ_ERR_BAD_ARGUMENT,
	HS_OBJ_ERR_BAD_HANDLER,
	HS_OBJ_ERR_NOT_IMPLEMENT,
	HS_OBJ_ERR_NOT_LICENSED,
	HS_OBJ_ERR_NOT_FOUND,
	HS_OBJ_ERR_NOT_LOADED,
	HS_OBJ_ERR_NOT_INITED,
	HS_OBJ_ERR_IS_DISABLED,
	HS_OBJ_ERR_NO_ACCESS,
	HS_OBJ_ERR_CANCELED,
	HS_OBJ_ERR_MODULE_FAULT,
	HS_OBJ_ERR_THREAD_FAULT,
	HS_OBJ_ERR_SYSTEM_FAULT,
	HS_OBJ_ERR_NOT_LINKED ,
	HS_COM_ERR_OPEN_FAULT ,
	HS_COM_ERR_LINK_TIMEOUT,
	HS_COM_ERR_WRITE_FAULT,
	HS_COM_ERR_READ_FAULT,
	HS_DEV_ERR_NOT_INSTALLED ,
	HS_DEV_ERR_NAME_NOT_SET ,
	HS_DEV_ERR_PORT_NOT_SET ,
	HS_DEV_ERR_STARTUP_FAULT ,
	HS_DEV_ERR_WAIT_TIMEOUT,
	HS_DEV_ERR_CMD_EXCEPTION,
	HS_DEV_ERR_CHECKSUM_FAULT,
	HS_DEV_ERR_PROTOCOL_FAULT,
	HS_DEV_ERR_EXECUTE_FAULT,
	HS_DEV_ERR_BAD_STATE,
	HS_CRD_ERR_BAD_CARD_DATA,
	HS_PIN_ERR_BAD_KEY_INDEX,
	HS_PIN_ERR_BAD_TRANSPORT_KEY,
	HS_PIN_ERR_NO_TRANSPORT_KEY,
	HS_PIN_ERR_NO_MASTER_KEY,
	HS_PIN_ERR_NO_WORK_KEY,
	HS_PIN_ERR_NO_CARD_PAN,
	HS_CSH_ERR_RESERVE,
	HS_CSH_ERR_NO_CURRENCY,
	HS_CSH_ERR_SHUTTER_FAULT,
	HS_CSH_ERR_SHUTTER_CLOSED,
	HS_CSH_ERR_CHECKER_FAULT,
	HS_CSH_ERR_BILL_JAMMED,
	HS_CSH_ERR_STACKER_FULL,
	HS_PRN_ERR_PAPER_OUT,
	HS_PRN_ERR_PAPER_JAM,
	HS_PRN_ERR_HARD_ERROR,
	HS_PRN_ERR_SOFT_ERROR,
	HS_DSP_ERR_EMPTY_STACK,
	HS_DSP_ERR_CANT_DISPENSE,
	HS_DSP_ERR_CASSETTE_MISS,
	HS_DSP_ERR_COUNTER_FAULT,
	HS_DSP_ERR_PICK_FAULT,
	HS_DSP_ERR_MOTION_FAULT,
	HS_OBJ_ERR_UNKNOWN
};

// Action event codes
enum HS_ActionEvents {
	HS_ACTION_NOTHING = 0,
	HS_ACTION_DEVICE_INIT,		//  1
	HS_ACTION_DEVICE_RESTART,	//  2
	HS_ACTION_DEVICE_START,		//  3
	HS_ACTION_DEVICE_STOP,		//  4
	HS_ACTION_DEVICE_ENABLE,	//  5
	HS_ACTION_DEVICE_DISABLE,	//  6
	HS_ACTION_DEVICE_RESET,		//  7
	HS_ACTION_DEVICE_STATUS,	//  8
	HS_ACTION_KEYBOARD_START,	//  9
	HS_ACTION_KEYBOARD_DONE,	// 10
	HS_ACTION_KEYBOARD_ABORT,	// 11
	HS_ACTION_CARD_READ_START,	// 12
	HS_ACTION_CARD_READ_DONE,	// 13
	HS_ACTION_CARD_READ_ABORT,	// 14
	HS_ACTION_SCANING_START,	// 15
	HS_ACTION_SCANING_DONE,		// 16
	HS_ACTION_SCANING_ABORT,	// 17
	HS_ACTION_TEXT_PRINT_START,	// 18
	HS_ACTION_TEXT_PRINT_DONE,	// 19
	HS_ACTION_PIN_ENTER_START,	// 20
	HS_ACTION_PIN_ENTER_DONE,	// 21
	HS_ACTION_PIN_ENTER_ABORT,	// 22
	HS_ACTION_DATA_ENTER_START,	// 23
	HS_ACTION_DATA_ENTER_DONE,	// 24
	HS_ACTION_DATA_ENTER_ABORT,	// 25
	HS_ACTION_WORK_KEY_SET,		// 26
	HS_ACTION_WORK_KEY_TEST,	// 27
	HS_ACTION_MASTER_KEY_SET,	// 28
	HS_ACTION_MASTER_KEY_TEST,	// 29
	HS_ACTION_CURRENCY_SET,		// 30
	HS_ACTION_ACCEPTING_START,	// 31
	HS_ACTION_ACCEPT_NOTE,		// 32
	HS_ACTION_ACCEPTING_STOP,	// 33
	HS_ACTION_NOTE_RETURNED,	// 34
	HS_ACTION_NOTE_STORED,		// 35
	HS_ACTION_STORE_CHECKED,	// 36
	HS_ACTION_STORE_CLEARED,	// 37
	HS_ACTION_CHIP_GET_ATR,		// 38
	HS_ACTION_CHIP_POWER_OFF,	// 39
	HS_ACTION_CHIP_DO_COMMAND,	// 40
	HS_ACTION_CHIP_READ_DATA,	// 41
	HS_ACTION_CHIP_WRITE_DATA,	// 42
	HS_ACTION_TURN_READ_INPUT,	// 43
	HS_ACTION_TURN_SWITCH_LED,	// 44
	HS_ACTION_TURN_TURN_RELAY,	// 45
	HS_ACTION_TURN_WATCH_DOG,	// 46
	HS_ACTION_SETUP_CALLBACK,	// 47
	HS_ACTION_CLEAR_CALLBACK,	// 48
	HS_ACTION_DISP_VEND_ITEM,	// 49
	HS_ACTION_CARD_IS_CAPTURED,	// 50
	HS_ACTION_CARD_ENTER_START,	// 51
	HS_ACTION_CARD_ENTER_STOP,	// 52
	HS_ACTION_CARD_EJECT_START,	// 53
	HS_ACTION_CARD_EJECT_STOP,	// 54
	HS_ACTION_STACK_CHECKED,	// 55
	HS_ACTION_STACK_ASSIGNED,	// 56
	HS_ACTION_NOTE_COUNT_START,	// 57
	HS_ACTION_NOTE_COUNT_STOP,	// 58
	HS_ACTION_NOTE_OUTPUT_START,// 59
	HS_ACTION_NOTE_OUTPUT_STOP,	// 60
	HS_ACTION_NOTE_DIVERT_START,// 61
	HS_ACTION_NOTE_DIVERT_STOP,	// 62
	HS_ACTION_COUNTS_CHECKED,	// 63
	HS_ACTION_COUNTS_SETUPED,	// 64
	HS_ACTION_NOTE_ASSIGNED,	// 65
	HS_ACTION_PED_CREATE_MAC,	// 66
	HS_ACTION_PED_ENCRYPT_DATA,	// 67
	HS_ACTION_PED_DECRYPT_DATA,	// 68

	HS_ACTION_LAST_COMMAND,
	HS_ACTION_USER_CARE = 128,
	HS_ACTION_DOOR_IS_OPENED,
	HS_ACTION_DOOR_IS_CLOSED,
	HS_ACTION_CARD_IN_GATE,
	HS_ACTION_CARD_IS_OUT,
	HS_ACTION_ITEM_INSIDE,
	HS_ACTION_
};


*/
