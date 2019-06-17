package duplex

import (
	"fmt"
	"sync"
)

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
