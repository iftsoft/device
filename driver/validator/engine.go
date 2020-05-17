package validator

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/driver/generic"
	"time"
)

type ValidatorEngine struct {
	generic.Simulator
	Batch       common.ValidatorBatch
	Accept      common.ValidatorAccept
	SerialNo    string
	// internal data
	CbValidator common.ValidatorCallback
	booker      common.ValidatorBooker
	config      *config.DeviceConfig
}

func (ve *ValidatorEngine) initEngine(cfg *config.DeviceConfig) *ValidatorEngine {
	ve.config = cfg
//	ve.protocol = GetValidatorProtocol(cfg.Linker)
	ve.Log = core.GetLogAgent(core.LogLevelTrace, "Engine")
	return ve
}

////////////////////////////////////////////////////////////////

func (ve *ValidatorEngine) DevStartup() error {
	var err error
	if ve.booker != nil {
		err = ve.booker.ReadNoteList(&ve.Batch)
	}
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

func (ve *ValidatorEngine) DevIdent() error {
	var err error
	ve.SerialNo = "Emul01234567"
	ve.Log.Debug("CashCode serial no is :%s", ve.SerialNo)
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
		ve.ProcessStage(stage)
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
	ve.Accept.Currency = 980
	ve.Accept.Nominal  = 2.0
	ve.Accept.Count    = 1
	ve.Accept.Amount   = 2.0
	_ = ve.RunNoteAccepted(&ve.Accept)
}
func (ve *ValidatorEngine) StepEscrowedDone() {
	ve.SetupMimic(valRejectingSteps)
}
func (ve *ValidatorEngine) StepStackingDone() {
	ve.SetupMimic(valWaitNoteSteps)
	_ = ve.RunCashIsStored(&ve.Accept)
}
func (ve *ValidatorEngine) StepReturningDone() {
	ve.SetupMimic(valWaitNoteSteps)
	_ = ve.RunCashReturned(&ve.Accept)
}
func (ve *ValidatorEngine) StepRejectingDone() {
	ve.SetupMimic(valWaitNoteSteps)
}


////////////////////////////////////////////////////////////////


func (ve *ValidatorEngine) RunValidatorStore(cmd string) error {
	var err error
	reply := &common.ValidatorStore{}
	reply.Command  = cmd
	reply.DevState = ve.DevState
	reply.ErrCode  = ve.DevError
	reply.ErrText  = ve.DevReply
	reply.Notes    = ve.Batch.Notes
	reply.BatchId  = ve.Batch.BatchId
	reply.State    = ve.Batch.State
	reply.Detail   = ve.Batch.Detail
	if ve.CbValidator != nil {
		err = ve.CbValidator.ValidatorStore(ve.DevName, reply)
	}
	if ve.Log != nil {
		ve.Log.Debug("ValidatorStore: %s", reply.String())
	}
	return err
}

func (ve *ValidatorEngine) RunNoteAccepted(value *common.ValidatorAccept) error {
	var err error
	if ve.CbValidator != nil {
		err = ve.CbValidator.NoteAccepted(ve.DevName, value)
	}
	if ve.Log != nil {
		ve.Log.Debug("NoteAccepted: %s", value.String())
	}
	return err
}

func (ve *ValidatorEngine) RunCashIsStored(value *common.ValidatorAccept) error {
	var err error
	if ve.CbValidator != nil {
		err = ve.CbValidator.CashIsStored(ve.DevName, value)
	}
	if ve.Log != nil {
		ve.Log.Debug("CashIsStored: %s", value.String())
	}
	return err
}

func (ve *ValidatorEngine) RunCashReturned(value *common.ValidatorAccept) error {
	var err error
	if ve.CbValidator != nil {
		err = ve.CbValidator.CashReturned(ve.DevName, value)
	}
	if ve.Log != nil {
		ve.Log.Debug("CashReturned: %s", value.String())
	}
	return err
}
