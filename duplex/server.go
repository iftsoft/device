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
	wg    sync.WaitGroup
}

func NewHandleSet() *HandleSet {
	hs := HandleSet{store: make(map[string]*DuplexHandler), count: 0}
	return &hs
}

func (hs *HandleSet) AddHandler(handle *DuplexHandler) string {
	if handle == nil {
		return ""
	}
	hs.wg.Add(1)
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
	hs.wg.Done()
	hs.mutex.Lock()
	defer hs.mutex.Unlock()
	delete(hs.store, name)
}

func (hs *HandleSet) StopAllHandlers() {
	hs.mutex.RLock()
	defer hs.mutex.RUnlock()
	for _, hnd := range hs.store {
		if hnd != nil {
			hnd.Stop()
		}
	}
}

func (hs *HandleSet) WaitAllHandlers() {
	hs.wg.Wait()
}

type DuplexServer struct {
	Config   *DuplexServerConfig
	listener *net.TCPListener
	scopeMap *ScopeSet
	handles  *HandleSet
	log      *core.LogAgent
}

func NewDuplexServer() *DuplexServer {
	ds := DuplexServer{
		Config:   nil,
		listener: nil,
		scopeMap: NewScoreSet(),
		handles:  NewHandleSet(),
		log:      core.GetLogAgent(core.LogLevelTrace, "Listener"),
	}
	return &ds
}

func (ds *DuplexServer) Listen(addr *net.TCPAddr) bool {
	var err error
	ds.listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		ds.log.Error("Unable to listen on port %s\n", addr.String())
		return false
	}
	ds.log.Info("Listen on %s", ds.listener.Addr().String())
	for {
		ds.log.Trace("Accept a connection request.")
		conn, err := ds.listener.AcceptTCP()
		if err != nil {
			ds.log.Error("Failed accepting a connection request:", err)
			_, err = ds.listener.SyscallConn()
			if err != nil {
				break
			} else {
				continue
			}
		}
		ds.log.Trace("Handle incoming messages.")
		go ds.handleMessages(conn)
	}
	return true
}

func (ds *DuplexServer) StopListen() {
	ds.log.Info("Stop listening on server.")
	_ = ds.listener.Close()
	ds.log.Trace("Closing all connections...")
	ds.handles.StopAllHandlers()
	ds.log.Trace("Waiting for running handlers.")
	ds.handles.WaitAllHandlers()
	ds.log.Info("All connections are closed.")
}

func (ds *DuplexServer) handleMessages(conn *net.TCPConn) {
	hand := GetDuplexHandler()
	name := ds.handles.AddHandler(hand)
	hand.Init(conn, name, ds.Config)
	hand.HandlerLoop(ds.handles)
}
