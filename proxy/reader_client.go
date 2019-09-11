package proxy

import (
	"encoding/json"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type ReaderClient struct {
	scopeItem *duplex.ScopeItem
	commands  common.ReaderManager
	log       *core.LogAgent
}

func NewReaderClient() *ReaderClient {
	rc := ReaderClient{
		scopeItem: duplex.NewScopeItem(duplex.ScopeReader),
		commands:  nil,
		log:       nil,
	}
	return &rc
}

func (rc *ReaderClient) GetScopeItem() *duplex.ScopeItem {
	return rc.scopeItem
}

func (rc *ReaderClient) Init(command common.ReaderManager, log *core.LogAgent) {
	rc.log = log
	rc.commands = command
	// init scope functions
	if rc.scopeItem != nil {
		rc.scopeItem.SetScopeFunc(common.CmdEnterCard, func(name string, dump []byte) {
			query := &common.DeviceQuery{}
			err := rc.decodeQuery(name, common.CmdEnterCard, dump, query)
			if err == nil && rc.commands != nil {
				err = rc.commands.EnterCard(name, query)
			}
		})
		rc.scopeItem.SetScopeFunc(common.CmdEjectCard, func(name string, dump []byte) {
			query := &common.DeviceQuery{}
			err := rc.decodeQuery(name, common.CmdEjectCard, dump, query)
			if err == nil && rc.commands != nil {
				err = rc.commands.EjectCard(name, query)
			}
		})
		rc.scopeItem.SetScopeFunc(common.CmdCaptureCard, func(name string, dump []byte) {
			query := &common.DeviceQuery{}
			err := rc.decodeQuery(name, common.CmdCaptureCard, dump, query)
			if err == nil && rc.commands != nil {
				err = rc.commands.CaptureCard(name, query)
			}
		})
		rc.scopeItem.SetScopeFunc(common.CmdReadCard, func(name string, dump []byte) {
			query := &common.DeviceQuery{}
			err := rc.decodeQuery(name, common.CmdReadCard, dump, query)
			if err == nil && rc.commands != nil {
				err = rc.commands.ReadCard(name, query)
			}
		})
		rc.scopeItem.SetScopeFunc(common.CmdChipGetATR, func(name string, dump []byte) {
			query := &common.DeviceQuery{}
			err := rc.decodeQuery(name, common.CmdChipGetATR, dump, query)
			if err == nil && rc.commands != nil {
				err = rc.commands.ChipGetATR(name, query)
			}
		})
		rc.scopeItem.SetScopeFunc(common.CmdChipPowerOff, func(name string, dump []byte) {
			query := &common.DeviceQuery{}
			err := rc.decodeQuery(name, common.CmdChipPowerOff, dump, query)
			if err == nil && rc.commands != nil {
				err = rc.commands.ChipPowerOff(name, query)
			}
		})
		rc.scopeItem.SetScopeFunc(common.CmdChipCommand, func(name string, dump []byte) {
			query := &common.ReaderChipQuery{}
			err := rc.decodeQuery(name, common.CmdChipCommand, dump, query)
			if err == nil && rc.commands != nil {
				err = rc.commands.ChipCommand(name, query)
			}
		})
		rc.scopeItem.SetScopeFunc(common.CmdReadPIN, func(name string, dump []byte) {
			query := &common.ReaderPinQuery{}
			err := rc.decodeQuery(name, common.CmdReadPIN, dump, query)
			if err == nil && rc.commands != nil {
				err = rc.commands.ReadPIN(name, query)
			}
		})
		rc.scopeItem.SetScopeFunc(common.CmdLoadMasterKey, func(name string, dump []byte) {
			query := &common.ReaderPinQuery{}
			err := rc.decodeQuery(name, common.CmdLoadMasterKey, dump, query)
			if err == nil && rc.commands != nil {
				err = rc.commands.LoadMasterKey(name, query)
			}
		})
		rc.scopeItem.SetScopeFunc(common.CmdLoadWorkKey, func(name string, dump []byte) {
			query := &common.ReaderPinQuery{}
			err := rc.decodeQuery(name, common.CmdLoadWorkKey, dump, query)
			if err == nil && rc.commands != nil {
				err = rc.commands.LoadWorkKey(name, query)
			}
		})
		rc.scopeItem.SetScopeFunc(common.CmdTestMasterKey, func(name string, dump []byte) {
			query := &common.ReaderPinQuery{}
			err := rc.decodeQuery(name, common.CmdTestMasterKey, dump, query)
			if err == nil && rc.commands != nil {
				err = rc.commands.TestMasterKey(name, query)
			}
		})
		rc.scopeItem.SetScopeFunc(common.CmdTestWorkKey, func(name string, dump []byte) {
			query := &common.ReaderPinQuery{}
			err := rc.decodeQuery(name, common.CmdTestWorkKey, dump, query)
			if err == nil && rc.commands != nil {
				err = rc.commands.TestWorkKey(name, query)
			}
		})
	}
}

func (rc *ReaderClient) decodeQuery(name string, cmd string, dump []byte, query interface{}) error {
	if rc.log != nil {
		rc.log.Dump("ReaderClient dev:%s take cmd:%s, pack:%s", name, cmd, string(dump))
	}
	return json.Unmarshal(dump, query)
}
