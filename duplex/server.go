package duplex

import (
	"errors"
	"fmt"
	"github.com/iftsoft/device/core"
	"net"
)

type DuplexServerConfig struct {
	Port int32 `yaml:"port"`
}

type DuplexServer struct {
	config   *DuplexServerConfig
	listener *net.TCPListener
	scopeMap *ScopeSet
	handles  *HandleSet
	log      *core.LogAgent
	exit     bool
}

func NewDuplexServer(config *DuplexServerConfig, log *core.LogAgent) *DuplexServer {
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

// Implementation of Transporter interface
func (ds *DuplexServer) SendPacket(pack *Packet, link string) error {
	hnd := ds.handles.GetHandler(link)
	if hnd != nil {
		return hnd.WritePacket(pack)
	}
	return errors.New("no such link")
}

func (ds *DuplexServer) AddScopeItem(item *ScopeItem) {
	if item != nil {
		ds.scopeMap.AddScope(item)
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
	ds.log.Info("Stop listening on server.")
	ds.exit = true
	_ = ds.listener.Close()
	ds.log.Trace("Closing all connections...")
	ds.handles.StopAllHandlers()
	ds.log.Trace("Waiting for running handlers.")
	ds.handles.WaitAllHandlers()
	ds.log.Info("All connections are closed.")
}

func (ds *DuplexServer) listenLoop() {
	if ds.listener == nil {
		return
	}
	ds.log.Info("Start listen on %s", ds.listener.Addr().String())
	for {
		conn, err := ds.listener.AcceptTCP()
		ds.log.Trace("Accept a connection request.")
		if err != nil {
			if ds.exit == true {
				break
			} else {
				ds.log.Error("Failed accepting a connection request:", err)
				continue
			}
		}
		ds.log.Trace("Handle incoming messages.")
		go ds.handleMessages(conn)
	}
	ds.log.Info("Stop listen on %s", ds.listener.Addr().String())
}

func (ds *DuplexServer) handleMessages(conn *net.TCPConn) {
	hand := GetDuplexHandler()
	name := ds.handles.AddHandler(hand)
	hand.Init(conn, name, ds.config, ds.scopeMap)
	hand.HandlerLoop(ds.handles)
}
