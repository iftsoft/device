package duplex

import (
	"sync"
)

type ScopeFunc func(dump []byte)

type ScopeItem struct {
	handler map[string]ScopeFunc
	mutex   sync.RWMutex
}

func NewScoreItem() *ScopeItem {
	si := ScopeItem{handler: make(map[string]ScopeFunc)}
	return &si
}

func (si *ScopeItem) SetScoreFunc(name string, proc ScopeFunc) {
	si.mutex.Lock()
	defer si.mutex.Unlock()
	si.handler[name] = proc
}

func (si *ScopeItem) GetScoreFunc(name string) ScopeFunc {
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

func NewScoreSet() *ScopeSet {
	ss := ScopeSet{store: make(map[PacketScope]*ScopeItem)}
	return &ss
}

func (ss *ScopeSet) SetScore(id PacketScope, score *ScopeItem) {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()
	ss.store[id] = score
}

func (ss *ScopeSet) GetScore(id PacketScope) *ScopeItem {
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()
	scope, ok := ss.store[id]
	if ok {
		return scope
	}
	return nil
}

func (ss *ScopeSet) GetScoreFunc(id PacketScope, name string) ScopeFunc {
	scope := ss.GetScore(id)
	if scope != nil {
		return scope.GetScoreFunc(name)
	}
	return nil
}
