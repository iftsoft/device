package duplex

import (
	"github.com/iftsoft/device/core"
	"net"
	"time"
)

type DuplexHandler struct {
	Duplex
	HndName  string
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
	go dh.readingLoop()
	dh.waitingLoop()
}

func (dh *DuplexHandler) NewPacket(pack *Packet) bool {
	//	dh.log.Trace("DuplexHandler NewPacket: %+v", pack)
	//proc := dh.scopeMap.GetScoreFunc(pack.Scope, pack.Command)
	//if proc == nil {
	//	return false
	//}
	//proc(pack.Content)
	//back := NewResponse(pack)
	//back.Command = "ServerResponse"
	//back.Content = []byte ("Server Response")
	//dh.WritePacket(back)
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

func (dh *DuplexHandler) SendRequest() {
	pack := NewRequest(ScopeSystem)
	pack.Command = "SeverRequest"
	pack.Content = []byte("Server Request")
	err := dh.WritePacket(pack)
	if err != nil {
		dh.log.Error("DuplexServer WritePacket error: %s", err)
	}
}
