package validator

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/driver/generic"
	"time"
)

type ValidatorEngine struct {
	generic.BaseValidator
	generic.Simulator
	booker      common.ValidatorBooker
	config      *config.DeviceConfig
	billIndex   int
}

func (ve *ValidatorEngine) initEngine(cfg *config.DeviceConfig) *ValidatorEngine {
	ve.config    = cfg
	ve.Log       = core.GetLogAgent(core.LogLevelTrace, "Engine")
	ve.billIndex = 0
	return ve
}

func (ve *ValidatorEngine) enterNote() {
	size := len(ve.Batch.Notes)
	if size <= 0 { return }
	for i:=0; i<size; i++ {
		index := ve.billIndex
		ve.billIndex ++
		if ve.billIndex >= size {
			ve.billIndex = 0
		}
		if ve.Accept.Currency == ve.Batch.Notes[index].Currency {
			ve.Accept.Nominal = ve.Batch.Notes[index].Nominal
			ve.Accept.Count   = 1
			ve.Accept.Amount  = ve.Accept.Nominal
			break
		}
	}
}

func (ve *ValidatorEngine) clearNote() {
	ve.Accept.Nominal = 0
	ve.Accept.Count   = 0
	ve.Accept.Amount  = 0
}


////////////////////////////////////////////////////////////////

func (ve *ValidatorEngine) DevStartup() error {
	var err error
	if ve.booker != nil {
		err = ve.booker.ReadNoteList(&ve.Batch)
	}
	ve.Accept.Currency = 980
	return err
}

func (ve *ValidatorEngine) DevCleanup() error {
	var err error
	return err
}

func (ve *ValidatorEngine) DevReset() error {
	ve.SetupMimic(valResetSteps)
	var err error
	for i := 0; i < 50; i++ {
		time.Sleep(200 * time.Millisecond)
		if ve.DevState == common.DevStateStandby {
			break
		}
	}
	return err
}

func (ve *ValidatorEngine) DevStatus() error {
	var err error
	return err
}

func (ve *ValidatorEngine) DevEnableBills(curr common.DevCurrency) error {
	ve.Log.Debug("ValidatorEngine Set currency %d - %s", curr, curr.String())
	ve.SetupMimic(valWaitNoteSteps)
	var err error
	return err
}

func (ve *ValidatorEngine) DevDisableBills() error {
	ve.SetupMimic(valStopWaitSteps)
	var err error
	return err
}

func (ve *ValidatorEngine) DevInitBillList() error {
	var err error
	if ve.booker != nil {
		err = ve.booker.CloseBatch(&ve.Batch)
		if err == nil {
			err = ve.booker.InitNoteList(valNoteListUah)
		}
		if err == nil {
			err = ve.booker.ReadNoteList(&ve.Batch)
		}
	}
	ve.Accept.Currency = 980
//	ve.Log.Debug(ve.Batch.String())
	return err
}

func (ve *ValidatorEngine) DevNoteAccept() error {
	ve.SetupMimic(valStackingSteps)
	var err error
	if ve.booker != nil {
		err = ve.booker.DepositNote(0, &ve.Accept)
	}
	return err
}

func (ve *ValidatorEngine) DevNoteReturn() error {
	ve.SetupMimic(valReturningSteps)
	var err error
	return err
}

func (ve *ValidatorEngine) DevCheckBatch() error {
	var err error
	if ve.booker != nil {
		err = ve.booker.ReadNoteList(&ve.Batch)
	}
	return err
}

func (ve *ValidatorEngine) DevClearBatch() error {
	var err error
	if ve.booker != nil {
		err = ve.booker.CloseBatch(&ve.Batch)
	}
	return err
}



////////////////////////////////////////////////////////////////

func (ve *ValidatorEngine) NextMimicStage() {
	stage := ve.GetMimicStep()
	if stage != nil {
		_ = ve.ProcessStage(&ve.BaseEngine, stage)
		switch stage.Value {
		case StepWaitNoteDone:
			ve.StepWaitNoteDone()
		case StepAcceptingDone:
			ve.StepAcceptingDone()
		case StepEscrowedDone:
			ve.StepEscrowedDone()
		case StepStackingDone:
			ve.StepStackingDone()
		case StepReturningDone:
			ve.StepReturningDone()
		case StepRejectingDone:
			ve.StepRejectingDone()
		default:
		}
	}
	return
}

func (ve *ValidatorEngine) StepWaitNoteDone() {
	ve.SetupMimic(valAcceptingSteps)
}
func (ve *ValidatorEngine) StepAcceptingDone() {
	ve.SetupMimic(valEscrowedSteps)
	ve.enterNote()
	_ = ve.RunNoteAccepted(&ve.Accept)
}
func (ve *ValidatorEngine) StepEscrowedDone() {
	ve.SetupMimic(valRejectingSteps)
	ve.clearNote()
}
func (ve *ValidatorEngine) StepStackingDone() {
	ve.SetupMimic(valWaitNoteSteps)
	_ = ve.RunCashIsStored(&ve.Accept)
	ve.clearNote()
}
func (ve *ValidatorEngine) StepReturningDone() {
	ve.SetupMimic(valWaitNoteSteps)
	_ = ve.RunCashReturned(&ve.Accept)
	ve.clearNote()
}
func (ve *ValidatorEngine) StepRejectingDone() {
	ve.SetupMimic(valWaitNoteSteps)
}


////////////////////////////////////////////////////////////////
