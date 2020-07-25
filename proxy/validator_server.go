package proxy

import (
	"encoding/json"
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type ValidatorServer struct {
	server    duplex.ServerManager
	callback  common.ValidatorCallback
	log       *core.LogAgent
}

func NewValidatorServer() *ValidatorServer {
	vs := ValidatorServer{
		server:    nil,
		callback:  nil,
		log:       nil,
	}
	return &vs
}


func (vs *ValidatorServer) Init(server duplex.ServerManager, callback common.ValidatorCallback, log *core.LogAgent) {
	vs.log = log
	vs.server = server
	vs.callback = callback
	if vs.server != nil {
		vs.server.AddDispatcher(duplex.ScopeValidator, vs)
	}
}

func (vs *ValidatorServer) EvalPacket(pack *duplex.Packet) error {
	if pack == nil {
		return errors.New("duplex Packet is nil")
	}
	switch pack.Command {
	case common.CmdNoteAccepted:
		reply := &common.ValidatorAccept{}
		err := vs.decodeReply(pack.DevName, pack.Command, pack.Content, reply)
		if err == nil && vs.callback != nil {
			err = vs.callback.NoteAccepted(pack.DevName, reply)
		}
		return err

	case common.CmdCashIsStored:
		reply := &common.ValidatorAccept{}
		err := vs.decodeReply(pack.DevName, pack.Command, pack.Content, reply)
		if err == nil && vs.callback != nil {
			err = vs.callback.CashIsStored(pack.DevName, reply)
		}
		return err

	case common.CmdCashReturned:
		reply := &common.ValidatorAccept{}
		err := vs.decodeReply(pack.DevName, pack.Command, pack.Content, reply)
		if err == nil && vs.callback != nil {
			err = vs.callback.CashReturned(pack.DevName, reply)
		}
		return err

	case common.CmdValidatorStore:
		reply := &common.ValidatorStore{}
		err := vs.decodeReply(pack.DevName, pack.Command, pack.Content, reply)
		if err == nil && vs.callback != nil {
			err = vs.callback.ValidatorStore(pack.DevName, reply)
		}
		return err

	default:
		vs.log.Warn("ValidatorServer EvalPacket: Unknown  command - %s", pack.Command)
		return errors.New("duplex Packet unknown command")
	}
}

func (vs *ValidatorServer) decodeReply(name string, cmd string, dump []byte, reply interface{}) (err error) {
	if vs.log != nil {
		vs.log.Dump("ValidatorServer dev:%s take cmd:%s, pack:%s", name, cmd, string(dump))
	}
	err = json.Unmarshal(dump, reply)
	return err
}

func (vs *ValidatorServer) SendValidatorCommand(name string, cmd string, query interface{}) error {
	if vs.server == nil {
		return errors.New("ServerManager is not set for ValidatorServer")
	}
	transport := vs.server.GetTransporter(name)
	if transport == nil {
		return errors.New("ValidatorServer can't get transport to device")
	}
	dump, err := json.Marshal(query)
	if err != nil {
		return err
	}
	if vs.log != nil {
		vs.log.Dump("ValidatorServer dev:%s send cmd:%s, pack:%s", name, cmd, string(dump))
	}
	pack := duplex.NewPacket(duplex.ScopeValidator, name, cmd, dump)
	err = transport.SendPacket(pack)
	return err
}
