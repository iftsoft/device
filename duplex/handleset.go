package duplex

import (
	"fmt"
	"sync"
)

type ClientManager interface {
	OnClientStarted(name string, info *GreetingInfo)
	OnClientStopped(name string)
}

type HandleSet struct {
	store   map[string]*DuplexHandler
	names   map[string]string
	manager ClientManager
	count   uint32
	mutex   sync.RWMutex
	wg      sync.WaitGroup
}

func NewHandleSet() *HandleSet {
	hs := HandleSet{
		store:   make(map[string]*DuplexHandler),
		names:   make(map[string]string),
		manager: nil,
		count:   0,
	}
	return &hs
}

func (hs *HandleSet) AddHandler() *DuplexHandler {
	handle := GetDuplexHandler()
	if handle == nil {
		return nil
	}

	hs.wg.Add(1)
	hs.mutex.Lock()
	defer hs.mutex.Unlock()

	hs.count++
	link := fmt.Sprintf("Link_%d", hs.count)
	handle.HndName = link
	hs.store[link] = handle
	return handle
}

func (hs *HandleSet) SetHandlerDevice(link, name string, info *GreetingInfo) {
	hs.mutex.Lock()
	hs.names[name] = link
	hs.mutex.Unlock()
	if hs.manager != nil {
		hs.manager.OnClientStarted(name, info)
	}
}

func (hs *HandleSet) GetHandler(name string) *DuplexHandler {
	hs.mutex.RLock()
	defer hs.mutex.RUnlock()
	link, ok1 := hs.names[name]
	if ok1 {
		handle, ok := hs.store[link]
		if ok {
			return handle
		}
	}
	return nil
}

func (hs *HandleSet) DelHandler(link, name string) {
	if hs.manager != nil {
		hs.manager.OnClientStopped(name)
	}
	hs.wg.Done()
	hs.mutex.Lock()
	defer hs.mutex.Unlock()
	delete(hs.store, link)
	delete(hs.names, name)
}

func (hs *HandleSet) StopAllHandlers() {
	hs.mutex.RLock()
	defer hs.mutex.RUnlock()
	for _, hnd := range hs.store {
		if hnd != nil {
			hnd.StopHandle(&hs.wg)
		}
	}
	hs.wg.Wait()
}
