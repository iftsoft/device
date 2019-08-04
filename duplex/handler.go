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
	Greeting *Greeting
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
		Greeting: &Greeting{},
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

func (dh *DuplexHandler) Init(conn *net.TCPConn, name string, cfg *DuplexServerConfig, scopes *ScopeSet) {
	_ = conn.SetNoDelay(true)
	_ = conn.SetKeepAlive(true)
	_ = conn.SetKeepAlivePeriod(5 * time.Second)
	dh.log = core.GetLogAgent(core.LogLevelTrace, name)
	dh.link.SetConnect(conn, dh.log)
	dh.HndName = name
	dh.Config = cfg
	dh.scopeMap = scopes
}

func (dh *DuplexHandler) HandlerLoop(hs *HandleSet) {
	defer hs.DelHandler(dh.HndName)
	defer dh.link.CloseConnect()
	// Get client greeting info
	err := dh.ReadGreeting()
	if err != nil {
		dh.log.Error("DuplexHandler ReadGreeting error: %s", err)
		return
	}
	go dh.readingLoop()
	dh.waitingLoop()
}

// Implementation of DuplexManager interface
func (dh *DuplexHandler) NewPacket(pack *Packet) bool {
	//	dh.log.Trace("DuplexHandler NewPacket: %+v", pack)
	dh.log.Trace("DuplexHandler NewPacket cmd:%s, dump:%s", pack.Command, string(pack.Content))
	proc := dh.scopeMap.GetScopeFunc(pack.Scope, pack.Command)
	if proc == nil {
		return false
	}
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
	dh.SendRequest()
}

// Implementation of Transporter interface
func (dh *DuplexHandler) SendPacket(pack *Packet) error {
	return dh.WritePacket(pack)
}

func (dh *DuplexHandler) SendRequest() {
	query := &common.SystemQuery{}
	query.DevName = "default"
	pack := NewRequest(ScopeSystem)
	pack.Command = "Inform"
	pack.Content, _ = json.Marshal(query)
	err := dh.WritePacket(pack)
	if err != nil {
		dh.log.Error("DuplexServer WritePacket error: %s", err)
	}
}

func (dh *DuplexHandler) ReadGreeting() error {
	conn := dh.link.GetConnect()
	if conn != nil {
		pack, err := conn.ReadPacket()
		if err != nil {
			return err
		} else if pack != nil {
			if pack.Command != commandGreeting {
				return errors.New("packet is not Greeting")
			}
			err = json.Unmarshal(pack.Content, dh.Greeting)
			if err != nil {
				return err
			}
		}
	} else {
		return errors.New("duplex DialTCP conn is nil")
	}
	return nil
}
