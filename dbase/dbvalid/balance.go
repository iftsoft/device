package dbvalid

import (
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/dbase"
	"time"
)

const (
	sqlBalanceCreate = `CREATE TABLE IF NOT EXISTS valid_balance (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	batch_id INTEGER NOT NULL,
    currency INTEGER NOT NULL DEFAULT 0,
    nominal INTEGER NOT NULL DEFAULT 0,
    count INTEGER NOT NULL DEFAULT 0,
    amount INTEGER NOT NULL DEFAULT 0,
    created VARCHAR(40),
    FOREIGN KEY (batch_id) REFERENCES valid_batch (id) ON UPDATE RESTRICT ON DELETE RESTRICT
);`
	sqlBalanceDelete = `DELETE FROM valid_balance WHERE batch_id = ?;`
	sqlBalanceSelect = `SELECT id, batch_id, currency, nominal, count, amount, created FROM valid_balance WHERE id = ?;`
	sqlBalanceSearch = `SELECT id, batch_id, currency, nominal, count, amount, created FROM valid_balance WHERE batch_id = ?;`
	sqlBalanceInsert = `INSERT INTO valid_balance (batch_id, extra_id, currency, nominal, count, amount, created) VALUES (?, ?, ?, ?, ?, ?, ?);`
)

type ObjBalance struct {
	Id       int64
	BatchId  int64
	Currency uint16
	Nominal  float32
	Count    uint16
	Amount   float32
	Created  time.Time
}
type ObjBalanceList []*ObjBalance

type QueryBalance struct {
	dbase.DBaseQuery
}

func NewQueryBalance(linker dbase.DBaseLinker, log *core.LogAgent) *QueryBalance {
	qry := &QueryBalance{}
	qry.InitQuery(linker, log)
	return qry
}


func (qry *QueryBalance)Select(bal *ObjBalance) (error) {
	param := make(dbase.ParamList, 1)
	param[0] = &bal.Id
	err := qry.RunSelectSql(sqlBalanceSelect, param, bal)
	return err
}

func (qry *QueryBalance)Search(batch_id int) (ObjBalanceList, error) {
	items := make(ObjBalanceList, 0)
	param := make(dbase.ParamList, 1)
	param[0] = &batch_id
	err := qry.RunSearchSql(sqlBalanceSearch, param, &items)
	return items, err
}

func (qry *QueryBalance)Delete(batch_id int) (int64, error) {
	param := make(dbase.ParamList, 1)
	param[0] = &batch_id
	err := qry.RunCommandSql(sqlBalanceDelete, param)
	return qry.RowsAffected(), err
}

func (qry *QueryBalance)Insert(bal *ObjBalance) error {
	param := make(dbase.ParamList, 6)
	param[0] = &bal.BatchId
	param[1] = &bal.Currency
	param[2] = &bal.Nominal
	param[3] = &bal.Count
	param[4] = &bal.Amount
	param[5] = &bal.Created
	err := qry.RunCommandSql(sqlBalanceInsert, param)
	if err == nil {
		bal.Id = qry.LastInsertId()
	}
	return err
}

func (qry *QueryBalance)InsertEx(bals ObjBalanceList) error {
	parList := make([]dbase.ParamList, len(bals))
	for i, bal := range bals {
		param := make(dbase.ParamList, 6)
		param[0] = &bal.BatchId
		param[1] = &bal.Currency
		param[2] = &bal.Nominal
		param[3] = &bal.Count
		param[4] = &bal.Amount
		param[5] = &bal.Created
		parList[i] = param
	}
	err := qry.RunPreparedSql(sqlBalanceInsert, parList)
	return err
}


