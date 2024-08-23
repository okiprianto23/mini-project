package tx_helper

import (
	"database/sql"
	"main-xyz/context"
	"time"
)

type TXHelper interface {
	InitTXService(
		ctx *context.ContextModel,
		inputStruct interface{},
		serveFunction ServiceFunction) (
		*txBuilder,
		error)
}

type ServiceFunction func(*context.ContextModel, *sql.Tx, interface{}, time.Time) (interface{}, error)
type AfterCommit func(*context.ContextModel, interface{})
