package dbvalid

import (
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/dbase"
	"time"
)

type ObjDeposit struct {
	Id       int
	BatchId  int
	ExtraId  int
	Device   string
	Currency uint16
	Nominal  float32
	Count    uint16
	Amount   float32
	Created  time.Time
}


type QueryDeposit struct {
	dbase.DBaseQuery
}

func NewQueryDeposit(linker dbase.DBaseLinker, log *core.LogAgent) *QueryDeposit {
	qry := &QueryDeposit{}
	qry.InitQuery(linker, log)
	return qry
}
