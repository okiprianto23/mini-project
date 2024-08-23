package transaction

import (
	"database/sql"
	"main-xyz/constanta"
	"main-xyz/context"
	"main-xyz/dto/out"
	error2 "main-xyz/error"
	"main-xyz/repository"
	"main-xyz/router"
	"time"
)

func (t transactionService) TransactionServiceInsert(ctx *context.ContextModel, _ router.URLParam, dto interface{}) (map[string]string, interface{}, error) {
	var (
		header map[string]string
		output interface{}
	)

	txBuilder, err := t.txHelper.InitTXService(ctx, dto, t.insertTransaction)
	if err != nil {
		return header, output, err
	}

	output, err = txBuilder.CompletedTXData()
	if err != nil {
		return header, output, err
	}

	return header, output, err
}

func (t transactionService) insertTransaction(ctx *context.ContextModel, tx *sql.Tx, dto interface{}, now time.Time) (output interface{}, err error) {

	inputStruct := t.parseDTO(dto)

	//check apakah consumernya ada dan get limit by tenor
	var resultLimit []interface{}
	resultLimit, err = t.creditLimitDAO.GetCreditLimitByConsumerID(ctx, inputStruct.ConsumerID)
	if err != nil {
		return
	}

	if len(resultLimit) < 1 {
		err = error2.ErrUnknownData.Param(constanta.ConsumerID)
		return
	}

	// jika ada maka kita hitung dulu harga yg mau di gunakan
	// otr + admin fee
	totalBarang := inputStruct.Otr + inputStruct.AdminFee

	// baru cari dari tenor terkecil buat ambil barangnya
	var tenorIndex int
	for i, cc := range resultLimit {
		temp := cc.(repository.CreditLimitModel)
		if temp.Limit.Float64 > totalBarang {
			tenorIndex = i
			continue
		}
	}

	if tenorIndex == 0 {
		err = error2.ErrUnknownData.Param("Credit")
		return
	}

	//set contact number
	inputStruct.ContractNumber, err = t.generateContractNumber("trc")
	if err != nil {
		return
	}

	//insert to transaction
	modelTransaction := repository.TransactionModel{
		ConsumerID:        sql.NullInt64{Int64: inputStruct.ConsumerID},
		LimitID:           resultLimit[tenorIndex].(repository.CreditLimitModel).ID,
		ContractNumber:    sql.NullString{String: inputStruct.ContractNumber},
		OTR:               sql.NullFloat64{Float64: inputStruct.Otr},
		AdminFee:          sql.NullFloat64{Float64: inputStruct.AdminFee},
		InstallmentAmount: resultLimit[tenorIndex].(repository.CreditLimitModel).MonthlyInstallments,
		InterestAmount:    resultLimit[tenorIndex].(repository.CreditLimitModel).InterestRate,
		AssetName:         sql.NullString{String: inputStruct.AssetName},
		TransactionDate:   sql.NullTime{Time: now},
		DefaultCreatedUpdated: repository.DefaultCreatedUpdated{
			UpdatedBy:     sql.NullInt64{Int64: ctx.AuthAccessTokenModel.ResourceUserID},
			UpdatedAt:     sql.NullTime{Time: now},
			UpdatedClient: sql.NullString{String: ctx.AuthAccessTokenModel.ClientID},
			CreatedBy:     sql.NullInt64{Int64: ctx.AuthAccessTokenModel.ResourceUserID},
			CreatedAt:     sql.NullTime{Time: now},
			CreatedClient: sql.NullString{String: ctx.AuthAccessTokenModel.ClientID},
		},
	}

	var resultID int64
	resultID, err = t.transactionDAO.InsertTransaction(ctx, tx, modelTransaction)
	if err != nil {
		return
	}

	// hitung nilai sisa credit
	remainingCredit := resultLimit[tenorIndex].(repository.CreditLimitModel).Limit.Float64 - totalBarang
	// updated credit limit
	err = t.creditLimitDAO.UpdateCreditLimitRemaining(ctx, tx, repository.CreditLimitModel{
		ID:             resultLimit[tenorIndex].(repository.CreditLimitModel).ID,
		ConsumerID:     sql.NullInt64{Int64: inputStruct.ConsumerID},
		RemainingLimit: sql.NullFloat64{Float64: remainingCredit},
	})
	if err != nil {
		return
	}

	output = out.DefaultResponsePayloadMessage{
		Status: out.DefaultResponsePayloadStatus{
			Code:    "OK",
			Message: t.bundles.ReadMessageBundle("transaction", "SUCCESS_INSERT_MESSAGE", ctx.AuthAccessTokenModel.Locale, nil),
		},
		Data: resultID,
	}

	return
}
