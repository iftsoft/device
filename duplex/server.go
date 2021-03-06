package duplex

import (
	"fmt"
	"github.com/iftsoft/device/core"
	"net"
)

type ServerManager interface {
	AddDispatcher(id PacketScope, scope Dispatcher)
	GetTransporter(name string) Transporter
}

type ServerConfig struct {
	Port int32 `yaml:"port"`
}

func GetDefaultServerConfig() *ServerConfig {
	srvCfg := &ServerConfig{
		Port: DuplexPort,
	}
	return srvCfg
}

func (cfg *ServerConfig) String() string {
	if cfg == nil { return "" }
	str := fmt.Sprintf("\nDuplex server config: "+
		"Port = %d.",
		cfg.Port)
	return str
}

type DuplexServer struct {
	config   *ServerConfig
	listener *net.TCPListener
	scopeMap *ScopeSet
	handles  *HandleSet
	log      *core.LogAgent
	exit     bool
}

func NewDuplexServer(config *ServerConfig, log *core.LogAgent) *DuplexServer {
	ds := DuplexServer{
		config:   config,
		listener: nil,
		scopeMap: NewScopeSet(),
		handles:  NewHandleSet(),
		log:      log,
		exit:     false,
	}
	return &ds
}

func (ds *DuplexServer) SetClientManager(manager ClientManager) {
	ds.handles.manager = manager
}

// Implementation of ServerManager interface
func (ds *DuplexServer) GetTransporter(name string) Transporter {
	hnd := ds.handles.GetHandler(name)
	return hnd
}

func (ds *DuplexServer) AddDispatcher(id PacketScope, scope Dispatcher) {
	if scope != nil {
		ds.scopeMap.AddScope(id, scope)
	}
}

func (ds *DuplexServer) StartListen() error {
	servAddr := fmt.Sprintf("localhost:%d", ds.config.Port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		ds.log.Error("DuplexClient ResolveTCPAddr: %s", err)
		return err
	}
	ds.listener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		ds.log.Error("Unable to listen on port %s\n", tcpAddr.String())
		return err
	}
	go ds.listenLoop()
	return nil
}

func (ds *DuplexServer) StopListen() {
	ds.log.Info("SysStop listening on server.")
	ds.exit = true
	_ = ds.listener.Close()
	ds.log.Trace("Closing all connections...")
	ds.handles.StopAllHandlers()
	ds.log.Info("All connections are closed.")
}

func (ds *DuplexServer) listenLoop() {
	if ds.listener == nil {
		return
	}
	ds.log.Info("SysStart listen on %s", ds.listener.Addr().String())
	for {
		conn, err := ds.listener.AcceptTCP()
		if err != nil {
			if ds.exit == true {
				break
			} else {
				ds.log.Error("Failed accepting a connection request:", err)
				continue
			}
		}
		ds.log.Debug("Accept a connection request.")
		go ds.handleMessages(conn)
	}
	ds.log.Info("SysStop listen on %s", ds.listener.Addr().String())
}

func (ds *DuplexServer) handleMessages(conn *net.TCPConn) {
	hand := ds.handles.AddHandler()
	if hand != nil {
		hand.Init(conn, ds.config, ds.scopeMap)
		hand.StartHandler(ds.handles)
	}
}
