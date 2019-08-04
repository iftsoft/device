package duplex

import (
	"encoding/json"
	"fmt"
	"github.com/iftsoft/device/core"
	"net"
	"time"
)

type DuplexClientConfig struct {
	Port    int32  `yaml:"port"`
	DevName string `yaml:"devName"`
}

type DuplexClient struct {
	Duplex
	config   *DuplexClientConfig
	scopeMap *ScopeSet
}

func NewDuplexClient(cfg *DuplexClientConfig) *DuplexClient {
	dc := &DuplexClient{
		Duplex: Duplex{
			link: LinkHolder{},
			done: make(chan struct{}),
			log:  core.GetLogAgent(core.LogLevelTrace, "Duplex"),
		},
		config:   cfg,
		scopeMap: NewScopeSet(),
	}
	dc.mngr = dc
	return dc
}

func (dc *DuplexClient) Start() {
	dc.log.Info("Starting client engine")
	go dc.clientLoop(dc.config.Port)
}

func (dc *DuplexClient) Stop() {
	dc.log.Info("Stopping client engine")
	close(dc.done)
}

func (dc *DuplexClient) AddScopeItem(item *ScopeItem) {
	if item != nil {
		dc.scopeMap.AddScope(item)
	}
}

// Implementation of Transporter interface
func (dc *DuplexClient) SendPacket(pack *Packet) error {
	return dc.WritePacket(pack)
}

// Implementation of DuplexManager interface
func (dc *DuplexClient) NewPacket(pack *Packet) bool {
	proc := dc.scopeMap.GetScopeFunc(pack.Scope, pack.Command)
	if proc == nil {
		dc.log.Trace("DuplexClient NewPacket: Unknown command - %s", pack.Command)
		return false
	}
	proc(pack.Content)
	return true
}

func (dc *DuplexClient) OnWriteError(err error) error {
	dc.log.Trace("DuplexClient OnWriteError: %s", err)
	dc.link.CloseConnect()
	return nil
}

func (dc *DuplexClient) OnReadError(err error) error {
	dc.log.Trace("DuplexClient OnReadError: %s", err)
	dc.link.CloseConnect()
	return nil
}

func (dc *DuplexClient) OnTimerTick(tm time.Time) {
	dc.log.Trace("DuplexClient OnTimerTick: %s", tm.Format(time.StampMilli))
	conn := dc.link.GetConnect()
	if conn == nil {
		dc.log.Error("DuplexClient DialTCP conn is nil")
		err := dc.connectToServer(dc.config.Port)
		if err != nil {
			dc.log.Error("DuplexClient Dial error: %s", err)
		}
	} else {
		//		dc.log.Info("Dialling connect %+v", conn)
		//		dc.SendRequest()
	}
}

func (dc *DuplexClient) clientLoop(port int32) {
	_ = dc.connectToServer(port)
	defer dc.link.CloseConnect()
	dc.waitingLoop()
}

func (dc *DuplexClient) connectToServer(port int32) error {
	err := dc.dialToAddress(port)
	if err != nil {
		return err
	}
	err = dc.sendGreeting()
	if err != nil {
		dc.link.CloseConnect()
		return err
	}
	go dc.readingLoop()
	return nil
}

func (dc *DuplexClient) dialToAddress(port int32) error {
	servAddr := fmt.Sprintf("localhost:%d", port)
	dc.log.Info("Dialling to %s", servAddr)
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		dc.log.Error("DuplexClient ResolveTCPAddr: %s", err)
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		dc.log.Error("DuplexClient DialTCP error: %s", err)
		return err
	}
	if conn == nil {
		dc.log.Error("DuplexClient DialTCP conn is nil")
		return err
	}
	dc.log.Info("Dialling connect %+v", conn)
	err = conn.SetNoDelay(true)
	err = conn.SetKeepAlive(true)
	err = conn.SetKeepAlivePeriod(5 * time.Second)
	dc.link.SetConnect(conn, dc.log)
	return nil
}

func (dc *DuplexClient) sendGreeting() error {
	hello := Greeting{}
	if dc.config != nil {
		hello.DevName = dc.config.DevName
		dc.log.Trace("DuplexClient SendGreeting for device: %s", hello.DevName)
	}
	dump, err := json.Marshal(&hello)
	if err != nil {
		dc.log.Error("DuplexClient SendGreeting error: %s", err)
		return err
	}
	pack := NewRequest(ScopeSystem)
	pack.Command = commandGreeting
	pack.Content = dump
	err = dc.WritePacket(pack)
	if err != nil {
		dc.log.Error("DuplexClient WritePacket error: %s", err)
	}
	return err
}
