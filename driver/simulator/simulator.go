package simulator

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"sync"
)

////////////////////////////////////////////////////////////////

type MimicStage struct {
	Delay  int
	Value  int
	State  common.EnumDevState
	Error  common.EnumDevError
	Prompt common.EnumDevPrompt
	Action common.EnumDevAction
	Inform string
	Reason string
}

type MimicStages []MimicStage

////////////////////////////////////////////////////////////////

type Mimicker struct {
	stages     MimicStages
	mutex      sync.Mutex
	index      int
	delay      int
}

func (m *Mimicker) SetStages(list MimicStages) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.stages = list
	m.index = 0
	m.delay = 0
}

func (m *Mimicker) GetStage() *MimicStage {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.index >= len(m.stages) {
		return nil
	}
	stage := &m.stages[m.index]
	if m.delay >= stage.Delay {
		m.delay = 0
		m.index++
		return stage
	}
	m.delay++
	return nil
}

////////////////////////////////////////////////////////////////

type Simulator struct {
	DevState  common.EnumDevState
	DevError  common.EnumDevError
	DevPrompt common.EnumDevPrompt
	DevAction common.EnumDevAction
	DevInform string
	DevReply  string
	// internal data
	mimicker  Mimicker
}

func (sm *Simulator) SetupMimic(list MimicStages) {
	sm.mimicker.SetStages(list)
}

func (sm *Simulator) ClearMimic() {
	sm.mimicker.SetStages(nil)
}

func (sm *Simulator) GetMimicStage() *MimicStage {
	return sm.mimicker.GetStage()
}

func (sm *Simulator) ClearDevice() {
	sm.DevState  = common.DevStateUndefined
	sm.DevError  = common.DevErrorSuccess
	sm.DevPrompt = common.DevPromptNone
	sm.DevAction = common.DevActionDoNothing
	sm.DevInform = ""
	sm.DevReply  = ""
}

func (sm *Simulator) ProcessStage(stage *MimicStage, callback common.DeviceCallback, log *core.LogAgent) {
	if stage == nil {
		return
	}
	sm.DevAction = stage.Action
	sm.RunStateChanged(stage.State, callback, log)
	sm.RunExecuteError(stage.Error, stage.Reason, callback, log)
	sm.RunActionPrompt(stage.Prompt, callback, log)
	sm.RunReaderReturn(stage.Inform, callback, log)
}

func (sm *Simulator) RunStateChanged(state common.EnumDevState, callback common.DeviceCallback, log *core.LogAgent) {
	// StateChanged processing
	if sm.DevState != state {
		query := &common.DeviceState{
			Action:   sm.DevAction,
			OldState: sm.DevState,
			NewState: state,
		}
		if callback != nil {
			_ = callback.StateChanged("", query)
		}
		if log != nil {
			log.Debug("StateChanged: Action - %s, Old state - %s, New state - %s",
				sm.DevAction.String(), sm.DevState.String(), state.String())
		}
		sm.DevState = state
	}
}

func (sm *Simulator) RunExecuteError(errCode common.EnumDevError, reason string, callback common.DeviceCallback, log *core.LogAgent) {
	sm.DevError  = errCode
	sm.DevReply  = common.NewError(errCode, reason).Error()
	// ExecuteError processing
	if sm.DevError != common.DevErrorSuccess {
		query := &common.DeviceError{
			Action:   sm.DevAction,
			DevState: sm.DevState,
			ErrCode:  sm.DevError,
			ErrText:  sm.DevReply,
		}
		if callback != nil {
			_ = callback.ExecuteError("", query)
		}
		if log != nil {
			log.Debug("ExecuteError: Action - %s, State - %s, Error %d - %s",
				sm.DevAction.String(), sm.DevState.String(), sm.DevError, sm.DevReply)
		}
	}
}

func (sm *Simulator) RunActionPrompt(prompt common.EnumDevPrompt, callback common.DeviceCallback, log *core.LogAgent) {
	sm.DevPrompt = prompt
	// ActionPrompt processing
	if sm.DevPrompt != common.DevPromptNone {
		query := &common.DevicePrompt{
			Action: sm.DevAction,
			Prompt: sm.DevPrompt,
		}
		if callback != nil {
			_ = callback.ActionPrompt("", query)
		}
		if log != nil {
			log.Debug("ActionPrompt: Action - %s, Prompt - %s",
				sm.DevAction.String(), sm.DevPrompt.String())
		}
	}
}

func (sm *Simulator) RunReaderReturn(inform string, callback common.DeviceCallback, log *core.LogAgent) {
	sm.DevInform = inform
	// ReaderReturn processing
	if sm.DevInform != "" {
		query := &common.DeviceInform{
			Action: sm.DevAction,
			Inform: sm.DevInform,
		}
		if callback != nil {
			_ = callback.ReaderReturn("", query)
		}
		if log != nil {
			log.Debug("ReaderReturn: Action - %s, Inform - %s",
				sm.DevAction.String(), sm.DevInform)
		}
	}
}

