package system

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
)

type SystemStub struct {
	callback common.SystemCallback
	log      *core.LogAgent
}

func NewSystemStub() *SystemStub {
	ss := SystemStub{
		callback: nil,
		log:      nil,
	}
	return &ss
}

func (ss *SystemStub) Init(callback common.SystemCallback, log *core.LogAgent) {
	ss.log = log
	ss.callback = callback
}

// Implemetation of common.SystemManager

func (st *SystemStub) Config(query *common.SystemQuery) error {
	if st.log != nil {
		st.log.Trace("SystemStub run cmd:Config, pack:%s", query.DevName)
	}
	reply := &common.SystemReply{}
	reply.DevName = query.DevName
	reply.Command = "Config"
	var err error
	if st.callback != nil {
		err = st.callback.CommandReply(reply)
	}
	return err
}

func (st *SystemStub) Inform(query *common.SystemQuery) error {
	if st.log != nil {
		st.log.Trace("SystemStub run cmd:Inform, pack:%s", query.DevName)
	}
	reply := &common.SystemReply{}
	reply.DevName = query.DevName
	reply.Command = "Inform"
	var err error
	if st.callback != nil {
		err = st.callback.CommandReply(reply)
	}
	return err
}

func (st *SystemStub) Start(query *common.SystemQuery) error {
	if st.log != nil {
		st.log.Trace("SystemStub run cmd:Start, pack:%s", query.DevName)
	}
	reply := &common.SystemReply{}
	reply.DevName = query.DevName
	reply.Command = "Start"
	var err error
	if st.callback != nil {
		err = st.callback.CommandReply(reply)
	}
	return err
}

func (st *SystemStub) Stop(query *common.SystemQuery) error {
	if st.log != nil {
		st.log.Trace("SystemStub run cmd:Stop, pack:%s", query.DevName)
	}
	reply := &common.SystemReply{}
	reply.DevName = query.DevName
	reply.Command = "Stop"
	var err error
	if st.callback != nil {
		err = st.callback.CommandReply(reply)
	}
	return err
}

func (st *SystemStub) Restart(query *common.SystemQuery) error {
	if st.log != nil {
		st.log.Trace("SystemStub run cmd:Restart, pack:%s", query.DevName)
	}
	reply := &common.SystemReply{}
	reply.DevName = query.DevName
	reply.Command = "Restart"
	var err error
	if st.callback != nil {
		err = st.callback.CommandReply(reply)
	}
	return err
}
