package duplex

import (
	"errors"
	"github.com/iftsoft/device/core"
	"net"
	"sync"
	"time"
)

const (
	aMagicByte1 byte = 0x39
	aMagicByte2 byte = 0x7B
	aMagicByte3 byte = 0xA2
	aMagicByte4 byte = 0x5F
	aHeaderSize int  = 8
)

type Connection struct {
	conn net.TCPConn
	log  *core.LogAgent
	lock sync.Mutex
}

func (c *Connection) Close() {
	_ = c.conn.Close()
}

func (c *Connection) WritePacket(pack *Packet) error {
	if pack == nil {
		return errors.New("packet pointer is nil")
	}
	dump := pack.Encode()
	c.lock.Lock()
	defer c.lock.Unlock()
	err := c.WriteBinary(dump)
	return err
}

func (c *Connection) WriteBinary(dump []byte) error {
	size := len(dump)
	if size == 0 {
		return nil
	}

	head := make([]byte, aHeaderSize)
	head[0] = aMagicByte1
	head[1] = aMagicByte2
	head[2] = aMagicByte3
	head[3] = aMagicByte4
	head[4], head[5], head[6], head[7] = UintToBytes(uint32(size))

	_, err := c.conn.Write(head)
	if err != nil {
		return err
	}
	_, err = c.conn.Write(dump)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) ReadPacket() (pack *Packet, err error) {
	var dump []byte
	dump, err = c.ReadBinary()
	if err != nil {
		return nil, err
	}
	pack = &Packet{}
	err = pack.Decode(dump)
	return pack, err
}

func (c *Connection) ReadBinary() (dump []byte, err error) {
	now := time.Now()
	now.Add(time.Millisecond)
	err = c.conn.SetReadDeadline(now)
	if err != nil {
		return nil, err
	}
	head := make([]byte, aHeaderSize)
	_, err = c.conn.Read(head)
	if err != nil {
		return nil, err
	}
	size := BytesToUint(head[4], head[5], head[6], head[7])
	dump = make([]byte, size)
	_, err = c.conn.Read(dump)
	if err != nil {
		return nil, err
	}
	return dump, nil
}
