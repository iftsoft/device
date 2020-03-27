package dbvalid

import (
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/dbase"
)

const (
	sqlNoteCreate = `CREATE TABLE IF NOT EXISTS valid_note (
	device VARCHAR(64) NOT NULL,
    currency INTEGER NOT NULL DEFAULT 0,
    nominal INTEGER NOT NULL DEFAULT 0,
    count INTEGER NOT NULL DEFAULT 0,
    amount INTEGER NOT NULL DEFAULT 0,
    UNIQUE (device, currency, nominal)
);`
	sqlNoteDelete = `DELETE FROM valid_note WHERE device = ?;`
	sqlNoteSelect = `SELECT device, currency, nominal, count, amount FROM valid_note WHERE device = ?, currency = ?, nominal = ?;`
	sqlNoteSearch = `SELECT device, currency, nominal, count, amount FROM valid_note WHERE device = ?;`
	sqlNoteInsert = `INSERT INTO valid_note (device, currency, nominal, count, amount) VALUES (?, ?, ?, ?, ?);`
	sqlNoteUpdate = `UPDATE valid_note SET count = count + ?, amount = amount + ? WHERE device = ?, currency = ?, nominal = ?;`
)

type ObjNote struct {
	Device   string
	Currency uint16
	Nominal  float32
	Count    uint16
	Amount   float32
}
type ObjNoteList []*ObjNote

type QueryNote struct {
	dbase.DBaseQuery
}

func NewQueryNote(linker dbase.DBaseLinker, log *core.LogAgent) *QueryNote {
	qry := &QueryNote{}
	qry.InitQuery(linker, log)
	return qry
}

func (qry *QueryNote)doSelect(note *ObjNote) error {
	param := make(dbase.ParamList, 3)
	param[0] = &note.Device
	param[1] = &note.Currency
	param[2] = &note.Nominal
	err := qry.RunSelectSql(sqlNoteSelect, param, note)
	return err
}

func (qry *QueryNote)doSearch(device string) (ObjNoteList, error) {
	items := make(ObjNoteList, 0)
	param := make(dbase.ParamList, 1)
	param[0] = &device
	err := qry.RunSearchSql(sqlNoteSearch, param, &items)
	return items, err
}

func (qry *QueryNote)doDelete(device string) (int64, error) {
	param := make(dbase.ParamList, 1)
	param[0] = &device
	err := qry.RunCommandSql(sqlNoteDelete, param)
	return qry.RowsAffected(), err
}

func (qry *QueryNote)doInsert(note *ObjNote) error {
	param := make(dbase.ParamList, 5)
	param[0] = &note.Device
	param[1] = &note.Currency
	param[2] = &note.Nominal
	param[3] = &note.Count
	param[4] = &note.Amount
	err := qry.RunCommandSql(sqlNoteInsert, param)
	return err
}

func (qry *QueryNote)doUpdate(note *ObjNote) error {
	param := make(dbase.ParamList, 5)
	param[0] = &note.Count
	param[1] = &note.Amount
	param[2] = &note.Device
	param[3] = &note.Currency
	param[4] = &note.Nominal
	err := qry.RunCommandSql(sqlNoteUpdate, param)
	return err
}

func (qry *QueryNote)doInsertEx(notes ObjNoteList) error {
	parList := make([]dbase.ParamList, len(notes))
	for i, note := range notes {
		param := make(dbase.ParamList, 5)
		param[0] = &note.Device
		param[1] = &note.Currency
		param[2] = &note.Nominal
		param[3] = &note.Count
		param[4] = &note.Amount
		parList[i] = param
	}
	err := qry.RunPreparedSql(sqlNoteInsert, parList)
	return err
}

func (qry *QueryNote)doUpdateEx(notes ObjNoteList) error {
	parList := make([]dbase.ParamList, len(notes))
	for i, note := range notes {
		param := make(dbase.ParamList, 5)
		param[0] = &note.Count
		param[1] = &note.Amount
		param[2] = &note.Device
		param[3] = &note.Currency
		param[4] = &note.Nominal
		parList[i] = param
	}
	err := qry.RunPreparedSql(sqlNoteUpdate, parList)
	return err
}

func (qry *QueryNote)doUpdateAccept(device string, data common.ValidatorAccept) error {
	note := &ObjNote{
		Device:   device,
		Currency: uint16(data.Currency),
		Nominal:  float32(data.Nominal),
		Count:    uint16(data.Count),
		Amount:   float32(data.Amount),
	}
	return qry.doUpdate(note)
}

func (qry *QueryNote)doInsertNotes(device string, notes common.ValidNoteList) error {
	objList := make(ObjNoteList, len(notes))
	for i, note := range notes {
		objList[i] = &ObjNote {
			Device:   device,
			Currency: uint16(note.Currency),
			Nominal:  float32(note.Nominal),
			Count:    uint16(note.Count),
			Amount:   float32(note.Amount),
		}
	}
	return qry.doInsertEx(objList)
}

