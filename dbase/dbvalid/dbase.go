package dbvalid

import (
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/dbase"
	"time"
)

const (
	errParamIsNil   = "parameter is nil pointer"
	errDBaseNotOpen = "database is not opened"
	errLinkerNotSet = "database linker is not set"
//	errStmtNotReady = "statement is not prepared"
)

type DBaseValidator struct{
	linker dbase.DBaseLinker
	device string
	log    *core.LogAgent
}

func NewDBaseValidator(linker dbase.DBaseLinker, device string) *DBaseValidator {
	db := &DBaseValidator{
		linker: linker,
		device: device,
		log: core.GetLogAgent(core.LogLevelTrace, "DbValid"),
	}
	return db
}

func (db *DBaseValidator) Open() error {
	if db.linker == nil {
		return errors.New(errLinkerNotSet)
	}
	err := db.linker.Open()
	return err
}

func (db *DBaseValidator) Close() (err error) {
	if db.linker != nil {
		err = db.linker.Close()
	}
	return err
}

func (db *DBaseValidator) CreateAllTables() error {
	qry := NewQueryValidator(db.linker, db.log)
	err := qry.CreateTableNote()
	err = qry.CreateTableBatch()
	err = qry.CreateTableDeposit()
	err = qry.CreateTableBalance()
	return err
}



func (db *DBaseValidator) GetLastBatch(batch *ObjBatch) error {
	qry := NewQueryBatch(db.linker, db.log)
	err := qry.doSelect(db.device, batch)
	if err != nil {
//		db.log.Error("InitNoteList select batch error: %s", err)
		return err
	}
	return err
}

func (db *DBaseValidator) CloseBatch(data *common.ValidatorBatch) error {
	if data == nil {
		return errors.New(errParamIsNil)
	}
	qryBt := NewQueryBatch(db.linker, db.log)
	qryNt := NewQueryNote(db.linker, db.log)
	qryDp := NewQueryDeposit(db.linker, db.log)
	qryBl := NewQueryBalance(db.linker, db.log)
	batch := &ObjBatch{}

	err := qryBt.doSelect(db.device, batch)
	if err != nil {
		return err
	}
	if batch.Id == 0 {
		err = qryBt.makeNewBranch(db.device, batch)
	}
	if err != nil {
		return err
	}
	data.BatchId = batch.Id
	data.State   = batch.State
	data.Notes, err = qryNt.doSearch(db.device)
	if err != nil {
		return err
	}
	if batch.State == common.StateEmpty {
		return err
	}
	depos, er := qryDp.doSearch(batch.Id)
	if er != nil {
		return er
	}
	batch.State, batch.Detail = checkBatch(data.Notes, depos)
	batch.Closed = time.Now().Format(timeFormat)
	err = qryBt.doUpdate(batch)
	if err != nil {
		return err
	}
	err = qryBl.doInsertNotes(batch.Id, data.Notes)
	if err != nil {
		return err
	}
	err = qryBt.makeNewBranch(db.device, batch)
	return err
}

func (db *DBaseValidator) InitNoteList(list common.ValidNoteList) error {
	qryBt := NewQueryBatch(db.linker, db.log)
	qryNt := NewQueryNote(db.linker, db.log)
	batch := &ObjBatch{}

	err := qryBt.doSelect(db.device, batch)
	if err != nil {
		return err
	}
	if batch.State != common.StateEmpty {
		err = errors.New("batch is not empty")
	}
	_, err = qryNt.doDelete(db.device)
	if err == nil {
		err = qryNt.doInsertNotes(db.device, list)
	}
	return err
}

func (db *DBaseValidator) ReadNoteList(data *common.ValidatorBatch) error {
	if data == nil {
		return errors.New(errParamIsNil)
	}
	qryBt := NewQueryBatch(db.linker, db.log)
	qryNt := NewQueryNote(db.linker, db.log)
	batch := &ObjBatch{}

	err := qryBt.doSelect(db.device, batch)
	if err != nil {
		return err
	}
	data.BatchId = batch.Id
	data.State   = batch.State
	data.Notes, err = qryNt.doSearch(db.device)
	return err
}

func (db *DBaseValidator) DepositNote(extraId int64, data *common.ValidatorAccept) error {
	if data == nil {
		return errors.New(errParamIsNil)
	}
	qryBt := NewQueryBatch(db.linker, db.log)
	qryNt := NewQueryNote(db.linker, db.log)
	qryDp := NewQueryDeposit(db.linker, db.log)
	batch := &ObjBatch{}

	err := qryBt.doSelect(db.device, batch)
	if err != nil {
		return err
	}
//	db.log.Info("Deposit batch - %s", batch.String())
	if batch.Id == 0 {
		err = errors.New("no active batch")
		return err
	}
	err = qryNt.doUpdateAccept(db.device, data)
	if err != nil {
		return err
	}
	depo, er := qryDp.doInsertAccept(batch.Id, extraId, data)
	if er != nil {
		return er
	}
	db.log.Info("Deposit note - %s", depo.String())
	batch.State = common.StateActive
	batch.Count += data.Count
	err = qryBt.doUpdate(batch)
	return err
}



