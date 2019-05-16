package duplex

import (
	"github.com/iftsoft/device/core"
	"net"
	"sync"
)

type LinkHolder struct {
	conn *Connection
	lock sync.Mutex
}

func (h *LinkHolder) GetConnect() *Connection {
	h.lock.Lock()
	defer h.lock.Unlock()
	return h.conn
}

func (h *LinkHolder) SetConnect(conn net.TCPConn, log *core.LogAgent) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.conn = &Connection{conn: conn, log: log}
}

func (h *LinkHolder) CloseConnect() {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.conn != nil {
		h.conn.Close()
	}
	h.conn = nil
}
