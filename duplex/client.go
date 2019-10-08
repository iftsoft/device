package duplex

import (
	"fmt"
	"github.com/iftsoft/device/core"
	"io"
	"net"
	"sync"
	"time"
)

type ClientConfig struct {
	Port    int32  `yaml:"port"`
	DevName string `yaml:"dev_name"`
}

func GetDefaultClientConfig() *ClientConfig {
	srvCfg := &ClientConfig{
		Port:    DuplexPort,
		DevName: "",
	}
	return srvCfg
}

func (cfg *ClientConfig) String() string {
	if cfg == nil { return "" }
	str := fmt.Sprintf("\nDuplex client config: "+
		"Port = %d, DevName = %s.",
		cfg.Port, cfg.DevName)
	return str
}

type DuplexClient struct {
	Duplex
	config   *ClientConfig
	scopeMap *ScopeSet
}

func NewDuplexClient(cfg *ClientConfig) *DuplexClient {
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

func (dc *DuplexClient) StartClient(wg *sync.WaitGroup) {
	wg.Add(1)
	dc.log.Info("Starting client engine")
	go dc.clientLoop(wg, dc.config.Port)
}

func (dc *DuplexClient) StopClient(wg *sync.WaitGroup) {
	wg.Done()
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
func (dc *DuplexClient) OnNewPacket(pack *Packet) bool {
	dc.log.Trace("DuplexClient OnNewPacket dev:%s, cmd:%s", pack.DevName, pack.Command)
	proc := dc.scopeMap.GetScopeFunc(pack.Scope, pack.Command)
	if proc == nil {
		dc.log.Warn("DuplexClient OnNewPacket: Unknown command - %s", pack.Command)
		return false
	}
	proc(pack.DevName, pack.Content)
	return true
}

func (dc *DuplexClient) OnWriteError(err error) error {
	dc.log.Debug("DuplexClient OnWriteError: %s", err)
	dc.link.CloseConnect()
	return nil
}

func (dc *DuplexClient) OnReadError(err error) error {
	dc.log.Debug("DuplexClient OnReadError: %s", err)
	dc.link.CloseConnect()
	return io.EOF
}

func (dc *DuplexClient) OnTimerTick(tm time.Time) {
	dc.log.Trace("DuplexClient OnTimerTick: %s", tm.Format(time.StampMilli))
	conn := dc.link.GetConnect()
	if conn == nil {
		dc.log.Warn("DuplexClient DialTCP connect is not open. Trying to dial...")
		err := dc.connectToServer(dc.config.Port)
		if err != nil {
			dc.log.Error("DuplexClient Dial error: %s", err)
		}
	} else {
		//		dc.log.Info("Dialling connect %+v", conn)
		//		dc.SendRequest()
	}
}

func (dc *DuplexClient) clientLoop(wg *sync.WaitGroup, port int32) {
	err := dc.connectToServer(port)
	if err != nil {
		return
	}
	defer dc.link.CloseConnect()
	go dc.readingLoop(wg)
	dc.waitingLoop(wg)
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
	return nil
}

func (dc *DuplexClient) dialToAddress(port int32) error {
	servAddr := fmt.Sprintf("localhost:%d", port)
	dc.log.Trace("Dialling to %s", servAddr)
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		dc.log.Error("DuplexClient ResolveTCPAddr: %s", err)
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		dc.log.Warn("DuplexClient DialTCP: %s", err)
		return err
	}
	if conn == nil {
		dc.log.Error("DuplexClient DialTCP conn is nil")
		return err
	}
	//	dc.log.Info("Dialling connect %+v", conn)
	err = conn.SetNoDelay(true)
	err = conn.SetKeepAlive(true)
	err = conn.SetKeepAlivePeriod(5 * time.Second)
	dc.link.SetConnect(conn, dc.log)
	return nil
}

func (dc *DuplexClient) sendGreeting() error {
	name := ""
	if dc.config != nil {
		name = dc.config.DevName
	}
	dc.log.Info("DuplexClient SendGreeting for device: %s", name)
	pack := NewPacket(ScopeSystem, name, commandGreeting, nil)
	err := dc.WritePacket(pack)
	if err != nil {
		dc.log.Error("DuplexClient WritePacket error: %s", err)
	}
	return err
}
