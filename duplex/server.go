package duplex

import (
	"fmt"
	"github.com/iftsoft/device/core"
	"net"
	"sync"
)

type DuplexServerConfig struct {
	Port int32 `yaml:"port"`
}

type HandleSet struct {
	store map[string]*DuplexHandler
	count uint32
	mutex sync.RWMutex
}

func NewHandleSet() HandleSet {
	hs := HandleSet{store: make(map[string]*DuplexHandler), count: 0}
	return hs
}

func (hs *HandleSet) AddHandler(handle *DuplexHandler) string {
	hs.mutex.Lock()
	defer hs.mutex.Unlock()
	hs.count++
	name := fmt.Sprintf("link_%d", hs.count)
	hs.store[name] = handle
	return name
}

func (hs *HandleSet) GetHandler(name string) *DuplexHandler {
	hs.mutex.RLock()
	defer hs.mutex.RUnlock()
	handle, ok := hs.store[name]
	if ok {
		return handle
	}
	return nil
}

func (hs *HandleSet) DelHandler(name string) {
	hs.mutex.Lock()
	defer hs.mutex.Unlock()
	delete(hs.store, name)
}

type DuplexServer struct {
	Config   *DuplexServerConfig
	listener net.Listener
	scopeMap ScopeSet
	handles  HandleSet
	log      *core.LogAgent
}

func (ds *DuplexServer) Listen(addr *net.TCPAddr) bool {
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		ds.log.Error("Unable to listen on port %s\n", addr.String())
		return false
	}
	ds.log.Info("Listen on %s", listener.Addr().String())
	for {
		ds.log.Trace("Accept a connection request.")
		conn, err := listener.AcceptTCP()
		if err != nil {
			ds.log.Error("Failed accepting a connection request:", err)
			continue
		}
		ds.log.Trace("Handle incoming messages.")
		go ds.handleMessages(conn)
	}
	return true
}

func (ds *DuplexServer) handleMessages(conn *net.TCPConn) {
	hand := GetDuplexHandler()
	name := ds.handles.AddHandler(hand)
	hand.Init(conn, name, ds.Config)
	//	hand.HandlerLoop()
}
