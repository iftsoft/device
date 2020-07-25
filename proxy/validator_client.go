package proxy

import (
	"encoding/json"
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type ValidatorClient struct {
	commands  common.ValidatorManager
	log       *core.LogAgent
}

func NewValidatorClient() *ValidatorClient {
	vc := ValidatorClient{
		commands:  nil,
		log:       nil,
	}
	return &vc
}

func (vc *ValidatorClient) GetDispatcher() duplex.Dispatcher {
	return vc
}

func (vc *ValidatorClient) Init(command common.ValidatorManager, log *core.LogAgent) {
	vc.log = log
	vc.commands = command
}

func (vc *ValidatorClient) EvalPacket(pack *duplex.Packet) error {
	if pack == nil {
		return errors.New("duplex Packet is nil")
	}
	switch pack.Command {
	case common.CmdInitValidator:
		query := &common.ValidatorQuery{}
		err := vc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
		if err == nil && vc.commands != nil {
			err = vc.commands.InitValidator(pack.DevName, query)
		}
		return err

	case common.CmdDoValidate:
		query := &common.ValidatorQuery{}
		err := vc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
		if err == nil && vc.commands != nil {
			err = vc.commands.DoValidate(pack.DevName, query)
		}
		return err

	case common.CmdNoteAccept:
		query := &common.ValidatorQuery{}
		err := vc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
		if err == nil && vc.commands != nil {
			err = vc.commands.NoteAccept(pack.DevName, query)
		}
		return err

	case common.CmdNoteReturn:
		query := &common.ValidatorQuery{}
		err := vc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
		if err == nil && vc.commands != nil {
			err = vc.commands.NoteReturn(pack.DevName, query)
		}
		return err

	case common.CmdStopValidate:
		query := &common.ValidatorQuery{}
		err := vc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
		if err == nil && vc.commands != nil {
			err = vc.commands.StopValidate(pack.DevName, query)
		}
		return err

	case common.CmdCheckValidator:
		query := &common.ValidatorQuery{}
		err := vc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
		if err == nil && vc.commands != nil {
			err = vc.commands.CheckValidator(pack.DevName, query)
		}
		return err

	case common.CmdClearValidator:
		query := &common.ValidatorQuery{}
		err := vc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
		if err == nil && vc.commands != nil {
			err = vc.commands.ClearValidator(pack.DevName, query)
		}
		return err

	default:
		vc.log.Warn("ValidatorClient EvalPacket: Unknown  command - %s", pack.Command)
		return errors.New("duplex Packet unknown command")
	}
}

func (vc *ValidatorClient) decodeQuery(name string, cmd string, dump []byte, query interface{}) error {
	if vc.log != nil {
		vc.log.Dump("ValidatorClient dev:%s take cmd:%s, pack:%s", name, cmd, string(dump))
	}
	return json.Unmarshal(dump, query)
}
