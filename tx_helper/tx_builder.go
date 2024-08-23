package tx_helper

import (
	"database/sql"
	"main-xyz/context"
	"time"
)

type txBuilder struct {
	ctx            *context.ContextModel
	inputStruct    interface{}
	serverFunction ServiceFunction
	afterCommit    AfterCommit
	db             *sql.Tx
	time           time.Time
}

func (tb *txBuilder) AfterCommit(
	f AfterCommit,
) *txBuilder {
	tb.afterCommit = f
	return tb
}

func (tb *txBuilder) CompletedTXData() (
	interface{},
	error,
) {
	return tb.doTXData()
}

func (tb *txBuilder) doTXData() (
	output interface{},
	err error,
) {

	defer func() {
		if err != nil {
			_ = tb.db.Rollback()
		} else {
			if tb.db != nil {
				err = tb.db.Commit()
				if err != nil {
					return
				}
			}

			if output != nil && tb.afterCommit != nil {
				go tb.afterCommit(tb.ctx, output)
			}

		}
	}()

	output, err = tb.serverFunction(tb.ctx, tb.db, tb.inputStruct, tb.time)

	if err != nil {
		return
	}

	return
}
