package proxy

import (
	"encoding/json"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type DeviceClient struct {
	scopeItem *duplex.ScopeItem
	commands  common.DeviceManager
	log       *core.LogAgent
}

func NewDeviceClient() *DeviceClient {
	dc := DeviceClient{
		scopeItem: duplex.NewScopeItem(duplex.ScopeDevice),
		commands:  nil,
		log:       nil,
	}
	return &dc
}

func (dc *DeviceClient) GetScopeItem() *duplex.ScopeItem {
	return dc.scopeItem
}

func (dc *DeviceClient) Init(command common.DeviceManager, log *core.LogAgent) {
	dc.log = log
	dc.commands = command
	// init scope functions
	if dc.scopeItem != nil {
		dc.scopeItem.SetScopeFunc(common.CmdDeviceCancel, func(name string, dump []byte) {
			query := &common.DeviceQuery{}
			err := dc.decodeQuery(name, common.CmdDeviceCancel, dump, query)
			if err == nil && dc.commands != nil {
				err = dc.commands.Cancel(name, query)
			}
		})
		dc.scopeItem.SetScopeFunc(common.CmdDeviceReset, func(name string, dump []byte) {
			query := &common.DeviceQuery{}
			err := dc.decodeQuery(name, common.CmdDeviceReset, dump, query)
			if err == nil && dc.commands != nil {
				err = dc.commands.Reset(name, query)
			}
		})
		dc.scopeItem.SetScopeFunc(common.CmdDeviceStatus, func(name string, dump []byte) {
			query := &common.DeviceQuery{}
			err := dc.decodeQuery(name, common.CmdDeviceStatus, dump, query)
			if err == nil && dc.commands != nil {
				err = dc.commands.Status(name, query)
			}
		})
		dc.scopeItem.SetScopeFunc(common.CmdRunAction, func(name string, dump []byte) {
			query := &common.DeviceQuery{}
			err := dc.decodeQuery(name, common.CmdRunAction, dump, query)
			if err == nil && dc.commands != nil {
				err = dc.commands.RunAction(name, query)
			}
		})
		dc.scopeItem.SetScopeFunc(common.CmdStopAction, func(name string, dump []byte) {
			query := &common.DeviceQuery{}
			err := dc.decodeQuery(name, common.CmdStopAction, dump, query)
			if err == nil && dc.commands != nil {
				err = dc.commands.StopAction(name, query)
			}
		})
	}
}

func (dc *DeviceClient) decodeQuery(name string, cmd string, dump []byte, query interface{}) error {
	if dc.log != nil {
		dc.log.Dump("DeviceClient dev:%s take cmd:%s, pack:%s", name, cmd, string(dump))
	}
	return json.Unmarshal(dump, query)
}
