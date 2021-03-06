package duplex

import (
	"errors"
	"github.com/iftsoft/device/core"
	"io"
	"net"
	"sync"
)

const (
	aMagicByte1 byte = 0x39
	aMagicByte2 byte = 0x7B
	aMagicByte3 byte = 0xA2
	aMagicByte4 byte = 0x5F
	aHeaderSize int  = 8
)

type Connection struct {
	conn *net.TCPConn
	log  *core.LogAgent
	lock sync.Mutex
	exit bool
}

func (c *Connection) Close() {
	c.exit = true
	_ = c.conn.Close()
}

func (c *Connection) WritePacket(pack *Packet) error {
	if pack == nil {
		return errors.New("packet pointer is nil")
	}
	dump := pack.Encode()
	//	c.log.Trace("Write packet dump: %+v", dump)
	c.lock.Lock()
	defer c.lock.Unlock()
	err := c.WriteBinary(dump)
	return err
}

func (c *Connection) WriteBinary(dump []byte) error {
	size := len(dump)
	if size == 0 || c.exit {
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
	//	c.log.Trace("Read packet conn: %+v", c)
	dump, err = c.ReadBinary()
	if err != nil {
		if err != io.EOF && c.exit == false {
			c.log.Error("Read packet error: %s", err)
		}
		return nil, err
	}
	pack = &Packet{}
	if dump != nil {
		//		c.log.Dump("Read packet dump: %+v", dump)
		err = pack.Decode(dump)
	}
	return pack, err
}

func (c *Connection) ReadBinary() (dump []byte, err error) {
	head := make([]byte, aHeaderSize)
	n := 0
	n, err = c.conn.Read(head)
	if err != nil {
		if err != io.EOF && c.exit == false {
			c.log.Error("Connection Read header error: %s", err)
		}
		return nil, err
	}
	if n != aHeaderSize {
		c.log.Warn("Connection Read header size: %d of %d bytes", n, aHeaderSize)
		return nil, errors.New("wrong header size")
	}
	size := BytesToUint(head[4], head[5], head[6], head[7])
	dump = make([]byte, size)
	n, err = c.conn.Read(dump)
	if err != nil {
		c.log.Error("Connection Read binary error: %s", err)
		return nil, err
	}
	if n != int(size) {
		c.log.Warn("Connection Read binary size: %d of %d bytes", n, aHeaderSize)
		return nil, errors.New("wrong binary size")
	}
	return dump, nil
}
