package proxy

import (
	"encoding/json"
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type ReaderServer struct {
	server    duplex.ServerManager
	callback  common.ReaderCallback
	log       *core.LogAgent
}

func NewReaderServer() *ReaderServer {
	rs := ReaderServer{
		server:    nil,
		callback:  nil,
		log:       nil,
	}
	return &rs
}

func (rs *ReaderServer) Init(server duplex.ServerManager, callback common.ReaderCallback, log *core.LogAgent) {
	rs.log = log
	rs.server = server
	rs.callback = callback
	if rs.server != nil {
		rs.server.AddDispatcher(duplex.ScopeReader, rs)
	}
}

func (rs *ReaderServer) EvalPacket(pack *duplex.Packet) error {
	if pack == nil {
		return errors.New("duplex Packet is nil")
	}
	switch pack.Command {
	case common.CmdCardPosition:
		reply := &common.ReaderCardPos{}
		err := rs.decodeReply(pack.DevName, pack.Command, pack.Content, reply)
		if err == nil && rs.callback != nil {
			err = rs.callback.CardPosition(pack.DevName, reply)
		}
		return err

	case common.CmdCardDescription:
		reply := &common.ReaderCardInfo{}
		err := rs.decodeReply(pack.DevName, pack.Command, pack.Content, reply)
		if err == nil && rs.callback != nil {
			err = rs.callback.CardDescription(pack.DevName, reply)
		}
		return err

	case common.CmdChipResponse:
		reply := &common.ReaderChipReply{}
		err := rs.decodeReply(pack.DevName, pack.Command, pack.Content, reply)
		if err == nil && rs.callback != nil {
			err = rs.callback.ChipResponse(pack.DevName, reply)
		}
		return err

	default:
		rs.log.Warn("ReaderServer EvalPacket: Unknown  command - %s", pack.Command)
		return errors.New("duplex Packet unknown command")
	}
}

func (rs *ReaderServer) decodeReply(name string, cmd string, dump []byte, reply interface{}) (err error) {
	if rs.log != nil {
		rs.log.Dump("ReaderServer dev:%s take cmd:%s, pack:%s", name, cmd, string(dump))
	}
	err = json.Unmarshal(dump, reply)
	return err
}

func (rs *ReaderServer) SendReaderCommand(name string, cmd string, query interface{}) error {
	if rs.server == nil {
		return errors.New("ServerManager is not set for ReaderServer")
	}
	transport := rs.server.GetTransporter(name)
	if transport == nil {
		return errors.New("ReaderServer can't get transport to device")
	}
	dump, err := json.Marshal(query)
	if err != nil {
		return err
	}
	if rs.log != nil {
		rs.log.Dump("ReaderServer dev:%s send cmd:%s, pack:%s", name, cmd, string(dump))
	}
	pack := duplex.NewPacket(duplex.ScopeReader, name, cmd, dump)
	err = transport.SendPacket(pack)
	return err
}
