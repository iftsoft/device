package duplex

import (
	"github.com/iftsoft/device/common"
	"sync"
)

const commandGreeting = "Greeting"

type ScopeFunc func(name string, dump []byte)

type Transporter interface {
	SendPacket(pack *Packet) error
}

type Dispatcher interface {
	EvalPacket(pack *Packet) error
}

type GreetingInfo struct {
	DevType   common.DevTypeMask	`json:"devType"`	// Implemented device types
	Supported common.DevScopeMask	`json:"supported"`	// Manager interfaces that driver supported
	Required  common.DevScopeMask	`json:"required"`	// Callback interfaces that driver required
}


type ScopeSet struct {
	store map[PacketScope]Dispatcher
	mutex sync.RWMutex
}

func NewScopeSet() *ScopeSet {
	ss := ScopeSet{store: make(map[PacketScope]Dispatcher)}
	return &ss
}

func (ss *ScopeSet) AddScope(id PacketScope, scope Dispatcher) {
	if scope == nil {
		return
	}
	ss.mutex.Lock()
	defer ss.mutex.Unlock()
	ss.store[id] = scope
}

func (ss *ScopeSet) GetScope(id PacketScope) Dispatcher {
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()
	scope, ok := ss.store[id]
	if ok {
		return scope
	}
	return nil
}

