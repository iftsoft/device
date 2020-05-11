package validator

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/driver/simulator"
)

const (
	StepDoNothing int = iota
	StepWaitNoteDone
	StepAcceptingDone
	StepEscrowedDone
	StepStackingDone
	StepReturningDone
	StepRejectingDone
)

var valResetStages = simulator.MimicStages{
	{ 0, 0, common.DevStateWorking, 0, common.DevPromptUnitWork, common.DevActionInitialization, "", "" },
	{ 10, 0, common.DevStateReady, 0, 0, common.DevActionInitialization, "", "" },
	{ 1, 0, common.DevStateStandby, 0, 0, common.DevActionDoNothing, "", "" },
}

var valWaitNoteStages = simulator.MimicStages{
	{ 0, 0, common.DevStateWorking, 0, common.DevPromptCashInsertBill, common.DevActionNoteWaiting, "", "" },
	{ 1, 0, common.DevStateWaiting, 0, 0, common.DevActionNoteWaiting, "", "" },
	{ 50, StepWaitNoteDone, common.DevStateWaiting, 0, 0, 0, "", "" },
}

var valAcceptingStages = simulator.MimicStages{
	{ 0, 0, common.DevStateCashAccepting, 0, common.DevPromptCashAccepting, common.DevActionNoteAccepting, "", "" },
	{ 1, 0, common.DevStateCashAccepting, 0, 0, common.DevActionNoteAccepting, "", "" },
	{ 10, StepAcceptingDone, common.DevStateCashAccepting, 0, 0, 0, "", "" },
}

var valEscrowedStages = simulator.MimicStages{
	{ 0, 0, common.DevStateCashEscrowed, 0, common.DevPromptCashEscrowed, common.DevActionNoteAccepting, "", "" },
	{ 1, 0, common.DevStateCashEscrowed, 0, 0, common.DevActionNoteAccepting, "", "" },
	{ 30, StepEscrowedDone, common.DevStateCashEscrowed, 0, 0, 0, "", "" },
}

var valStackingStages = simulator.MimicStages{
	{ 0, 0, common.DevStateCashStacking, 0, common.DevPromptCashStacking, common.DevActionNoteStacking, "", "" },
	{ 1, 0, common.DevStateCashStacking, 0, 0, common.DevActionNoteStacking, "", "" },
	{ 20, StepStackingDone, common.DevStateCashStacked, 0, 0, 0, "", "" },
}

var valReturningStages = simulator.MimicStages{
	{ 0, 0, common.DevStateCashReturning, 0, common.DevPromptCashReturning, common.DevActionNoteReturning, "", "" },
	{ 1, 0, common.DevStateCashReturning, 0, 0, common.DevActionNoteReturning, "", "" },
	{ 10, StepReturningDone, common.DevStateCashReturned, 0, 0, 0, "", "" },
}

var valRejectingStages = simulator.MimicStages{
	{ 0, 0, common.DevStateCashRejecting, 0, common.DevPromptCashReturning, common.DevActionNoteRejecting, "", "" },
	{ 1, 0, common.DevStateCashRejecting, 0, 0, common.DevActionNoteRejecting, "", "" },
	{ 10, StepRejectingDone, common.DevStateCashRejecting, 0, 0, 0, "", "" },
}

var valStopWaitStages = simulator.MimicStages{
	{ 0, 0, common.DevStateWorking, 0, common.DevPromptUnitDone, common.DevActionDeviceStopping, "", "" },
	{ 5, 0, common.DevStateReady, 0, 0, common.DevActionDeviceStopping, "", "" },
	{ 1, 0, common.DevStateStandby, 0, 0, common.DevActionDoNothing, "", "" },
}

var valNoteListUah = common.ValidNoteList {
	{"", 980, 0, 1.0, 0.0, },
	{"", 980, 0, 2.0, 0.0, },
	{"", 980, 0, 5.0, 0.0, },
	{"", 980, 0, 10.0, 0.0, },
	{"", 980, 0, 20.0, 0.0, },
	{"", 980, 0, 50.0, 0.0, },
	{"", 980, 0, 100.0, 0.0, },
	{"", 980, 0, 200.0, 0.0, },
	{"", 980, 0, 500.0, 0.0, },
	{"", 980, 0, 1000.0, 0.0, },
}

