package dbvalid

import (
	"fmt"
	"github.com/iftsoft/device/common"
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
    opened VARCHAR(64),
    closed VARCHAR(64),
    detail TEXT
);`
	sqlBatchDelete = `DELETE FROM valid_batch WHERE id = ?;`
	sqlBatchSelect = `SELECT id, state, device, count, opened, closed, detail FROM valid_batch WHERE device = ? ORDER BY id desc LIMIT 1;`
	sqlBatchSearch = `SELECT id, state, device, count, opened, closed, detail FROM valid_batch WHERE device = ? ORDER BY id asc;`
	sqlBatchInsert = `INSERT INTO valid_batch (state, device, count, opened, closed, detail) VALUES (?, ?, ?, ?, ?, ?);`
	sqlBatchUpdate = `UPDATE valid_batch SET state = ?, count = ?, closed = ?, detail = ? WHERE id = ?;`
)

type ObjBatch struct {
	Id       int64
	Device   string
	State    common.BatchState
	Count    common.DevCounter
	Opened   string
	Closed   string
	Detail   string
}
type ObjBatchList []*ObjBatch

func (b *ObjBatch) String() string {
	if b == nil {
		return ""
	}
	str := fmt.Sprintf("Batch Id:%d, Device:%s, State:%s, Count:%d, Opened:%s, Closed:%s, Detail:%s",
		b.Id, b.Device, b.State.String(), b.Count, b.Opened, b.Closed, b.Detail)
	return str
}


type QueryBatch struct {
	dbase.DBaseQuery
}

func NewQueryBatch(linker dbase.DBaseLinker, log *core.LogAgent) *QueryBatch {
	qry := &QueryBatch{}
	qry.InitQuery(linker, log)
	return qry
}


func (qry *QueryBatch) doSelect(device string, batch *ObjBatch) error {
	param := make(dbase.ParamList, 1)
	param[0] = &device
	err := qry.RunSelectSql(sqlBatchSelect, param, batch)
	return err
}

func (qry *QueryBatch) doSearch(device string) (ObjBatchList, error) {
	items := make(ObjBatchList, 0)
	param := make(dbase.ParamList, 1)
	param[0] = &device
	err := qry.RunSearchSql(sqlBatchSearch, param, &items)
	return items, err
}

func (qry *QueryBatch) doDelete(id int) (int64, error) {
	param := make(dbase.ParamList, 1)
	param[0] = &id
	err := qry.RunCommandSql(sqlBatchDelete, param)
	return qry.RowsAffected(), err
}

func (qry *QueryBatch)doInsert(batch *ObjBatch) error {
	param := make(dbase.ParamList, 6)
	param[0] = &batch.State
	param[1] = &batch.Device
	param[2] = &batch.Count
	param[3] = &batch.Opened
	param[4] = &batch.Closed
	param[5] = &batch.Detail
	err := qry.RunCommandSql(sqlBatchInsert, param)
	if err == nil {
		batch.Id = qry.LastInsertId()
	}
	return err
}

func (qry *QueryBatch)doUpdate(batch *ObjBatch) error {
	param := make(dbase.ParamList, 5)
	param[0] = &batch.State
	param[1] = &batch.Count
	param[2] = &batch.Closed
	param[3] = &batch.Detail
	param[4] = &batch.Id
	err := qry.RunCommandSql(sqlBatchUpdate, param)
	return err
}

func (qry *QueryBatch)makeNewBranch(device string, batch *ObjBatch) error {
	batch.Id     = 0
	batch.State  = common.StateEmpty
	batch.Device = device
	batch.Count  = 0
	batch.Opened = time.Now().Format(timeFormat)
	batch.Closed = ""
	batch.Detail = ""
	return qry.doInsert(batch)
}


