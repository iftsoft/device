package proxy

import (
	"encoding/json"
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type DeviceClient struct {
	commands  common.DeviceManager
	log       *core.LogAgent
}

func NewDeviceClient() *DeviceClient {
	dc := DeviceClient{
		commands:  nil,
		log:       nil,
	}
	return &dc
}

func (dc *DeviceClient) GetDispatcher() duplex.Dispatcher {
	return dc
}

func (dc *DeviceClient) Init(command common.DeviceManager, log *core.LogAgent) {
	dc.log = log
	dc.commands = command
}

func (dc *DeviceClient) EvalPacket(pack *duplex.Packet) error {
	if pack == nil {
		return errors.New("duplex Packet is nil")
	}
	switch pack.Command {
	case common.CmdDeviceCancel:
		query := &common.DeviceQuery{}
		err := dc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
		if err == nil && dc.commands != nil {
			err = dc.commands.Cancel(pack.DevName, query)
		}
		return err

	case common.CmdDeviceReset:
			query := &common.DeviceQuery{}
			err := dc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
			if err == nil && dc.commands != nil {
				err = dc.commands.Reset(pack.DevName, query)
			}
		return err

	case common.CmdDeviceStatus:
			query := &common.DeviceQuery{}
			err := dc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
			if err == nil && dc.commands != nil {
				err = dc.commands.Status(pack.DevName, query)
			}
		return err

	case common.CmdRunAction:
			query := &common.DeviceQuery{}
			err := dc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
			if err == nil && dc.commands != nil {
				err = dc.commands.RunAction(pack.DevName, query)
			}
		return err

	case common.CmdStopAction:
			query := &common.DeviceQuery{}
			err := dc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
			if err == nil && dc.commands != nil {
				err = dc.commands.StopAction(pack.DevName, query)
			}
		return err

	default:
		dc.log.Warn("DeviceClient EvalPacket: Unknown  command - %s", pack.Command)
		return errors.New("duplex Packet unknown command")
	}
}

func (dc *DeviceClient) decodeQuery(name string, cmd string, dump []byte, query interface{}) error {
	if dc.log != nil {
		dc.log.Dump("DeviceClient dev:%s take cmd:%s, pack:%s", name, cmd, string(dump))
	}
	return json.Unmarshal(dump, query)
}
