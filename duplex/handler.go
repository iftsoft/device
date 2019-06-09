package duplex

import (
	"github.com/iftsoft/device/core"
	"net"
	"sync"
	"time"
)

type DuplexHandler struct {
	Duplex
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
		Config:   nil,
		scopeMap: nil,
	}
	dh.mngr = dh
	return dh
}

func (dh *DuplexHandler) Stop() {
	close(dh.done)
}

func (dh *DuplexHandler) Init(conn *net.TCPConn, name string, cfg *DuplexServerConfig) {
	dh.log = core.GetLogAgent(core.LogLevelTrace, name)
	dh.link.SetConnect(conn, dh.log)
	dh.Config = cfg
}

func (dh *DuplexHandler) NewPacket(pack *Packet) bool {
	proc := dh.scopeMap.GetScoreFunc(pack.Scope, pack.Command)
	if proc == nil {
		return false
	}
	proc(pack.Content)
	return true
}

func (dh *DuplexHandler) OnWriteError(err error) error {
	dh.log.Trace("DuplexHandler OnWriteError: %s", err)
	return nil
}

func (dh *DuplexHandler) OnReadError(err error) error {
	dh.log.Trace("DuplexHandler OnReadError: %s", err)
	return nil
}

func (dh *DuplexHandler) OnTimerTick(tm time.Time) {
	dh.log.Trace("DuplexHandler OnTimerTick: %s", tm.Format(time.StampMilli))
}

func (dh *DuplexHandler) HandlerLoop(wg *sync.WaitGroup) {
	defer wg.Done()
	dh.readingLoop()
}
