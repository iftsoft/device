package proxy

import (
	"encoding/json"
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type ReaderClient struct {
	commands  common.ReaderManager
	log       *core.LogAgent
}

func NewReaderClient() *ReaderClient {
	rc := ReaderClient{
		commands:  nil,
		log:       nil,
	}
	return &rc
}

func (rc *ReaderClient) GetDispatcher() duplex.Dispatcher {
	return rc
}

func (rc *ReaderClient) Init(command common.ReaderManager, log *core.LogAgent) {
	rc.log = log
	rc.commands = command
}

func (rc *ReaderClient) EvalPacket(pack *duplex.Packet) error {
	if pack == nil {
		return errors.New("duplex Packet is nil")
	}
	switch pack.Command {
	case common.CmdEnterCard:
		query := &common.DeviceQuery{}
		err := rc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
		if err == nil && rc.commands != nil {
			err = rc.commands.EnterCard(pack.DevName, query)
		}
		return err

	case common.CmdEjectCard:
		query := &common.DeviceQuery{}
		err := rc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
		if err == nil && rc.commands != nil {
			err = rc.commands.EjectCard(pack.DevName, query)
		}
		return err

	case common.CmdCaptureCard:
		query := &common.DeviceQuery{}
		err := rc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
		if err == nil && rc.commands != nil {
			err = rc.commands.CaptureCard(pack.DevName, query)
		}
		return err

	case common.CmdReadCard:
		query := &common.DeviceQuery{}
		err := rc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
		if err == nil && rc.commands != nil {
			err = rc.commands.ReadCard(pack.DevName, query)
		}
		return err

	case common.CmdChipGetATR:
		query := &common.DeviceQuery{}
		err := rc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
		if err == nil && rc.commands != nil {
			err = rc.commands.ChipGetATR(pack.DevName, query)
		}
		return err

	case common.CmdChipPowerOff:
		query := &common.DeviceQuery{}
		err := rc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
		if err == nil && rc.commands != nil {
			err = rc.commands.ChipPowerOff(pack.DevName, query)
		}
		return err

	case common.CmdChipCommand:
		query := &common.ReaderChipQuery{}
		err := rc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
		if err == nil && rc.commands != nil {
			err = rc.commands.ChipCommand(pack.DevName, query)
		}
		return err

	default:
		rc.log.Warn("ReaderClient EvalPacket: Unknown  command - %s", pack.Command)
		return errors.New("duplex Packet unknown command")
	}
}

func (rc *ReaderClient) decodeQuery(name string, cmd string, dump []byte, query interface{}) error {
	if rc.log != nil {
		rc.log.Dump("ReaderClient dev:%s take cmd:%s, pack:%s", name, cmd, string(dump))
	}
	return json.Unmarshal(dump, query)
}
