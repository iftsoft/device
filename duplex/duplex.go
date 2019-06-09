package duplex

import (
	"errors"
	"github.com/iftsoft/device/core"
	"time"
)

const DuplexPort int32 = 9380

type DuplexManager interface {
	NewPacket(pack *Packet) bool
	//	OnNoConnect(err error) error
	OnWriteError(err error) error
	OnReadError(err error) error
	OnTimerTick(tm time.Time)
}

type Duplex struct {
	link LinkHolder
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
	conn := d.link.GetConnect()
	if conn == nil {
		return errors.New("connection is closed")
	}
	d.log.Trace("Write packet: %+v", pack)
	err := conn.WritePacket(pack)
	if err != nil {
		err = d.mngr.OnWriteError(err)
	}
	return err
}

func (d *Duplex) ReadPacket() error {
	if d.mngr == nil {
		return errors.New("duplex manager is not set")
	}
	conn := d.link.GetConnect()
	if conn == nil {
		return errors.New("connection is closed")
	}
	pack, err := conn.ReadPacket()
	d.log.Trace("Read packet: %+v", pack)
	if err != nil {
		err = d.mngr.OnReadError(err)
	} else {
		d.mngr.NewPacket(pack)
	}
	return err
}

func (d *Duplex) readingLoop() {
	tick := time.NewTicker(10 * time.Millisecond)
	defer tick.Stop()
	for {
		select {
		case <-d.done:
			return
		case tm := <-tick.C:
			d.mngr.OnTimerTick(tm)
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

/*
// keepDoingSomething will keep trying to doSomething() until either
// we get a result from doSomething() or the timeout expires
func keepDoingSomething() (bool, error) {
	timeout := time.After(5 * time.Second)
	tick := time.Tick(500 * time.Millisecond)
	// Keep trying until we're timed out or got a result or got an error
	for {
		select {
		// Got a timeout! fail with a timeout error
		case <-timeout:
			return false, errors.New("timed out")
		// Got a tick, we should check on doSomething()
		case <-tick:
			ok, err := doSomething()
			// Error from doSomething(), we should bail
			if err != nil {
				return false, err
			// doSomething() worked! let's finish up
			} else if ok {
				return true, nil
			}
			// doSomething() didn't work yet, but it didn't fail, so let's try again
			// this will exit up to the for loop
		}
	}
}
*/
