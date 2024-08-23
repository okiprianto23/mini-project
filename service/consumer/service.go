package consumer

import (
	"database/sql"
	"main-xyz/constanta"
	"main-xyz/context"
	"main-xyz/dto/out"
	error2 "main-xyz/error"
	"main-xyz/repository"
	"main-xyz/router"
	"main-xyz/utils"
	"main-xyz/utils/compare"
	"time"
)

func (c consumerService) InsertConsumer(ctx *context.ContextModel, _ router.URLParam, dto interface{}) (map[string]string, interface{}, error) {
	var (
		header map[string]string
		output interface{}
	)

	txBuilder, err := c.txHelper.InitTXService(ctx, dto, c.insertConsumer)
	if err != nil {
		return header, output, err
	}

	output, err = txBuilder.CompletedTXData()
	if err != nil {
		return header, output, err
	}

	return header, output, err
}

func (c consumerService) insertConsumer(ctx *context.ContextModel, tx *sql.Tx, dto interface{}, now time.Time) (output interface{}, err error) {

	inputStructMultipart := c.parseMultipartDTO(dto)
	inputStruct := inputStructMultipart.Consumer
	inputStruct.KTPPhoto = inputStructMultipart.KTPPhoto.FullPath
	inputStruct.SelfiePhoto = inputStructMultipart.SelfiePhoto.FullPath

	// validation data yg dimasukan apakah sudah valida atau belum
	var userDB repository.UserAdmin
	userDB, err = c.userDAO.GetUserByID(inputStruct.UserID)
	if err != nil {
		return
	}

	// check user id valid
	if compare.IsEmptyInt64(userDB.ID.Int64) {
		err = error2.ErrUnknownData.Param(constanta.UserID)
		return
	}

	// check apakah consumer dengan NIK yg sudah ada
	var isNIKExist bool
	isNIKExist, err = c.consumerDAO.CheckNIKConsumer(ctx, inputStruct.NIK)

	if isNIKExist {
		err = error2.ErrDataUsed.Param(constanta.Nik)
		return
	}

	//jika berhasil lakukan insert ke database
	//seblum itu kita rubah dto in menjadi repo model
	resultModel := c.consumerConvertToRepo(ctx, inputStruct, now)

	//insert to table consumer
	var resultID int64
	resultID, err = c.consumerDAO.InsertConsumer(ctx, tx, resultModel)
	if err != nil {
		return
	}

	resultModel.ID.Int64 = resultID

	// setelah di insert jalan kan untuk memasukan credit limit berdasarkan salary
	err = c.logicForCreditLimit(ctx, tx, resultModel, now)
	if err != nil {
		return
	}

	output = out.DefaultResponsePayloadMessage{
		Status: out.DefaultResponsePayloadStatus{
			Code:    "OK",
			Message: c.bundles.ReadMessageBundle("consumer", "SUCCESS_INSERT_MESSAGE", ctx.AuthAccessTokenModel.Locale, nil),
		},
		Data: resultID,
	}

	return
}

func (c consumerService) GetListConsumer(ctx *context.ContextModel, _ router.URLParam, dto interface{}) (header map[string]string, output interface{}, err error) {
	var dbResult []interface{}
	// get list query user
	dbResult, err = c.consumerDAO.GetListConsumer(ctx)
	if err != nil {
		return
	}

	var result []out.ConsumerOut
	for _, d := range dbResult {
		temp := d.(repository.ConsumerModel)
		result = append(result, out.ConsumerOut{
			ID:          temp.ID.Int64,
			UUID:        temp.UUID.String,
			NIK:         temp.NIK.String,
			FullName:    temp.FullName.String,
			LegalName:   temp.LegalName.String,
			BirthPlace:  temp.BirthPlace.String,
			BirthDate:   temp.BirthDate.Time.Format(time.DateOnly),
			Salary:      temp.Salary.Float64,
			KTPPhoto:    temp.KTPPhoto.String,
			SelfiePhoto: temp.SelfiePhoto.String,
		})
	}

	output = out.DefaultResponsePayloadMessage{
		Status: out.DefaultResponsePayloadStatus{
			Code:    "OK",
			Message: c.bundles.ReadMessageBundle("consumer", "SUCCESS_LIST_MESSAGE", ctx.AuthAccessTokenModel.Locale, nil),
		},
		Data: result,
	}

	return
}

func (c consumerService) GetListCreditByConsumerID(ctx *context.ContextModel, param router.URLParam, _ interface{}) (header map[string]string, output interface{}, err error) {
	var id int
	id, err = utils.CheckIDParam(param)
	if err != nil {
		return
	}

	var dbResult []interface{}
	// get list query user
	dbResult, err = c.creditLimitDAO.GetCreditLimitByConsumerID(ctx, int64(id))
	if err != nil {
		return
	}

	var result []out.CreditLimitOut
	for _, d := range dbResult {
		temp := d.(repository.CreditLimitModel)
		result = append(result, out.CreditLimitOut{
			ID:                  temp.ID.Int64,
			UUID:                temp.UUID.String,
			ConsumerID:          temp.ConsumerID.Int64,
			MonthlyInstallments: temp.MonthlyInstallments.Float64,
			InterestRate:        temp.InterestRate.Float64,
			Tenor:               temp.Tenor.Int64,
			LimitAmount:         temp.Limit.Float64,
		})
	}

	output = out.DefaultResponsePayloadMessage{
		Status: out.DefaultResponsePayloadStatus{
			Code:    "OK",
			Message: c.bundles.ReadMessageBundle("consumer", "SUCCESS_LIST_MESSAGE", ctx.AuthAccessTokenModel.Locale, nil),
		},
		Data: result,
	}

	return
}
