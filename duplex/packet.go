package duplex

import "errors"

const packetHeaderSize int = 16

type PacketType byte
type PacketScope byte

const (
	PacketNotify PacketType = iota
	PackerRequest
	PacketResponse
	PacketCallback
	PacketBackword
)

const (
	ScopeSystem PacketScope = iota
	ScopeConfig
	ScopeDevice
)

type Packet struct {
	Type    PacketType
	Scope   PacketScope
	Counter uint32
	Options uint32
	Command string
	Content []byte
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
	head[0] = byte(p.Type)
	head[1] = byte(p.Scope)
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
	p.Type = PacketType(head[0])
	p.Scope = PacketScope(head[1])
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
