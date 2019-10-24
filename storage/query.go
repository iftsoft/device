package storage

import (
	"github.com/iftsoft/device/core"
	"database/sql"
	"reflect"
)

///////////////////////////////////////////////////////////////////////
//
// Common linker to database storage
//
type DBaseLink struct {
	conn DBaseLinker	// Interface for database connection
	stmt *sql.Stmt
	Done int64			// Rows affected or fetch count
	log  *core.LogAgent
}

// Set connection to DAO object
func (link *DBaseLink) SetConnection(conn DBaseLinker) {
	link.conn = conn
	link.stmt = nil
	link.Done = 0
	link.log  = core.GetLogAgent(core.LogLevelTrace, "DBlink")
}

///////////////////////////////////////////////////////////////////////
//
// Execute select SQL
func (link *DBaseLink) executeSelectSql(sqlText string, params []interface{}, unit interface{}) (err error) {
	defer core.PanicRecover(&err, link.log)

	var row *sql.Rows
	row, err = link.conn.DoQuery(sqlText, params...)
	if err == nil {
		link.Done, err = fetchSelectedRow(row, unit)
	}
	link.log.Debug("SQL SelectOne return %s", core.GetErrorText(err))
	return err
}

// Execute search SQL
func (link *DBaseLink) executeSearchSql(sqlText string, params []interface{}, list interface{}) (err error) {
	defer core.PanicRecover(&err, link.log)

	var row *sql.Rows
	row, err = link.conn.DoQuery(sqlText, params...)
	if err == nil {
		link.Done, err = fetchSearchedRows(row, list)
	}
	link.log.Debug("SQL SelectAll return %s", core.GetErrorText(err))
	return err
}

// Execute command SQL
func (link *DBaseLink) executeCommandSql(sqlText string, params []interface{}) (err error) {
	defer core.PanicRecover(&err, link.log)

	var res sql.Result
	res, err = link.conn.DoExec(sqlText, params...)
	if err == nil {
		link.Done, err = res.RowsAffected()
	}
	link.log.Debug("SQL ExecQuery return %s", core.GetErrorText(err))
	return err
}

// Execute Return SQL
func (link *DBaseLink) executeReturnSql(sqlText string, params []interface{}, m map[string]reflect.Value) (err error) {
	defer core.PanicRecover(&err, link.log)

	var row *sql.Rows
	row, err = link.conn.DoQuery(sqlText, params...)
	if err == nil {
		link.Done, err = fetchReturnedRow(row, m)
	}
	link.log.Debug("SQL ReturnRow return %s", core.GetErrorText(err))
	return err
}
