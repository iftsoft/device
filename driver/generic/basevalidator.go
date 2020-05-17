package generic

import (
	"github.com/iftsoft/device/common"
)

type BaseValidator struct {
	BaseEngine
	Batch       common.ValidatorBatch
	Accept      common.ValidatorAccept
	CbValidator common.ValidatorCallback
}



func (bv *BaseValidator) RunValidatorStore(cmd string) error {
	var err error
	reply := &common.ValidatorStore{}
	reply.Command  = cmd
	reply.DevState = bv.DevState
	reply.ErrCode  = bv.DevError
	reply.ErrText  = bv.DevReply
	reply.Notes    = bv.Batch.Notes
	reply.BatchId  = bv.Batch.BatchId
	reply.State    = bv.Batch.State
	reply.Detail   = bv.Batch.Detail
	if bv.CbValidator != nil {
		err = bv.CbValidator.ValidatorStore(bv.DevName, reply)
	}
	if bv.Log != nil {
		bv.Log.Debug("Callback ValidatorStore: %s", reply.String())
	}
	return err
}

func (bv *BaseValidator) RunNoteAccepted(value *common.ValidatorAccept) error {
	var err error
	if bv.CbValidator != nil {
		err = bv.CbValidator.NoteAccepted(bv.DevName, value)
	}
	if bv.Log != nil {
		bv.Log.Debug("Callback NoteAccepted: %s", value.String())
	}
	return err
}

func (bv *BaseValidator) RunCashIsStored(value *common.ValidatorAccept) error {
	var err error
	if bv.CbValidator != nil {
		err = bv.CbValidator.CashIsStored(bv.DevName, value)
	}
	if bv.Log != nil {
		bv.Log.Debug("Callback CashIsStored: %s", value.String())
	}
	return err
}

func (bv *BaseValidator) RunCashReturned(value *common.ValidatorAccept) error {
	var err error
	if bv.CbValidator != nil {
		err = bv.CbValidator.CashReturned(bv.DevName, value)
	}
	if bv.Log != nil {
		bv.Log.Debug("Callback CashReturned: %s", value.String())
	}
	return err
}

