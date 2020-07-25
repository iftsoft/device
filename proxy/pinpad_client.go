package proxy

import (
	"encoding/json"
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type PinPadClient struct {
	commands  common.PinPadManager
	log       *core.LogAgent
}

func NewPinPadClient() *PinPadClient {
	ppc := PinPadClient{
		commands:  nil,
		log:       nil,
	}
	return &ppc
}

func (ppc *PinPadClient) GetDispatcher() duplex.Dispatcher {
	return ppc
}

func (ppc *PinPadClient) Init(command common.PinPadManager, log *core.LogAgent) {
	ppc.log = log
	ppc.commands = command
}

func (ppc *PinPadClient) EvalPacket(pack *duplex.Packet) error {
	if pack == nil {
		return errors.New("duplex Packet is nil")
	}
	switch pack.Command {
	case common.CmdReadPIN:
		query := &common.ReaderPinQuery{}
		err := ppc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
		if err == nil && ppc.commands != nil {
			err = ppc.commands.ReadPIN(pack.DevName, query)
		}
		return err

	case common.CmdLoadMasterKey:
		query := &common.ReaderPinQuery{}
		err := ppc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
		if err == nil && ppc.commands != nil {
			err = ppc.commands.LoadMasterKey(pack.DevName, query)
		}
		return err

	case common.CmdLoadWorkKey:
		query := &common.ReaderPinQuery{}
		err := ppc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
		if err == nil && ppc.commands != nil {
			err = ppc.commands.LoadWorkKey(pack.DevName, query)
		}
		return err

	case common.CmdTestMasterKey:
		query := &common.ReaderPinQuery{}
		err := ppc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
		if err == nil && ppc.commands != nil {
			err = ppc.commands.TestMasterKey(pack.DevName, query)
		}
		return err

	case common.CmdTestWorkKey:
		query := &common.ReaderPinQuery{}
		err := ppc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
		if err == nil && ppc.commands != nil {
			err = ppc.commands.TestWorkKey(pack.DevName, query)
		}
		return err

	default:
		ppc.log.Warn("PinPadClient EvalPacket: Unknown  command - %s", pack.Command)
		return errors.New("duplex Packet unknown command")
	}
}

func (ppc *PinPadClient) decodeQuery(name string, cmd string, dump []byte, query interface{}) error {
	if ppc.log != nil {
		ppc.log.Dump("PinPadClient dev:%s take cmd:%s, pack:%s", name, cmd, string(dump))
	}
	return json.Unmarshal(dump, query)
}
