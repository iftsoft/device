package duplex

import (
	"sync"
)

type ScopeFunc func(dump []byte)

type Transporter interface {
	SendPacket(pack *Packet, link string) error
	AddScopeItem(item *ScopeItem)
}

type ScopeItem struct {
	scopeId PacketScope
	handler map[string]ScopeFunc
	mutex   sync.RWMutex
}

func NewScopeItem(id PacketScope) *ScopeItem {
	si := ScopeItem{scopeId: id, handler: make(map[string]ScopeFunc)}
	return &si
}

func (si *ScopeItem) SetScopeFunc(name string, proc ScopeFunc) {
	si.mutex.Lock()
	defer si.mutex.Unlock()
	si.handler[name] = proc
}

func (si *ScopeItem) GetScopeFunc(name string) ScopeFunc {
	si.mutex.RLock()
	defer si.mutex.RUnlock()
	proc, ok := si.handler[name]
	if ok {
		return proc
	}
	return nil
}

type ScopeSet struct {
	store map[PacketScope]*ScopeItem
	mutex sync.RWMutex
}

func NewScopeSet() *ScopeSet {
	ss := ScopeSet{store: make(map[PacketScope]*ScopeItem)}
	return &ss
}

func (ss *ScopeSet) AddScope(scope *ScopeItem) {
	if scope == nil {
		return
	}
	ss.mutex.Lock()
	defer ss.mutex.Unlock()
	ss.store[scope.scopeId] = scope
}

func (ss *ScopeSet) GetScope(id PacketScope) *ScopeItem {
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()
	scope, ok := ss.store[id]
	if ok {
		return scope
	}
	return nil
}

func (ss *ScopeSet) GetScopeFunc(id PacketScope, name string) ScopeFunc {
	scope := ss.GetScope(id)
	if scope != nil {
		return scope.GetScopeFunc(name)
	}
	return nil
}
