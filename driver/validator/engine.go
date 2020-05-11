package validator

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/config"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/driver/simulator"
	"time"
)

type ValidatorEngine struct {
	simulator.Simulator
//	BillIndex  byte
//	BillEnable uint32
//	BillSecure uint32
//	BillList   BillList
	Batch     common.ValidatorBatch
	SerialNo  string
	// internal data
	devName   string
	device    common.DeviceCallback
	validator common.ValidatorCallback
	booker    common.ValidatorBooker
	config    *config.DeviceConfig
//	protocol  *ValidatorProtocol
	log       *core.LogAgent
}

func (ve *ValidatorEngine) initEngine(cfg *config.DeviceConfig) *ValidatorEngine {
	ve.config = cfg
//	ve.protocol = GetValidatorProtocol(cfg.Linker)
	ve.log = core.GetLogAgent(core.LogLevelTrace, "Engine")
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
	ve.SetupMimic(valResetStages)
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
	ve.log.Debug("CashCode serial no is :%s", ve.SerialNo)
	return err
}

func (ve *ValidatorEngine) DevEnableBills(curr common.DevCurrency) error {
	ve.log.Debug("ValidatorEngine Set currency %d - %s", curr, curr.String())
	ve.SetupMimic(valWaitNoteStages)
	var err error
	return err
}

func (ve *ValidatorEngine) DevDisableBills() error {
	ve.SetupMimic(valStopWaitStages)
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
//	ve.log.Debug(ve.Batch.String())
	return err
}

func (ve *ValidatorEngine) DevNoteAccept() error {
	ve.SetupMimic(valStackingStages)
	var err error
	return err
}

func (ve *ValidatorEngine) DevNoteReturn() error {
	ve.SetupMimic(valReturningStages)
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


////////////////////////////////////////////////////////////////

func (ve *ValidatorEngine) NextMimicStage() {
	stage := ve.GetMimicStage()
	if stage != nil {
		ve.ProcessStage(stage, ve.device, ve.log)
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
}
func (ve *ValidatorEngine) StepAcceptingDone() {
}
func (ve *ValidatorEngine) StepEscrowedDone() {
}
func (ve *ValidatorEngine) StepStackingDone() {
}
func (ve *ValidatorEngine) StepReturningDone() {
}
func (ve *ValidatorEngine) StepRejectingDone() {
}

