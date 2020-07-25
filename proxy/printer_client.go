package proxy

import (
	"encoding/json"
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type PrinterClient struct {
	commands  common.PrinterManager
	log       *core.LogAgent
}

func NewPrinterClient() *PrinterClient {
	pc := PrinterClient{
		commands:  nil,
		log:       nil,
	}
	return &pc
}

func (pc *PrinterClient) GetDispatcher() duplex.Dispatcher {
	return pc
}

func (pc *PrinterClient) Init(command common.PrinterManager, log *core.LogAgent) {
	pc.log = log
	pc.commands = command
}

func (pc *PrinterClient) EvalPacket(pack *duplex.Packet) error {
	if pack == nil {
		return errors.New("duplex Packet is nil")
	}
	switch pack.Command {
	case common.CmdInitPrinter:
		query := &common.PrinterSetup{}
		err := pc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
		if err == nil && pc.commands != nil {
			err = pc.commands.InitPrinter(pack.DevName, query)
		}
		return err

	case common.CmdPrintText:
		query := &common.PrinterQuery{}
		err := pc.decodeQuery(pack.DevName, pack.Command, pack.Content, query)
		if err == nil && pc.commands != nil {
			err = pc.commands.PrintText(pack.DevName, query)
		}
		return err

	default:
		pc.log.Warn("PrinterClient EvalPacket: Unknown  command - %s", pack.Command)
		return errors.New("duplex Packet unknown command")
	}
}

func (pc *PrinterClient) decodeQuery(name string, cmd string, dump []byte, query interface{}) error {
	if pc.log != nil {
		pc.log.Dump("PrinterClient dev:%s take cmd:%s, pack:%s", name, cmd, string(dump))
	}
	return json.Unmarshal(dump, query)
}
