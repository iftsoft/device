package duplex

import (
	"fmt"
	"github.com/iftsoft/device/core"
	"net"
	"time"
)

type DuplexClientConfig struct {
	Port int32 `yaml:"port"`
}

type DuplexClient struct {
	Duplex
	Config   *DuplexClientConfig
	scopeMap ScopeSet
}

func GetDuplexClient() *DuplexClient {
	dc := &DuplexClient{
		Duplex: Duplex{
			link: LinkHolder{},
			done: make(chan struct{}),
			log:  nil,
		},
		Config:   nil,
		scopeMap: NewScoreSet(),
	}
	dc.mngr = dc
	return dc
}

func (dc *DuplexClient) Stop() {
	close(dc.done)
}

func (dc *DuplexClient) Init(cfg *DuplexClientConfig) {
	dc.log = core.GetLogAgent(core.LogLevelTrace, "Duplex")
	dc.Config = cfg
}

func (dc *DuplexClient) NewPacket(pack *Packet) bool {
	proc := dc.scopeMap.GetScoreFunc(pack.Scope, pack.Command)
	if proc == nil {
		return false
	}
	proc(pack.Content)
	return true
}

func (dc *DuplexClient) OnWriteError(err error) error {
	dc.log.Trace("DuplexClient OnWriteError: %s", err)
	return nil
}

func (dc *DuplexClient) OnReadError(err error) error {
	dc.log.Trace("DuplexClient OnReadError: %s", err)
	return nil
}

func (dc *DuplexClient) OnTimerTick(tm time.Time) {
	dc.log.Trace("DuplexClient OnTimerTick: %s", tm.Format(time.StampMilli))
	conn := dc.link.GetConnect()
	if conn == nil {
		dc.dialToAddress(dc.Config.Port)
	}
}

func (dc *DuplexClient) dialToAddress(port int32) error {
	servAddr := fmt.Sprintf("localhost:%d", port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		dc.log.Error("DuplexClient ResolveTCPAddr: %s", err)
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	dc.link.SetConnect(conn, dc.log)
	return nil
}

func (dc *DuplexClient) ClientLoop() {
	dc.readingLoop()
}
