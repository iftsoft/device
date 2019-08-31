package handler

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
)

type ObjectState struct {
	DevName string
	log     *core.LogAgent
}

func NewObjectState() *ObjectState {
	os := ObjectState{
		DevName: "",
		log:     nil,
	}
	return &os
}

func (os *ObjectState) Init(name string, log *core.LogAgent) {
	os.DevName = name
	os.log = log
}

// Implementation of common.SystemCallback
func (os *ObjectState) SystemReply(name string, reply *common.SystemReply) error {
	if os.log != nil {
		os.log.Debug("ObjectState dev:%s get cmd:%s",
			name, reply.Command)
	}
	return nil
}

// Implementation of common.DeviceCallback
func (os *ObjectState) DeviceReply(name string, reply *common.DeviceReply) error {
	if os.log != nil {
		os.log.Debug("ObjectState.DeviceReply dev:%s, cmd:%s, state:%d",
			name, reply.Command, reply.DevState)
	}
	return nil
}

func (os *ObjectState) ExecuteError(name string, reply *common.DeviceError) error {
	if os.log != nil {
		os.log.Debug("ObjectState.ExecuteError dev:%s, action:%d, error:%d - %s",
			name, reply.Action, reply.ErrCode, reply.ErrText)
	}
	return nil
}

func (os *ObjectState) StateChanged(name string, reply *common.DeviceState) error {
	if os.log != nil {
		os.log.Debug("ObjectState.StateChanged dev:%s, old state:%d, new state:%d",
			name, reply.OldState, reply.NewState)
	}
	return nil
}

func (os *ObjectState) ActionPrompt(name string, reply *common.DevicePrompt) error {
	if os.log != nil {
		os.log.Debug("ObjectState.ActionPrompt dev:%s, action:%d, prompt:%d",
			name, reply.Action, reply.Prompt)
	}
	return nil
}
