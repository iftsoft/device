package dbvalid

import (
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/dbase"
)

type QueryValidator struct {
	dbase.DBaseQuery
}

func NewQueryValidator(linker dbase.DBaseLinker, log *core.LogAgent) *QueryValidator {
	qry := &QueryValidator{}
	qry.InitQuery(linker, log)
	return qry
}

func (dao *QueryValidator)CreateTableNote() error {
	param := make(dbase.ParamList, 0)
	err := dao.RunCommandSql(sqlNoteCreate, param)
	return err
}

func (dao *QueryValidator)CreateTableBatch() error {
	param := make(dbase.ParamList, 0)
	err := dao.RunCommandSql(sqlBatchCreate, param)
	return err
}



