package system

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
)

type SystemBack struct {
	log *core.LogAgent
}

func NewSystemBack() *SystemBack {
	sb := SystemBack{
		log: nil,
	}
	return &sb
}

func (sb *SystemBack) Init(log *core.LogAgent) {
	sb.log = log
}

func (sb *SystemBack) CommandReply(name string, reply *common.SystemReply) error {
	if sb.log != nil {
		sb.log.Trace("SystemBack dev:%s get cmd:%s data:%s",
			name, reply.Command, reply.DevName)
	}
	return nil
}
