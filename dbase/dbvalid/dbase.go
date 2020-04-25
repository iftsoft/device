package dbvalid

import (
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/dbase"
)

const (
	errDBaseNotOpen = "Database is not opened"
	errLinkerNotSet = "Database linker is not set"
	errStmtNotReady = "Statement is not prepared"
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
//	batch := ObjBatch{}
	qry := NewQueryBatch(db.linker, db.log)
	err := qry.doSelect(db.device, batch)
	if err != nil {
//		db.log.Error("InitNoteList select batch error: %s", err)
		return err
	}
	return err
}

func (db *DBaseValidator) closeLastBatch() error {
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
		return err
	}
	notes, err := qryNt.doSearch(db.device)
	if err != nil {
		return err
	}
	depos, err := qryDp.doSearch(batch.Id)
	if err != nil {
		return err
	}
	batch.State = checkBatch(notes, depos)
	if len(notes) > 0 {
		err = qryBl.doInsertNotes(batch.Id, notes)
	}
	err = qryBt.doUpdate(batch)
	if err != nil {
		return err
	}
	err = qryBt.makeNewBranch(db.device, batch)
	return err
}

func (db *DBaseValidator) InitNoteList(list common.ValidNoteList) error {
//	err := db.closeLastBatch()
//	if err != nil {
//		return err
//	}
	qryNt := NewQueryNote(db.linker, db.log)
	_, err := qryNt.doDelete(db.device)
	if err == nil {
		err = qryNt.doInsertNotes(db.device, list)
	}
	return err
}

func (db *DBaseValidator) ReadNoteList()(common.ValidNoteList, error) {
	qry := NewQueryNote(db.linker, db.log)
	notes, err := qry.doSearch(db.device)
	return notes, err
}

func (db *DBaseValidator) DepositNote(extraId int64, data *common.ValidatorAccept) error {
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
	_, err = qryDp.doInsertAccept(batch.Id, extraId, data)
	return err
}

func (db *DBaseValidator) SaveNoteList(list common.ValidNoteList) error {
//	objList := make(ObjNoteList, len(list))
//	for i, note := range list {
//		objList[i] = &ObjNote {
//			Device:   db.device,
//			Currency: uint16(note.Currency),
//			Nominal:  float32(note.Nominal),
//			Count:    uint16(note.Count),
//			Amount:   float32(note.Amount),
//		}
//	}
	qry := NewQueryNote(db.linker, db.log)
	err := qry.doUpdateEx(list)
	return err
}



