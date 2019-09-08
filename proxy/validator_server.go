package proxy

import (
	"encoding/json"
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type ValidatorServer struct {
	scopeItem *duplex.ScopeItem
	server    duplex.ServerManager
	callback  common.ValidatorCallback
	log       *core.LogAgent
}

func NewValidatorServer() *ValidatorServer {
	vs := ValidatorServer{
		scopeItem: duplex.NewScopeItem(duplex.ScopeValidator),
		server:    nil,
		callback:  nil,
		log:       nil,
	}
	return &vs
}

func (vs *ValidatorServer) GetScopeItem() *duplex.ScopeItem {
	return vs.scopeItem
}

func (vs *ValidatorServer) Init(server duplex.ServerManager, callback common.ValidatorCallback, log *core.LogAgent) {
	vs.log = log
	vs.server = server
	vs.callback = callback
	if vs.scopeItem != nil {
		vs.scopeItem.SetScopeFunc(common.CmdNoteAccepted, func(name string, dump []byte) {
			reply := &common.ValidatorAccept{}
			err := vs.decodeReply(name, common.CmdNoteAccepted, dump, reply)
			if err == nil && vs.callback != nil {
				err = vs.callback.NoteAccepted(name, reply)
			}
		})
		vs.scopeItem.SetScopeFunc(common.CmdCashIsStored, func(name string, dump []byte) {
			reply := &common.ValidatorAccept{}
			err := vs.decodeReply(name, common.CmdCashIsStored, dump, reply)
			if err == nil && vs.callback != nil {
				err = vs.callback.CashIsStored(name, reply)
			}
		})
		vs.scopeItem.SetScopeFunc(common.CmdCashReturned, func(name string, dump []byte) {
			reply := &common.ValidatorAccept{}
			err := vs.decodeReply(name, common.CmdCashReturned, dump, reply)
			if err == nil && vs.callback != nil {
				err = vs.callback.CashReturned(name, reply)
			}
		})
		vs.scopeItem.SetScopeFunc(common.CmdValidatorStore, func(name string, dump []byte) {
			reply := &common.ValidatorStore{}
			err := vs.decodeReply(name, common.CmdValidatorStore, dump, reply)
			if err == nil && vs.callback != nil {
				err = vs.callback.ValidatorStore(name, reply)
			}
		})
		if vs.server != nil {
			vs.server.AddScopeItem(vs.scopeItem)
		}
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
