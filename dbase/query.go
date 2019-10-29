package dbase

import (
	"database/sql"
	"errors"
	"github.com/iftsoft/device/core"
	"reflect"
)

type ParamList []interface{}

///////////////////////////////////////////////////////////////////////
//
// Common linker to database storage
//
type DBaseQuery struct {
	linker DBaseLinker	// Interface for database connection
	stmt   *sql.Stmt
	res    sql.Result
	count  int64
	log    *core.LogAgent
}

// Set connection to Query object
func (qry *DBaseQuery) InitQuery(link DBaseLinker, log *core.LogAgent) {
	qry.linker = link
	qry.stmt   = nil
	qry.res    = nil
	qry.count  = 0
	qry.log    = log
}

func (qry *DBaseQuery)GetCounter() int64 {
	return qry.count
}

func (qry *DBaseQuery)RowsAffected() int64 {
	if qry.res == nil {
		return -1
	}
	out, err := qry.res.RowsAffected()
	if err == nil {
		return out
	}
	return -1
}

func (qry *DBaseQuery)LastInsertId() int64 {
	if qry.res == nil {
		return -1
	}
	out, err := qry.res.LastInsertId()
	if err == nil {
		return out
	}
	return -1
}

///////////////////////////////////////////////////////////////////////
//
// Execute select SQL
func (qry *DBaseQuery) RunSelectSql(sqlText string, params ParamList, unit interface{}) (err error) {
	defer core.PanicRecover(&err, qry.log)
	if qry.linker == nil {
		return errors.New(errLinkerNotSet)
	}
	var row *sql.Rows
	row, err = qry.linker.DoQuery(sqlText, params...)
	if err == nil {
		qry.count, err = fetchSelectedRow(row, unit)
	}
	qry.log.Debug("SQL SelectOne return %s", core.GetErrorText(err))
	return err
}

// Execute search SQL
func (qry *DBaseQuery) RunSearchSql(sqlText string, params ParamList, list interface{}) (err error) {
	defer core.PanicRecover(&err, qry.log)
	if qry.linker == nil {
		return errors.New(errLinkerNotSet)
	}
	var row *sql.Rows
	row, err = qry.linker.DoQuery(sqlText, params...)
	if err == nil {
		qry.count, err = fetchSearchedRows(row, list)
	}
	qry.log.Debug("SQL SelectAll return %s", core.GetErrorText(err))
	return err
}

// Execute Return SQL
func (qry *DBaseQuery) RunReturnSql(sqlText string, params ParamList, m map[string]reflect.Value) (err error) {
	defer core.PanicRecover(&err, qry.log)
	if qry.linker == nil {
		return errors.New(errLinkerNotSet)
	}
	var row *sql.Rows
	row, err = qry.linker.DoQuery(sqlText, params...)
	if err == nil {
		qry.count, err = fetchReturnedRow(row, m)
	}
	qry.log.Debug("SQL ReturnRow return %s", core.GetErrorText(err))
	return err
}

// Execute command SQL
func (qry *DBaseQuery) RunCommandSql(sqlText string, params ParamList) (err error) {
	defer core.PanicRecover(&err, qry.log)
	if qry.linker == nil {
		return errors.New(errLinkerNotSet)
	}
	qry.count = 0
	qry.res, err = qry.linker.DoExec(sqlText, params...)
	qry.log.Debug("SQL ExecQuery return %s", core.GetErrorText(err))
	return err
}

// Execute prepared SQL statement
func (qry *DBaseQuery) RunPreparedSql(sqlText string, list []ParamList) (err error) {
	defer core.PanicRecover(&err, qry.log)
	if qry.linker == nil {
		return errors.New(errLinkerNotSet)
	}
	qry.count = 0
	qry.stmt, err = qry.linker.Prepare(sqlText)
	if err != nil {
		return err
	}
	defer func(){
		_ = qry.stmt.Close()
		qry.stmt = nil
	}()
	for i, params := range list {
		qry.res, err = qry.stmt.Exec(params...)
		if err == nil {
			qry.count += 1
			qry.log.Trace("SQL ExecPrepared line %d, return %s",
				i+1, core.GetErrorText(err))
		} else {
			break
		}
	}
	qry.log.Debug("SQL ExecPrepared return %s", core.GetErrorText(err))
	return err
}

