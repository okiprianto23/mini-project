package tx_helper

import (
	"database/sql"
	"main-xyz/context"
	"time"
)

func NewTXHelper(
	db *sql.DB,
) TXHelper {
	return &txHelper{
		db: db,
	}
}

type txHelper struct {
	db *sql.DB
}

func (t *txHelper) InitTXService(
	ctx *context.ContextModel,
	inputStruct interface{},
	serverFunction ServiceFunction,
) (
	*txBuilder,
	error,
) {

	builder := &txBuilder{
		ctx:            ctx,
		inputStruct:    inputStruct,
		time:           time.Now(),
		serverFunction: serverFunction,
	}

	if t.db != nil {
		tx, err := t.db.Begin()

		if err != nil {
			return nil, err
		}
		builder.db = tx
	}

	return builder, nil
}
