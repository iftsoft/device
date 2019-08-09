package duplex

import (
	"encoding/json"
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"net"
	"time"
)

type DuplexHandler struct {
	Duplex
	HndName  string
	DevName  string
	Config   *DuplexServerConfig
	scopeMap *ScopeSet
}

func GetDuplexHandler() *DuplexHandler {
	dh := &DuplexHandler{
		Duplex: Duplex{
			link: LinkHolder{},
			done: make(chan struct{}),
			log:  nil,
		},
		HndName:  "",
		DevName:  "",
		Config:   nil,
		scopeMap: nil,
	}
	dh.mngr = dh
	return dh
}

func (dh *DuplexHandler) Stop() {
	dh.link.CloseConnect()
	close(dh.done)
}

func (dh *DuplexHandler) Init(conn *net.TCPConn, cfg *DuplexServerConfig, scopes *ScopeSet) {
	_ = conn.SetNoDelay(true)
	_ = conn.SetKeepAlive(true)
	_ = conn.SetKeepAlivePeriod(5 * time.Second)
	dh.log = core.GetLogAgent(core.LogLevelTrace, dh.HndName)
	dh.link.SetConnect(conn, dh.log)
	dh.Config = cfg
	dh.scopeMap = scopes
}

func (dh *DuplexHandler) HandlerLoop(hs *HandleSet) {
	defer hs.DelHandler(dh.HndName, dh.DevName)
	defer dh.link.CloseConnect()
	// Get client greeting info
	err := dh.ReadGreeting()
	if err != nil {
		dh.log.Error("DuplexHandler ReadGreeting error: %s", err)
		return
	}
	hs.SetHandlerDevice(dh.HndName, dh.DevName)
	dh.log.Info("DuplexHandler %s started for device %s", dh.HndName, dh.DevName)
	defer dh.log.Info("DuplexHandler %s stopped for device %s", dh.HndName, dh.DevName)

	go dh.readingLoop()
	dh.waitingLoop()
}

// Implementation of DuplexManager interface
func (dh *DuplexHandler) NewPacket(pack *Packet) bool {
	//	dh.log.Trace("DuplexHandler NewPacket: %+v", pack)
	dh.log.Trace("DuplexHandler NewPacket dev:%s, cmd:%s, dump:%s", pack.DevName, pack.Command, string(pack.Content))
	proc := dh.scopeMap.GetScopeFunc(pack.Scope, pack.Command)
	if proc == nil {
		dh.log.Trace("DuplexHandler NewPacket: Unknown command - %s", pack.Command)
		return false
	}
	proc(pack.DevName, pack.Content)
	return true
}

func (dh *DuplexHandler) OnWriteError(err error) error {
	dh.log.Trace("DuplexHandler OnWriteError: %s", err)
	dh.Stop()
	return err
}

func (dh *DuplexHandler) OnReadError(err error) error {
	dh.log.Trace("DuplexHandler OnReadError: %s", err)
	dh.Stop()
	return err
}

func (dh *DuplexHandler) OnTimerTick(tm time.Time) {
	dh.log.Trace("DuplexHandler OnTimerTick: %s", tm.Format(time.StampMilli))
	//	dh.SendRequest()
}

// Implementation of Transporter interface
func (dh *DuplexHandler) SendPacket(pack *Packet) error {
	return dh.WritePacket(pack)
}

func (dh *DuplexHandler) SendRequest() {
	query := &common.SystemQuery{}
	query.DevName = "default"
	data, _ := json.Marshal(query)
	pack := NewPacket(ScopeSystem, "default", "Inform", data)
	err := dh.WritePacket(pack)
	if err != nil {
		dh.log.Error("DuplexServer WritePacket error: %s", err)
	}
}

func (dh *DuplexHandler) ReadGreeting() error {
	conn := dh.link.GetConnect()
	if conn == nil {
		return errors.New("duplex DialTCP conn is nil")
	}
	pack, err := conn.ReadPacket()
	if err != nil {
		return err
	}
	if pack.Command != commandGreeting {
		return errors.New("packet is not Greeting")
	}
	dh.DevName = pack.DevName
	return nil
}
