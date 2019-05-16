package duplex

import (
	"errors"
	"github.com/iftsoft/device/core"
	"sync"
)

const DuplexPort int32 = 9380

type DuplexServerConfig struct {
	Port int32 `yaml:"port"`
}

type DuplexClientConfig struct {
	Port int32 `yaml:"port"`
}

type DuplexManager interface {
	NewPacket(pack *Packet) bool
	OnWriteError(d *Duplex, err error) error
	OnReadError(d *Duplex, err error) error
	OnTimerTick(d *Duplex)
}

type Duplex struct {
	Link LinkHolder
	mngr DuplexManager
	done chan struct{}
	log  *core.LogAgent
}

func (d *Duplex) WritePacket(pack *Packet) error {
	if d.mngr == nil {
		return errors.New("duplex manager is not set")
	}
	if pack == nil {
		return errors.New("packet pointer is nil")
	}
	conn := d.Link.GetConnect()
	if conn == nil {
		return errors.New("connection is closed")
	}
	err := conn.WritePacket(pack)
	if err != nil {
		err = d.mngr.OnWriteError(d, err)
	}
	return err
}

func (d *Duplex) ReadPacket() error {
	if d.mngr == nil {
		return errors.New("duplex manager is not set")
	}
	conn := d.Link.GetConnect()
	if conn == nil {
		return errors.New("connection is closed")
	}
	pack, err := conn.ReadPacket()
	if err != nil {
		err = d.mngr.OnReadError(d, err)
	} else {
		d.mngr.NewPacket(pack)
	}
	return err
}

func (d *Duplex) readingLoop(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-d.done:
			return
		default:
			err := d.ReadPacket()
			if err != nil {
				//} else if err == io.EOF {
				//	break
				//} else {
				return
			}
		}
	}
}
