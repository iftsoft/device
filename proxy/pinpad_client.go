package proxy

import (
	"encoding/json"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type PinPadClient struct {
	scopeItem *duplex.ScopeItem
	commands  common.PinPadManager
	log       *core.LogAgent
}

func NewPinPadClient() *PinPadClient {
	ppc := PinPadClient{
		scopeItem: duplex.NewScopeItem(duplex.ScopePinPad),
		commands:  nil,
		log:       nil,
	}
	return &ppc
}

func (ppc *PinPadClient) GetScopeItem() *duplex.ScopeItem {
	return ppc.scopeItem
}

func (ppc *PinPadClient) Init(command common.PinPadManager, log *core.LogAgent) {
	ppc.log = log
	ppc.commands = command
	// init scope functions
	if ppc.scopeItem != nil {
		ppc.scopeItem.SetScopeFunc(common.CmdReadPIN, func(name string, dump []byte) {
			query := &common.ReaderPinQuery{}
			err := ppc.decodeQuery(name, common.CmdReadPIN, dump, query)
			if err == nil && ppc.commands != nil {
				err = ppc.commands.ReadPIN(name, query)
			}
		})
		ppc.scopeItem.SetScopeFunc(common.CmdLoadMasterKey, func(name string, dump []byte) {
			query := &common.ReaderPinQuery{}
			err := ppc.decodeQuery(name, common.CmdLoadMasterKey, dump, query)
			if err == nil && ppc.commands != nil {
				err = ppc.commands.LoadMasterKey(name, query)
			}
		})
		ppc.scopeItem.SetScopeFunc(common.CmdLoadWorkKey, func(name string, dump []byte) {
			query := &common.ReaderPinQuery{}
			err := ppc.decodeQuery(name, common.CmdLoadWorkKey, dump, query)
			if err == nil && ppc.commands != nil {
				err = ppc.commands.LoadWorkKey(name, query)
			}
		})
		ppc.scopeItem.SetScopeFunc(common.CmdTestMasterKey, func(name string, dump []byte) {
			query := &common.ReaderPinQuery{}
			err := ppc.decodeQuery(name, common.CmdTestMasterKey, dump, query)
			if err == nil && ppc.commands != nil {
				err = ppc.commands.TestMasterKey(name, query)
			}
		})
		ppc.scopeItem.SetScopeFunc(common.CmdTestWorkKey, func(name string, dump []byte) {
			query := &common.ReaderPinQuery{}
			err := ppc.decodeQuery(name, common.CmdTestWorkKey, dump, query)
			if err == nil && ppc.commands != nil {
				err = ppc.commands.TestWorkKey(name, query)
			}
		})
	}
}

func (ppc *PinPadClient) decodeQuery(name string, cmd string, dump []byte, query interface{}) error {
	if ppc.log != nil {
		ppc.log.Dump("PinPadClient dev:%s take cmd:%s, pack:%s", name, cmd, string(dump))
	}
	return json.Unmarshal(dump, query)
}
