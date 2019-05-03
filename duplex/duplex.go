package duplex

import (
	"errors"
	"net"
)

const DuplexPort int32 = 9380

type DuplexServerConfig struct {
	Port int32 `yaml:"port"`
}

type DuplexClientConfig struct {
	Port int32 `yaml:"port"`
}

type Duplex struct {
	Conn net.TCPConn
}

func (d *Duplex) WritePacket(pack *Packet) error {
	if pack == nil {
		return errors.New("packet pointer is nil")
	}
	return nil
}

func (d *Duplex) ReadPacket() (pack *Packet, err error) {
	pack = &Packet{}
	return pack, nil
}
