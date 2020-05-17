package generic

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
)

type BaseEngine struct {
	Log       *core.LogAgent
	DevName   string
	DevState  common.EnumDevState
	DevError  common.EnumDevError
	DevPrompt common.EnumDevPrompt
	DevAction common.EnumDevAction
	DevInform string
	DevReply  string
	CbDevice  common.DeviceCallback
}



func (be *BaseEngine) RunDeviceReply(cmd string) error {
	// StateChanged processing
	var err error
	reply := &common.DeviceReply{}
	reply.Command  = cmd
	reply.Action   = be.DevAction
	reply.DevState = be.DevState
	reply.ErrCode  = be.DevError
	reply.ErrText  = be.DevReply
	if be.CbDevice != nil {
		err = be.CbDevice.DeviceReply(be.DevName, reply)
	}
	if be.Log != nil {
		be.Log.Debug("DeviceReply: %s", reply.String())
	}
	return err
}

func (be *BaseEngine) RunStateChanged(state common.EnumDevState) error {
	// StateChanged processing
	var err error
	if be.DevState != state {
		query := &common.DeviceState{
			Action:   be.DevAction,
			OldState: be.DevState,
			NewState: state,
		}
		if be.CbDevice != nil {
			err = be.CbDevice.StateChanged(be.DevName, query)
		}
		if be.Log != nil {
			be.Log.Debug("StateChanged: %s", query.String())
		}
		be.DevState = state
	}
	return err
}

func (be *BaseEngine) RunExecuteError(errCode common.EnumDevError, reason string) error {
	be.DevError  = errCode
	be.DevReply  = common.NewError(errCode, reason).Error()
	// ExecuteError processing
	var err error
	if be.DevError != common.DevErrorSuccess {
		query := &common.DeviceError{
			Action:   be.DevAction,
			DevState: be.DevState,
			ErrCode:  be.DevError,
			ErrText:  be.DevReply,
		}
		if be.CbDevice != nil {
			err = be.CbDevice.ExecuteError(be.DevName, query)
		}
		if be.Log != nil {
			be.Log.Debug("ExecuteError: %s", query.String())
		}
	}
	return err
}

func (be *BaseEngine) RunActionPrompt(prompt common.EnumDevPrompt) error {
	be.DevPrompt = prompt
	// ActionPrompt processing
	var err error
	if be.DevPrompt != common.DevPromptNone {
		query := &common.DevicePrompt{
			Action: be.DevAction,
			Prompt: be.DevPrompt,
		}
		if be.CbDevice != nil {
			err = be.CbDevice.ActionPrompt(be.DevName, query)
		}
		if be.Log != nil {
			be.Log.Debug("ActionPrompt: %s", query.String())
		}
	}
	return err
}

func (be *BaseEngine) RunReaderReturn(inform string) error {
	be.DevInform = inform
	// ReaderReturn processing
	var err error
	if be.DevInform != "" {
		query := &common.DeviceInform{
			Action: be.DevAction,
			Inform: be.DevInform,
		}
		if be.CbDevice != nil {
			err = be.CbDevice.ReaderReturn(be.DevName, query)
		}
		if be.Log != nil {
			be.Log.Debug("ReaderReturn: %s", query.String())
		}
	}
	return err
}


