package dbvalid

import (
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/dbase"
	"time"
)

const (
	sqlBatchCreate = `CREATE TABLE IF NOT EXISTS valid_batch (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
    state INTEGER NOT NULL DEFAULT 0,
	device VARCHAR(64) NOT NULL,
    count INTEGER NOT NULL DEFAULT 0,
    opened TEXT,
    closed TEXT
);`
	sqlBatchDelete = `DELETE FROM valid_batch WHERE device = ?;`
	sqlBatchSelect = `SELECT id, state, device, count, opened, closed FROM valid_batch WHERE device = ? order by id desc;`
	sqlBatchSearch = `SELECT id, state, device, count, opened, closed FROM valid_batch WHERE device = ? order by id asc;`
	sqlBatchInsert = `INSERT INTO valid_batch (state, device, count, opened, closed) VALUES (?, ?, ?, ?, ?);`
	sqlBatchUpdate = `UPDATE valid_batch SET state = ?, count = ?, closed = ? WHERE id = ?;`
)

const (
	StateUndefined uint16 = iota
	StateActive
	StateClosed
)
type ObjBatch struct {
	Id       int
	State    uint16
	Device   string
	Count    uint16
	Opened   time.Time
	Closed   time.Time
}


type QueryBatch struct {
	dbase.DBaseQuery
}

func NewQueryBatch(linker dbase.DBaseLinker, log *core.LogAgent) *QueryBatch {
	qry := &QueryBatch{}
	qry.InitQuery(linker, log)
	return qry
}


