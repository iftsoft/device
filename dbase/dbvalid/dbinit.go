package dbvalid

import (
	"fmt"
	"github.com/iftsoft/device/common"
	"github.com/iftsoft/device/core"
	"github.com/iftsoft/device/dbase"
)

const timeFormat = "2006-01-02T15:04:05.999999Z07:00"


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

func (dao *QueryValidator)CreateTableDeposit() error {
	param := make(dbase.ParamList, 0)
	err := dao.RunCommandSql(sqlDepositCreate, param)
	return err
}

func (dao *QueryValidator)CreateTableBalance() error {
	param := make(dbase.ParamList, 0)
	err := dao.RunCommandSql(sqlBalanceCreate, param)
	return err
}


func checkBatch(notes common.ValidNoteList, depos ObjDepositList) (common.BatchState, string) {
	if  len(notes) == 0 &&
		len(depos) == 0 {
		return common.StateEmpty, ""
	}
	type currItem struct {
		note_cnt common.DevCounter
		note_sum common.DevAmount
		depo_cnt common.DevCounter
		depo_sum common.DevAmount
	}
	type currMap map[common.DevCurrency]*currItem
	set := currMap{}

	// process deposit list
	for _, note := range notes {
		item, ok := set[note.Currency]
		if !ok {
			item = &currItem{}
			set[note.Currency] = item
		}
		item.note_cnt += note.Count
		item.note_sum += note.Amount
	}

	// process deposit list
	for _, depo := range depos {
		item, ok := set[depo.Currency]
		if !ok {
			item = &currItem{}
			set[depo.Currency] = item
		}
		item.depo_cnt += depo.Count
		item.depo_sum += depo.Amount
	}

	// Compare items data
	state := common.StateCorrect
	brief := ""
	for curr, item := range set {
		if  item.note_cnt != item.depo_cnt ||
			item.note_sum != item.depo_sum {
			state = common.StateMismatch
			brief += fmt.Sprintf("Missmatch currency %3d (%s): Notes %4d /%9.2f != Depos %4d /%9.2f - %s; ",
				curr, curr.IsoCode(), item.note_cnt, item.note_sum, item.depo_cnt, item.depo_sum, curr.String() )
		}
	}
	return state, brief
}


