package generic

import (
	"github.com/iftsoft/device/common"
	"sync"
)

////////////////////////////////////////////////////////////////

type MimicStep struct {
	Delay  int
	Value  int
	State  common.EnumDevState
	Error  common.EnumDevError
	Prompt common.EnumDevPrompt
	Action common.EnumDevAction
	Inform string
	Reason string
}

type MimicSteps []MimicStep

////////////////////////////////////////////////////////////////

type Mimicker struct {
	stages MimicSteps
	mutex  sync.Mutex
	index  int
	delay  int
}

func (m *Mimicker) SetSteps(list MimicSteps) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.stages = list
	m.index = 0
	m.delay = 0
}

func (m *Mimicker) GetStep() *MimicStep {
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
	mimicker  Mimicker
}

func (sm *Simulator) SetupMimic(list MimicSteps) {
	sm.mimicker.SetSteps(list)
}

func (sm *Simulator) ClearMimic() {
	sm.mimicker.SetSteps(nil)
}

func (sm *Simulator) GetMimicStep() *MimicStep {
	return sm.mimicker.GetStep()
}

func (sm *Simulator) ProcessStage(be *BaseEngine, stage *MimicStep) error {
	if stage == nil {
		return nil
	}
	var err error
	be.DevAction = stage.Action
	err = be.RunStateChanged(stage.State)
	err = be.RunExecuteError(stage.Error, stage.Reason)
	err = be.RunActionPrompt(stage.Prompt)
	err = be.RunReaderReturn(stage.Inform)
	return err
}
