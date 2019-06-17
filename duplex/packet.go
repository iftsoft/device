package duplex

import (
	"errors"
	"github.com/iftsoft/device/core"
)

const (
	packetVersion    PacketVersion = 0
	packetHeaderSize int           = 16
)

type PacketType byte
type PacketScope byte
type PacketVersion byte

const (
	PacketNotify PacketType = iota
	PackerRequest
	PacketResponse

//	PacketCallback
//	PacketBackword
)

const (
	ScopeSystem PacketScope = iota
	ScopeConfig
	ScopeDevice
)

type Packet struct {
	Version PacketVersion
	Type    PacketType
	Scope   PacketScope
	Counter uint32
	Options uint32
	Command string
	Content []byte
}

func NewRequest(scope PacketScope) *Packet {
	p := Packet{
		Version: packetVersion,
		Type:    PackerRequest,
		Scope:   scope,
		Counter: 0,
		Options: 0,
	}
	return &p
}

func NewResponse(query *Packet) *Packet {
	p := Packet{
		Version: packetVersion,
		Type:    PacketResponse,
		Scope:   query.Scope,
		Counter: query.Counter,
		Options: 0,
		Command: query.Command,
	}
	return &p
}

func (p *Packet) Print(log *core.LogAgent, text string) {
	if p != nil && log != nil {
		log.Trace("%s packet Type:%d, Scope:%d, Command:%s, Counter:%d, Data len:%d",
			text, p.Type, p.Scope, p.Command, p.Counter, len(p.Content))
	}
}

func (p *Packet) Encode() []byte {
	head := make([]byte, packetHeaderSize)
	task := []byte(p.Command)
	len1 := len(task)
	if len1 > 250 {
		task = append(task[:250], byte(0))
		len1 = len(task)
	}
	size := len(p.Content)
	head[0] = byte(p.Version)
	head[1] = byte(p.Type)
	head[2] = byte(p.Scope)
	head[3] = byte(len1)
	head[4], head[5], head[6], head[7] = UintToBytes(p.Options)
	head[8], head[9], head[10], head[11] = UintToBytes(p.Counter)
	head[12], head[13], head[14], head[15] = UintToBytes(uint32(size))
	dump := append(head, task...)
	dump = append(dump, p.Content...)
	return dump
}

func (p *Packet) Decode(dump []byte) error {
	err := errors.New("wrong packet size")
	full := len(dump)
	if full < packetHeaderSize {
		return err
	}
	head := dump[:packetHeaderSize]
	len1 := int(head[3])
	if full < packetHeaderSize+len1 {
		return err
	}
	task := dump[packetHeaderSize : packetHeaderSize+len1]
	p.Version = PacketVersion(head[0])
	p.Type = PacketType(head[1])
	p.Scope = PacketScope(head[2])
	p.Options = BytesToUint(head[4], head[5], head[6], head[7])
	p.Counter = BytesToUint(head[8], head[9], head[10], head[11])
	p.Command = string(task)
	size := int(BytesToUint(head[12], head[13], head[14], head[15]))
	if full < packetHeaderSize+len1+size {
		return err
	}
	if size > 0 {
		p.Content = dump[packetHeaderSize+len1:]
	}
	return nil
}

func UintToBytes(val uint32) (b0, b1, b2, b3 byte) {
	b0 = byte(val % 256)
	val = val >> 8
	b1 = byte(val % 256)
	val = val >> 8
	b2 = byte(val % 256)
	val = val >> 8
	b3 = byte(val % 256)
	return b0, b1, b2, b3
}

func BytesToUint(b0, b1, b2, b3 byte) (val uint32) {
	val = uint32(b3)
	val = val << 8
	val += uint32(b2)
	val = val << 8
	val += uint32(b1)
	val = val << 8
	val += uint32(b0)
	return val
}
