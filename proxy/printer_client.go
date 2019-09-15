package proxy

import (
	"encoding/json"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type PrinterClient struct {
	scopeItem *duplex.ScopeItem
	commands  common.PrinterManager
	log       *core.LogAgent
}

func NewPrinterClient() *PrinterClient {
	pc := PrinterClient{
		scopeItem: duplex.NewScopeItem(duplex.ScopePrinter),
		commands:  nil,
		log:       nil,
	}
	return &pc
}

func (pc *PrinterClient) GetScopeItem() *duplex.ScopeItem {
	return pc.scopeItem
}

func (pc *PrinterClient) Init(command common.PrinterManager, log *core.LogAgent) {
	pc.log = log
	pc.commands = command
	// init scope functions
	if pc.scopeItem != nil {
		pc.scopeItem.SetScopeFunc(common.CmdInitPrinter, func(name string, dump []byte) {
			query := &common.PrinterSetup{}
			err := pc.decodeQuery(name, common.CmdInitPrinter, dump, query)
			if err == nil && pc.commands != nil {
				err = pc.commands.InitPrinter(name, query)
			}
		})
		pc.scopeItem.SetScopeFunc(common.CmdPrintText, func(name string, dump []byte) {
			query := &common.PrinterQuery{}
			err := pc.decodeQuery(name, common.CmdPrintText, dump, query)
			if err == nil && pc.commands != nil {
				err = pc.commands.PrintText(name, query)
			}
		})
	}
}

func (pc *PrinterClient) decodeQuery(name string, cmd string, dump []byte, query interface{}) error {
	if pc.log != nil {
		pc.log.Dump("PrinterClient dev:%s take cmd:%s, pack:%s", name, cmd, string(dump))
	}
	return json.Unmarshal(dump, query)
}
