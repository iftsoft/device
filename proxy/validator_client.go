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
		vc.scopeItem.SetScopeFunc(common.CmdNoteReturn, func(name string, dump []byte) {
			query := &common.ValidatorQuery{}
			err := vc.decodeQuery(name, common.CmdNoteReturn, dump, query)
			if err == nil && vc.commands != nil {
				err = vc.commands.NoteReturn(name, query)
			}
		})
		vc.scopeItem.SetScopeFunc(common.CmdStopValidate, func(name string, dump []byte) {
			query := &common.ValidatorQuery{}
			err := vc.decodeQuery(name, common.CmdStopValidate, dump, query)
			if err == nil && vc.commands != nil {
				err = vc.commands.StopValidate(name, query)
			}
		})
		vc.scopeItem.SetScopeFunc(common.CmdCheckValidator, func(name string, dump []byte) {
			query := &common.ValidatorQuery{}
			err := vc.decodeQuery(name, common.CmdCheckValidator, dump, query)
			if err == nil && vc.commands != nil {
				err = vc.commands.CheckValidator(name, query)
			}
		})
		vc.scopeItem.SetScopeFunc(common.CmdClearValidator, func(name string, dump []byte) {
			query := &common.ValidatorQuery{}
			err := vc.decodeQuery(name, common.CmdClearValidator, dump, query)
			if err == nil && vc.commands != nil {
				err = vc.commands.ClearValidator(name, query)
			}
		})
	}
}

func (vc *ValidatorClient) decodeQuery(name string, cmd string, dump []byte, query interface{}) error {
	if vc.log != nil {
		vc.log.Dump("ValidatorClient dev:%s take cmd:%s, pack:%s", name, cmd, string(dump))
	}
	return json.Unmarshal(dump, query)
}
