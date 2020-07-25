package proxy

import (
	"encoding/json"
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type SystemClient struct {
//	scopeItem *duplex.ScopeItem
	commands  common.SystemManager
	log       *core.LogAgent
}

func NewSystemClient() *SystemClient {
	sc := SystemClient{
//		scopeItem: duplex.NewScopeItem(duplex.ScopeSystem),
		commands:  nil,
		log:       nil,
	}
	return &sc
}

func (sc *SystemClient) GetDispatcher() duplex.Dispatcher {
	return sc
}

func (sc *SystemClient) Init(command common.SystemManager, log *core.LogAgent) {
	sc.log = log
	sc.commands = command
}

func (sc *SystemClient) EvalPacket(pack *duplex.Packet) error {
	if pack == nil {
		return errors.New("duplex Packet is nil")
	}
	switch pack.Command {
	case common.CmdSystemTerminate:
			query := &common.SystemQuery{}
			err := sc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
			if err == nil && sc.commands != nil {
				err = sc.commands.Terminate(pack.DevName, query)
			}
		return err

	case common.CmdSystemInform:
			query := &common.SystemQuery{}
			err := sc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
			if err == nil && sc.commands != nil {
				err = sc.commands.SysInform(pack.DevName, query)
			}
		return err

	case common.CmdSystemStart:
			query := &common.SystemConfig{}
			err := sc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
			if err == nil && sc.commands != nil {
				err = sc.commands.SysStart(pack.DevName, query)
			}
		return err

	case common.CmdSystemStop:
			query := &common.SystemQuery{}
			err := sc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
			if err == nil && sc.commands != nil {
				err = sc.commands.SysStop(pack.DevName, query)
			}
		return err

	case common.CmdSystemRestart:
			query := &common.SystemConfig{}
			err := sc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
			if err == nil && sc.commands != nil {
				err = sc.commands.SysRestart(pack.DevName, query)
			}
		return err

	default:
		sc.log.Warn("SystemClient EvalPacket: Unknown  command - %s", pack.Command)
		return errors.New("duplex Packet unknown command")
	}
}

func (sc *SystemClient) decodeQuery(name string, cmd string, dump []byte, query interface{}) error {
	if sc.log != nil {
		sc.log.Dump("SystemClient dev:%s take cmd:%s, pack:%s", name, cmd, string(dump))
	}
	return json.Unmarshal(dump, query)
}

