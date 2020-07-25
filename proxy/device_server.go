package proxy

import (
	"encoding/json"
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type DeviceServer struct {
	server    duplex.ServerManager
	callback  common.DeviceCallback
	log       *core.LogAgent
}

func NewDeviceServer() *DeviceServer {
	ss := DeviceServer{
		server:    nil,
		callback:  nil,
		log:       nil,
	}
	return &ss
}


func (ss *DeviceServer) Init(server duplex.ServerManager, callback common.DeviceCallback, log *core.LogAgent) {
	ss.log = log
	ss.server = server
	ss.callback = callback
	if ss.server != nil {
		ss.server.AddDispatcher(duplex.ScopeDevice, ss)
	}
}

func (ss *DeviceServer) EvalPacket(pack *duplex.Packet) error {
	if pack == nil {
		return errors.New("duplex Packet is nil")
	}
	switch pack.Command {
	case common.CmdDeviceReply:
		reply := &common.DeviceReply{}
		err := ss.decodeReply(pack.DevName, pack.Command, pack.Content, reply)
		if err == nil && ss.callback != nil {
			err = ss.callback.DeviceReply(pack.DevName, reply)
		}
		return err

	case common.CmdExecuteError:
			reply := &common.DeviceError{}
			err := ss.decodeReply(pack.DevName, pack.Command, pack.Content, reply)
			if err == nil && ss.callback != nil {
				err = ss.callback.ExecuteError(pack.DevName, reply)
			}
		return err

	case common.CmdStateChanged:
			reply := &common.DeviceState{}
			err := ss.decodeReply(pack.DevName, pack.Command, pack.Content, reply)
			if err == nil && ss.callback != nil {
				err = ss.callback.StateChanged(pack.DevName, reply)
			}
		return err

	case common.CmdActionPrompt:
			reply := &common.DevicePrompt{}
			err := ss.decodeReply(pack.DevName, pack.Command, pack.Content, reply)
			if err == nil && ss.callback != nil {
				err = ss.callback.ActionPrompt(pack.DevName, reply)
			}
		return err

	case common.CmdReaderReturn:
			reply := &common.DeviceInform{}
			err := ss.decodeReply(pack.DevName, pack.Command, pack.Content, reply)
			if err == nil && ss.callback != nil {
				err = ss.callback.ReaderReturn(pack.DevName, reply)
			}
		return err

	default:
		ss.log.Warn("DeviceServer EvalPacket: Unknown  command - %s", pack.Command)
		return errors.New("duplex Packet unknown command")
	}
}

func (ss *DeviceServer) decodeReply(name string, cmd string, dump []byte, reply interface{}) (err error) {
	if ss.log != nil {
		ss.log.Dump("DeviceServer dev:%s take cmd:%s, pack:%s", name, cmd, string(dump))
	}
	err = json.Unmarshal(dump, reply)
	return err
}

func (ss *DeviceServer) SendDeviceCommand(name string, cmd string, query interface{}) error {
	if ss.server == nil {
		return errors.New("ServerManager is not set for DeviceServer")
	}
	transport := ss.server.GetTransporter(name)
	if transport == nil {
		return errors.New("DeviceServer can't get transport to device")
	}
	dump, err := json.Marshal(query)
	if err != nil {
		return err
	}
	if ss.log != nil {
		ss.log.Dump("DeviceServer dev:%s send cmd:%s, pack:%s", name, cmd, string(dump))
	}
	pack := duplex.NewPacket(duplex.ScopeDevice, name, cmd, dump)
	err = transport.SendPacket(pack)
	return err
}
