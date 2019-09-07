package proxy

import (
	"encoding/json"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type ValidatorClient struct {
	scopeItem *duplex.ScopeItem
	commands  common.ValidatorManager
	log       *core.LogAgent
}

func NewValidatorClient() *ValidatorClient {
	vc := ValidatorClient{
		scopeItem: duplex.NewScopeItem(duplex.ScopeValidator),
		commands:  nil,
		log:       nil,
	}
	return &vc
}

func (vc *ValidatorClient) GetScopeItem() *duplex.ScopeItem {
	return vc.scopeItem
}

func (vc *ValidatorClient) Init(command common.ValidatorManager, log *core.LogAgent) {
	vc.log = log
	vc.commands = command
	// init scope functions
	if vc.scopeItem != nil {
		vc.scopeItem.SetScopeFunc(common.CmdInitValidator, func(name string, dump []byte) {
			query := &common.ValidatorQuery{}
			err := vc.decodeQuery(name, common.CmdInitValidator, dump, query)
			if err == nil && vc.commands != nil {
				err = vc.commands.InitValidator(name, query)
			}
		})
		vc.scopeItem.SetScopeFunc(common.CmdDoValidate, func(name string, dump []byte) {
			query := &common.ValidatorQuery{}
			err := vc.decodeQuery(name, common.CmdDoValidate, dump, query)
			if err == nil && vc.commands != nil {
				err = vc.commands.DoValidate(name, query)
			}
		})
		vc.scopeItem.SetScopeFunc(common.CmdNoteAccept, func(name string, dump []byte) {
			query := &common.ValidatorQuery{}
			err := vc.decodeQuery(name, common.CmdNoteAccept, dump, query)
			if err == nil && vc.commands != nil {
				err = vc.commands.NoteAccept(name, query)
			}
		})
		vc.scopeItem.SetScopeFunc(common.CmdNoteReject, func(name string, dump []byte) {
			query := &common.ValidatorQuery{}
			err := vc.decodeQuery(name, common.CmdNoteReject, dump, query)
			if err == nil && vc.commands != nil {
				err = vc.commands.NoteReject(name, query)
			}
		})
		vc.scopeItem.SetScopeFunc(common.CmdStopValidate, func(name string, dump []byte) {
			query := &common.ValidatorQuery{}
			err := vc.decodeQuery(name, common.CmdStopValidate, dump, query)
			if err == nil && vc.commands != nil {
				err = vc.commands.StopValidate(name, query)
			}
		})
		vc.scopeItem.SetScopeFunc(common.CmdCheckStore, func(name string, dump []byte) {
			query := &common.ValidatorQuery{}
			err := vc.decodeQuery(name, common.CmdCheckStore, dump, query)
			if err == nil && vc.commands != nil {
				err = vc.commands.CheckStore(name, query)
			}
		})
		vc.scopeItem.SetScopeFunc(common.CmdClearStore, func(name string, dump []byte) {
			query := &common.ValidatorQuery{}
			err := vc.decodeQuery(name, common.CmdClearStore, dump, query)
			if err == nil && vc.commands != nil {
				err = vc.commands.ClearStore(name, query)
			}
		})
	}
}

func (vc *ValidatorClient) decodeQuery(name string, cmd string, dump []byte, query interface{}) error {
	if vc.log != nil {
		vc.log.Dump("ValidatorClient dev:%s take cmd:%s, pack:%s", name, cmd, string(dump))
	}
	query = &common.DeviceQuery{}
	return json.Unmarshal(dump, query)
}
