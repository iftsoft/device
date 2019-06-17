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
	//	d.log.Info("Duplex WritePacket %+v", pack)
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
	pack.Print(d.log, "Write")
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
	if conn != nil {
		pack, err := conn.ReadPacket()
		if err != nil {
			return d.mngr.OnReadError(err)
		} else if pack != nil {
			pack.Print(d.log, "Read ")
			d.mngr.NewPacket(pack)
		}
	} else {
		return errors.New("duplex DialTCP conn is nil")
	}
	return nil
}

func (d *Duplex) readingLoop() {
	for {
		err := d.ReadPacket()
		if err != nil {
			d.log.Error("Duplex ReadPacket error: %s", err)
			return
		}
	}
}

func (d *Duplex) waitingLoop() {
	d.log.Info("Duplex loop is started")
	defer d.log.Info("Duplex loop is stopped")

	tick := time.NewTicker(1000 * time.Millisecond)
	defer tick.Stop()

	for {
		select {
		case <-d.done:
			defer d.log.Info("Duplex loop stopping")
			return
		case tm := <-tick.C:
			d.log.Info("Duplex loop timer tick")
			d.mngr.OnTimerTick(tm)
			//default:
			//	err := d.ReadPacket()
			//	if err != nil {
			//		d.log.Error("Duplex ReadPacket error: %s", err)
			//		//} else if err == io.EOF {
			//		//	break
			//		//} else {
			//		return
			//	}
			//			time.Sleep(100*time.Microsecond)
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
