package dbvalid

import (
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/dbase"
	"time"
)

type ObjBalance struct {
	Id       int
	BatchId  int
	Device   string
	Currency uint16
	Nominal  float32
	Count    uint16
	Amount   float32
	Created  time.Time
}

type QueryBalance struct {
	dbase.DBaseQuery
}

func NewQueryBalance(linker dbase.DBaseLinker, log *core.LogAgent) *QueryBalance {
	qry := &QueryBalance{}
	qry.InitQuery(linker, log)
	return qry
}


