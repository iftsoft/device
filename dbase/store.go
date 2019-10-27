package dbase

import (
	"database/sql"
	"errors"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

const (
	errDBaseNotOpen = "Database is not opened"
	errLinkerNotSet = "Database linker is not set"
	errStmtNotReady = "Statement is not prepared"
)

type DBaseLinker interface {
	Open() error
	Close() error
	Begin() error
	Commit() error
	Rollback() error
	Prepare(query string) (*sql.Stmt, error)
	DoQuery(query string, args ...interface{}) (*sql.Rows, error)
	DoExec(query string, args ...interface{}) (sql.Result, error)
}

///////////////////////////////////////////////////////////////////////
//
// Database storage descriptor
//
type DBaseStore struct {
	file   string
	base   *sql.DB
	tran   *sql.Tx
	log    *core.LogAgent
}

func GetNewDBaseStore(cfg *StorageConfig) *DBaseStore {
	s := &DBaseStore{
		file: cfg.FileName,
		base: nil,
		tran: nil,
		log:  core.GetLogAgent(core.LogLevelTrace, "Sqlite"),
	}
	return s
}

// Open Database store from config data
func (s *DBaseStore) Open() (err error){
	err = s.Close()
	s.log.Info("Open database: %s", s.file)
	s.base, err = sql.Open("sqlite3", s.file)
	if err == nil {
		s.base.SetMaxIdleConns(3)
		s.base.SetMaxOpenConns(7)
		s.base.SetConnMaxLifetime(time.Duration(3) * time.Second)
	} else {
		s.log.Error("Return error: %s", err.Error())
	}
	return common.ExtendError(common.DevErrorDatabaseFault, err)
}

// Close Database store
func (s *DBaseStore) Close() (err error) {
	if s.tran != nil {
		err = s.tran.Commit()
		s.tran = nil
	}
	if s.base != nil {
		err = s.base.Close()
		s.base = nil
	}
	return common.ExtendError(common.DevErrorDatabaseFault, err)
}

// Get underlay sql.DB
func (s *DBaseStore) GetSqlDB() *sql.DB {
	return s.base
}


// Begin database transaction
func (s *DBaseStore) Begin() (err error) {
	if s.tran != nil {
		s.log.Warn("Transaction is already started")
		return err
	}
	s.log.Debug("Transaction begin")
	if s.base != nil {
		s.tran, err = s.base.Begin()
	} else {
		err = errors.New(errDBaseNotOpen)
	}
	if err != nil {
		s.log.Error("Transaction begin: %s", err.Error())
	}
	return common.ExtendError(common.DevErrorDatabaseFault, err)
}

// Commit database transaction
func (s *DBaseStore) Commit() (err error) {
	s.log.Debug("Transaction commit")
	if s.tran != nil {
		err = s.tran.Commit()
		if err != nil {
			s.log.Error("Transaction commit error: %s", err.Error())
		}
		s.tran = nil
	} else {
		err = errors.New("transaction was not started (commit)")
	}
	return common.ExtendError(common.DevErrorDatabaseFault, err)
}

// Rollback database transaction
func (s *DBaseStore) Rollback() (err error) {
	s.log.Debug("Transaction rollback")
	if s.tran != nil {
		err = s.tran.Rollback()
		if err != nil {
			s.log.Error("Transaction rollback error: %s", err.Error())
		}
		s.tran = nil
	} else {
		err = errors.New("transaction was not started (rollback)")
	}
	return common.ExtendError(common.DevErrorDatabaseFault, err)
}


// Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the returned statement.
func (s *DBaseStore) Prepare(query string) (*sql.Stmt, error) {
	if s.tran != nil {
		return s.tran.Prepare(query)
	}
	if s.base != nil {
		return s.base.Prepare(query)
	}
	return nil, errors.New(errDBaseNotOpen)
}

// doQuery executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (s *DBaseStore) DoQuery(query string, args ...interface{}) (*sql.Rows, error) {
	if s.tran != nil {
		return s.tran.Query(query, args...)
	}
	if s.base != nil {
		return s.base.Query(query, args...)
	}
	return nil, errors.New(errDBaseNotOpen)
}

// doExec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (s *DBaseStore) DoExec(query string, args ...interface{}) (sql.Result, error) {
	if s.tran != nil {
		return s.tran.Exec(query, args...)
	}
	if s.base != nil {
		return s.base.Exec(query, args...)
	}
	return nil, errors.New(errDBaseNotOpen)
}



