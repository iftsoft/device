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
	ScopeDevice
	ScopeReader
	ScopeValidator
	ScopeMax
)

var listScopeName = []string{
	"System",
	"Device",
	"Reader",
	"Validator",
	"",
}

func GetScopeName(scope PacketScope) string {
	if scope < ScopeMax {
		return listScopeName[scope]
	}
	return "Unknown"
}

type Packet struct {
	Version PacketVersion
	Type    PacketType
	Scope   PacketScope
	Counter uint32
	Options uint32
	DevName string
	Command string
	Content []byte
}

func NewPacket(scope PacketScope, name string, cmd string, data []byte) *Packet {
	p := Packet{
		Version: packetVersion,
		Type:    PackerRequest,
		Scope:   scope,
		DevName: name,
		Command: cmd,
		Counter: 0,
		Options: 0,
		Content: data,
	}
	return &p
}

func (p *Packet) Print(log *core.LogAgent, text string) {
	if p != nil && log != nil {
		log.Dump("%s packet Scope:%s, Device:%s, Command:%s, Data len:%d, Content:%s",
			text, GetScopeName(p.Scope), p.DevName, p.Command, len(p.Content), string(p.Content))
	}
}

func (p *Packet) Encode() []byte {
	head := make([]byte, packetHeaderSize)
	name := []byte(p.DevName)
	task := []byte(p.Command)
	len1 := len(name)
	len2 := len(task)
	size := len(p.Content)

	//	head[0] = byte(p.Version)
	head[0] = byte(p.Type)
	head[1] = byte(p.Scope)
	head[2] = byte(len1)
	head[3] = byte(len2)
	head[4], head[5], head[6], head[7] = UintToBytes(p.Options)
	head[8], head[9], head[10], head[11] = UintToBytes(p.Counter)
	head[12], head[13], head[14], head[15] = UintToBytes(uint32(size))
	dump := append(head, name...)
	dump = append(dump, task...)
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
	len1 := int(head[2])
	len2 := int(head[3])
	if full < packetHeaderSize+len1+len2 {
		return err
	}
	name := dump[packetHeaderSize : packetHeaderSize+len1]
	task := dump[packetHeaderSize+len1 : packetHeaderSize+len1+len2]
	//	p.Version = PacketVersion(head[0])
	p.Type = PacketType(head[0])
	p.Scope = PacketScope(head[1])
	p.Options = BytesToUint(head[4], head[5], head[6], head[7])
	p.Counter = BytesToUint(head[8], head[9], head[10], head[11])
	p.DevName = string(name)
	p.Command = string(task)
	size := int(BytesToUint(head[12], head[13], head[14], head[15]))
	if full < packetHeaderSize+len1+len2+size {
		return err
	}
	if size > 0 {
		p.Content = dump[packetHeaderSize+len1+len2:]
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
