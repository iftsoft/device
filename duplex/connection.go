package duplex

import (
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

//type DumpReader interface {
//	PacketDump(dump []byte) bool
//}

type Connection struct {
	conn net.TCPConn
	//	reader DumpReader
	//	done   chan struct{}
	wrLock sync.Mutex
	rdLock sync.Mutex
}

//func NewConnection(conn net.TCPConn, reader DumpReader, done chan struct{}) *Connection {
//	return &Connection{conn: conn, reader: reader, done: done}
//}

func (c *Connection) WriteBinary(dump []byte) error {
	size := len(dump)
	if size == 0 {
		return nil
	}
	c.wrLock.Lock()
	defer c.wrLock.Unlock()

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

func (c *Connection) ReadBinary() (dump []byte, err error) {
	c.rdLock.Lock()
	defer c.rdLock.Unlock()

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

//func (c *Connection) readingLoop(wg *sync.WaitGroup) {
//	defer wg.Done()
//	for {
//		select {
//		case <-c.done:
//			return
//		default:
//			dump, err := c.ReadBinary()
//			if err != nil {
//				c.reader.PacketDump(dump)
//			} else if err == io.EOF {
//				break
//			} else {
//				return
//			}
//		}
//	}
//}
