package proxy

import (
	"encoding/json"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/duplex"
)

type DeviceClient struct {
	scopeItem *duplex.ScopeItem
	transport duplex.Transporter
	commands  common.DeviceManager
	log       *core.LogAgent
}

func NewDeviceClient() *DeviceClient {
	dc := DeviceClient{
		scopeItem: duplex.NewScopeItem(duplex.ScopeDevice),
		transport: nil,
		commands:  nil,
		log:       nil,
	}
	return &dc
}

func (dc *DeviceClient) Init(trans duplex.Transporter, command common.DeviceManager, log *core.LogAgent) {
	dc.log = log
	dc.transport = trans
	dc.commands = command
	// init scope functions
	if dc.scopeItem != nil {
		dc.scopeItem.SetScopeFunc(common.CmdDeviceCancel, func(name string, dump []byte) {
			query, err := dc.decodeQuery(name, common.CmdDeviceCancel, dump)
			if err == nil && dc.commands != nil {
				err = dc.commands.Cancel(name, query)
			}
		})
		dc.scopeItem.SetScopeFunc(common.CmdDeviceReset, func(name string, dump []byte) {
			query, err := dc.decodeQuery(name, common.CmdDeviceReset, dump)
			if err == nil && dc.commands != nil {
				err = dc.commands.Reset(name, query)
			}
		})
		dc.scopeItem.SetScopeFunc(common.CmdDeviceStatus, func(name string, dump []byte) {
			query, err := dc.decodeQuery(name, common.CmdDeviceStatus, dump)
			if err == nil && dc.commands != nil {
				err = dc.commands.Status(name, query)
			}
		})
	}
}

func (dc *DeviceClient) decodeQuery(name string, cmd string, dump []byte) (query *common.DeviceQuery, err error) {
	if dc.log != nil {
		dc.log.Dump("DeviceClient for dev:%s get cmd:%s, pack:%s", name, cmd, string(dump))
	}
	query = &common.DeviceQuery{}
	err = json.Unmarshal(dump, query)
	return query, err
}

func (dc *DeviceClient) GetScopeItem() *duplex.ScopeItem {
	return dc.scopeItem
}

// Implementation of common.DeviceCallback
func (dc *DeviceClient) DeviceReply(name string, reply *common.DeviceReply) error {
	return dc.encodeReply(name, common.CmdDeviceReply, reply)
}
func (dc *DeviceClient) ExecuteError(name string, reply *common.DeviceError) error {
	return dc.encodeReply(name, common.CmdExecuteError, reply)
}
func (dc *DeviceClient) StateChanged(name string, reply *common.DeviceState) error {
	return dc.encodeReply(name, common.CmdStateChanged, reply)
}
func (dc *DeviceClient) ActionPrompt(name string, reply *common.DevicePrompt) error {
	return dc.encodeReply(name, common.CmdActionPrompt, reply)
}

func (dc *DeviceClient) encodeReply(name string, cmd string, reply interface{}) error {
	dump, err := json.Marshal(reply)
	if err != nil {
		return err
	}
	if dc.log != nil {
		dc.log.Dump("DeviceClient dev:%s put cmd:%s pack:%s", name, cmd, string(dump))
	}
	pack := duplex.NewPacket(duplex.ScopeDevice, name, cmd, dump)
	if dc.transport != nil {
		err = dc.transport.SendPacket(pack)
	}
	return err
}
