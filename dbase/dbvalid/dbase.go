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
	return err
}

func (db *DBaseValidator) InitNoteList(list common.ValidNoteList) error {
	objList := make(ObjNoteList, len(list))
	for i, note := range list {
		objList[i] = &ObjNote {
			Device:   db.device,
			Currency: uint16(note.Currency),
			Nominal:  float32(note.Nominal),
			Count:    uint16(note.Count),
			Amount:   float32(note.Amount),
		}
	}
	qry := NewQueryNote(db.linker, db.log)
	_, err := qry.Delete(db.device)
	if err == nil {
		err = qry.InsertEx(objList)
	}
	return err
}



