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

type DBaseValidator struct {
	linker dbase.DBaseLinker
	device string
	log    *core.LogAgent
}

func NewDBaseValidator(linker dbase.DBaseLinker, device string) *DBaseValidator {
	db := &DBaseValidator{
		linker: linker,
		device: device,
		log:    core.GetLogAgent(core.LogLevelTrace, "DbValid"),
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
	db.log.Debug("Validator database - CreateAllTables")
	err := db.linker.Begin()
	if err == nil {
		qry := NewQueryValidator(db.linker, db.log)
		err = qry.CreateTableNote()
		if err == nil {
			err = qry.CreateTableBatch()
		}
		if err == nil {
			err = qry.CreateTableDeposit()
		}
		if err == nil {
			err = qry.CreateTableBalance()
		}
		if err == nil {
			err = db.linker.Commit()
		} else {
			db.log.Error("Can't create tables: %s", err)
			_ = db.linker.Rollback()
		}
	}
	return err
}

func (db *DBaseValidator) GetLastBatch(batch *ObjBatch) error {
	db.log.Debug("Validator database - GetLastBatch")
	qry := NewQueryBatch(db.linker, db.log)
	err := qry.doSelect(db.device, batch)
	if err != nil {
		db.log.Error("Can't get last batch: %s", err)
		return err
	}
	db.log.Dump("Current batch - %s", batch.String())
	return err
}

func (db *DBaseValidator) CloseBatch(data *common.ValidatorBatch) error {
	if data == nil {
		return errors.New(errParamIsNil)
	}
	db.log.Debug("Validator database - CloseBatch")
	err := db.linker.Begin()
	if err == nil {
		err = db.tryCloseBatch(data)
		if err == nil {
			err = db.linker.Commit()
		} else {
			db.log.Error("Can't close batch: %s", err)
			_ = db.linker.Rollback()
		}
	}
	return err
}

func (db *DBaseValidator) tryCloseBatch(data *common.ValidatorBatch) error {
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
		db.log.Dump("Empty batch - %s", data.String())
		return err
	}
	depos, er := qryDp.doSearch(batch.Id)
	if er != nil {
		return er
	}
	batch.State, batch.Detail = checkBatch(data.Notes, depos)
	batch.Closed = time.Now().Format(timeFormat)
	data.State   = batch.State
	data.Detail  = batch.Detail
	db.log.Dump("Closed batch - %s", data.String())
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
	db.log.Debug("Validator database - InitNoteList")
	err := db.linker.Begin()
	if err == nil {
		err = db.tryInitNoteList(list)
		if err == nil {
			err = db.linker.Commit()
		} else {
			db.log.Error("Can't init note list: %s", err)
			_ = db.linker.Rollback()
		}
	}
	return err
}

func (db *DBaseValidator) tryInitNoteList(list common.ValidNoteList) error {
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
	db.log.Debug("Validator database - ReadNoteList")
	err := db.tryReadNoteList(data)
	if err == nil {
		db.log.Dump("Current batch - %s", data.String())
	} else {
		db.log.Error("Can't read note list: %s", err)
	}
	return err
}

func (db *DBaseValidator) tryReadNoteList(data *common.ValidatorBatch) error {
	qryBt := NewQueryBatch(db.linker, db.log)
	qryNt := NewQueryNote(db.linker, db.log)
	batch := &ObjBatch{}

	err := qryBt.doSelect(db.device, batch)
	if err != nil {
		return err
	}
	data.BatchId = batch.Id
	data.State = batch.State
	data.Notes, err = qryNt.doSearch(db.device)
	return err
}


func (db *DBaseValidator) DepositNote(extraId int64, data *common.ValidatorAccept) error {
	if data == nil {
		return errors.New(errParamIsNil)
	}
	db.log.Debug("Validator database - DepositNote %7.2f %s", data.Amount, data.Currency.IsoCode())
	err := db.linker.Begin()
	if err == nil {
		err = db.tryDepositNote(extraId, data)
		if err == nil {
			err = db.linker.Commit()
		} else {
			db.log.Error("Can't deposit note: %s", err)
			_ = db.linker.Rollback()
		}
	}
	return err
}

func (db *DBaseValidator) tryDepositNote(extraId int64, data *common.ValidatorAccept) error {
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
	db.log.Dump("Deposit note - %s", depo.String())
	batch.State = common.StateActive
	batch.Count += data.Count
	err = qryBt.doUpdate(batch)
	return err
}
