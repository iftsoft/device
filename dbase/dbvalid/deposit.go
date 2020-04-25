package dbvalid

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/dbase"
	"time"
)

const (
	sqlDepositCreate = `CREATE TABLE IF NOT EXISTS valid_deposit (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	batch_id INTEGER NOT NULL,
	extra_id INTEGER NOT NULL DEFAULT 0,
    currency INTEGER NOT NULL DEFAULT 0,
    nominal REAL NOT NULL DEFAULT 0,
    count INTEGER NOT NULL DEFAULT 0,
    amount REAL NOT NULL DEFAULT 0,
    created VARCHAR(64),
    FOREIGN KEY (batch_id) REFERENCES valid_batch (id) ON UPDATE RESTRICT ON DELETE RESTRICT
);`
	sqlDepositDelete = `DELETE FROM valid_deposit WHERE batch_id = ?;`
	sqlDepositSelect = `SELECT id, batch_id, extra_id, currency, nominal, count, amount, created FROM valid_deposit WHERE id = ?;`
	sqlDepositSearch = `SELECT id, batch_id, extra_id, currency, nominal, count, amount, created FROM valid_deposit WHERE batch_id = ?;`
	sqlDepositInsert = `INSERT INTO valid_deposit (batch_id, extra_id, currency, nominal, count, amount, created) VALUES (?, ?, ?, ?, ?, ?, ?);`
)

type ObjDeposit struct {
	Id       int64
	BatchId  int64
	ExtraId  int64
	Currency common.DevCurrency
	Nominal  common.DevAmount
	Count    common.DevCounter
	Amount   common.DevAmount
	Created  time.Time
}
type ObjDepositList []*ObjDeposit


type QueryDeposit struct {
	dbase.DBaseQuery
}

func NewQueryDeposit(linker dbase.DBaseLinker, log *core.LogAgent) *QueryDeposit {
	qry := &QueryDeposit{}
	qry.InitQuery(linker, log)
	return qry
}


func (qry *QueryDeposit)doSelect(id int64, depo *ObjDeposit) error {
	param := make(dbase.ParamList, 1)
	param[0] = &id
	err := qry.RunSelectSql(sqlDepositSelect, param, depo)
	return err
}

func (qry *QueryDeposit)doSearch(batchId int64) (ObjDepositList, error) {
	items := make(ObjDepositList, 0)
	param := make(dbase.ParamList, 1)
	param[0] = &batchId
	err := qry.RunSearchSql(sqlDepositSearch, param, &items)
	return items, err
}

func (qry *QueryDeposit)doDelete(batchId int64) (int64, error) {
	param := make(dbase.ParamList, 1)
	param[0] = &batchId
	err := qry.RunCommandSql(sqlDepositDelete, param)
	return qry.RowsAffected(), err
}

func (qry *QueryDeposit)doInsert(depo *ObjDeposit) error {
	param := make(dbase.ParamList, 7)
	param[0] = &depo.BatchId
	param[1] = &depo.ExtraId
	param[2] = &depo.Currency
	param[3] = &depo.Nominal
	param[4] = &depo.Count
	param[5] = &depo.Amount
	param[6] = &depo.Created
	err := qry.RunCommandSql(sqlDepositInsert, param)
	if err == nil {
		depo.Id = qry.LastInsertId()
	}
	return err
}

func (qry *QueryDeposit)doInsertEx(depos ObjDepositList) error {
	parList := make([]dbase.ParamList, len(depos))
	for i, depo := range depos {
		param := make(dbase.ParamList, 7)
		param[0] = &depo.BatchId
		param[1] = &depo.ExtraId
		param[2] = &depo.Currency
		param[3] = &depo.Nominal
		param[4] = &depo.Count
		param[5] = &depo.Amount
		param[6] = &depo.Created
		parList[i] = param
	}
	err := qry.RunPreparedSql(sqlDepositInsert, parList)
	return err
}

func (qry *QueryDeposit)doInsertAccept(batchId, extraId int64, data *common.ValidatorAccept) (*ObjDeposit, error) {
	depo := &ObjDeposit{
		BatchId:   batchId,
		ExtraId:   extraId,
		Currency:  data.Currency,
		Nominal:   data.Nominal,
		Count:     data.Count,
		Amount:    data.Amount,
		Created:   time.Now(),
	}
	err := qry.doInsert(depo)
	return depo, err
}

