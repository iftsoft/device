package proxy

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
func (os *ObjectState) CommandReply(name string, reply *common.SystemReply) error {
	if os.log != nil {
		os.log.Trace("ObjectState dev:%s get cmd:%s data:%s",
			name, reply.Command, reply.DevName)
	}
	return nil
}
